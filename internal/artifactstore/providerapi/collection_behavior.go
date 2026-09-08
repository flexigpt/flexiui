package providerapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/definition"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

// Collection is a provider-safe view of a persisted Collection.
//
// It deliberately excludes repositories, SQLite handles, source config, and
// feature-private runtime state.
type Collection struct {
	ID          collection.CollectionID   `json:"id"`
	RootID      root.RootID               `json:"rootID"`
	Kind        collection.CollectionKind `json:"kind"`
	DisplayName string                    `json:"displayName"`
	Description string                    `json:"description,omitempty"`
	Enabled     bool                      `json:"enabled"`
	Revision    uint64                    `json:"revision"`
	Data        json.RawMessage           `json:"data"`
}

func (c Collection) Clone() Collection {
	output := c
	output.Data = append(json.RawMessage(nil), c.Data...)
	return output
}

// Attachment is a provider-safe view of a persisted Collection attachment.
type Attachment struct {
	RootID       root.RootID               `json:"rootID"`
	CollectionID collection.CollectionID   `json:"collectionID"`
	SourceID     source.SourceID           `json:"sourceID"`
	Role         collection.AttachmentRole `json:"role"`
	Enabled      bool                      `json:"enabled"`
	Revision     uint64                    `json:"revision"`
	Data         json.RawMessage           `json:"data"`
}

func (a Attachment) Clone() Attachment {
	output := a
	output.Data = append(json.RawMessage(nil), a.Data...)
	return output
}

// Source is a provider-safe summary of an attached Source.
//
// It deliberately excludes Source.Config. A collection provider decides
// discovery semantics from source identity and kind, while Artifact Store owns
// source configuration, snapshot opening, and filesystem access.
type Source struct {
	ID          source.SourceID     `json:"id"`
	RootID      root.RootID         `json:"rootID"`
	StorageKey  basespec.StorageKey `json:"storageKey"`
	Kind        source.SourceKind   `json:"kind"`
	DisplayName string              `json:"displayName"`
	Enabled     bool                `json:"enabled"`
	Revision    uint64              `json:"revision"`
}

// Occurrence is the provider-safe observation supplied to automatic adoption.
type Occurrence struct {
	RootID             root.RootID                 `json:"rootID"`
	CollectionID       collection.CollectionID     `json:"collectionID"`
	SourceID           source.SourceID             `json:"sourceID"`
	Locator            basespec.Locator            `json:"locator"`
	SubresourceLocator basespec.SubresourceLocator `json:"subresourceLocator,omitempty"`
	Kind               artifact.ArtifactKind       `json:"kind"`
}

// AdoptionInput contains one valid source occurrence that Artifact Store is
// considering for automatic artifact adoption.
type AdoptionInput struct {
	Collection Collection
	Attachment Attachment
	Occurrence Occurrence
	Definition definition.Definition
}

func (i AdoptionInput) Clone() AdoptionInput {
	output := i
	output.Collection = i.Collection.Clone()
	output.Attachment = i.Attachment.Clone()
	output.Definition = i.Definition.Clone()
	return output
}

// AdoptionDecision is provider-owned semantic intent. Artifact Store owns
// artifact ID generation, artifact construction, validation, persistence,
// reconciliation, revisions, and timestamps.
type AdoptionDecision struct {
	Adopt       bool
	Name        string
	Enabled     bool
	Data        json.RawMessage
	Diagnostics []diagnostic.Diagnostic
}

func (d AdoptionDecision) Clone() AdoptionDecision {
	output := d
	output.Data = append(json.RawMessage(nil), d.Data...)
	output.Diagnostics = diagnostic.Clone(d.Diagnostics)
	return output
}

func (d AdoptionDecision) Validate() error {
	if err := diagnostic.Validate(d.Diagnostics); err != nil {
		return err
	}
	if !d.Adopt {
		return nil
	}
	return basespec.ValidateRequiredText(
		"automatically adopted artifact name",
		d.Name,
		basespec.MaxDisplayNameBytes,
	)
}

// CollectionBehavior is the common Artifact Store inbound behavior for one
// CollectionKind. A behavior must additionally implement exactly one planning
// role below.
//
// It receives immutable generic views and returns declarations. It must not
// perform direct Artifact Store mutations or depend on system.Components.
type CollectionBehavior interface {
	CollectionLifecycle

	CollectionKind() collection.CollectionKind

	// Revision must change whenever provider behavior can alter discovery
	// scope, decoder selection, automatic-adoption eligibility, or default
	// local artifact data.
	Revision() string

	DecideAutomaticAdoption(
		ctx context.Context,
		input AdoptionInput,
	) (AdoptionDecision, error)
}

// CollectionPlanner is the normal provider planning role. It receives only
// persisted generic views and returns a declarative discovery plan.
type CollectionPlanner interface {
	CollectionBehavior

	BuildDiscoveryPlan(
		ctx context.Context,
		collection Collection,
		attachments []Attachment,
		sources []Source,
	) (Plan, error)
}

// DocumentPlanningBehavior is the exceptional planning role for collection
// kinds whose declared discovery scope depends on one source-owned canonical
// document. Workspace implements this for workspace.json.
//
// The supplied reader is intentionally narrower than source.Runtime. It can
// read only canonical documents through Artifact Store controls.
type DocumentPlanningBehavior interface {
	CollectionBehavior

	BuildDiscoveryPlanWithDocuments(
		ctx context.Context,
		collection Collection,
		attachments []Attachment,
		sources []Source,
		documents PlanningDocumentReader,
	) (Plan, error)
}

func ValidateCollectionBehavior(
	behavior CollectionBehavior,
) error {
	if behavior == nil {
		return fmt.Errorf(
			"%w: collection behavior is nil",
			basespec.ErrInvalid,
		)
	}
	if err := behavior.CollectionKind().Validate(); err != nil {
		return err
	}
	_, plainPlanner := behavior.(CollectionPlanner)
	_, documentPlanner := behavior.(DocumentPlanningBehavior)
	if plainPlanner == documentPlanner {
		return fmt.Errorf(
			"%w: collection behavior must implement exactly one planning role",
			basespec.ErrInvalid,
		)
	}
	return basespec.ValidateRequiredText(
		"collection behavior revision",
		behavior.Revision(),
		basespec.MaxVersionBytes,
	)
}
