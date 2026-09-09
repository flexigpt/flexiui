package resource

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/definition"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
)

type ResolveOptions struct {
	VerifySourceContent bool `json:"verifySourceContent"`
}

// ResolvedArtifact contains the Store-verified resource chain for an
// available Artifact.
type ResolvedArtifact struct {
	Artifact         artifact.Artifact     `json:"artifact"`
	Collection       collection.Collection `json:"collection"`
	Definition       definition.Definition `json:"definition"`
	Occurrence       catalog.Occurrence    `json:"occurrence"`
	Source           source.Summary        `json:"source"`
	CatalogRevision  uint64                `json:"catalogRevision"`
	SourceGeneration string                `json:"sourceGeneration"`
}

func (r ResolvedArtifact) Validate() error {
	if err := r.Artifact.Validate(); err != nil {
		return err
	}
	if err := r.Collection.Validate(); err != nil {
		return err
	}
	if err := r.Definition.Validate(); err != nil {
		return err
	}
	if err := r.Occurrence.Validate(); err != nil {
		return err
	}
	if err := r.Source.Validate(); err != nil {
		return err
	}
	if r.CatalogRevision == 0 {
		return fmt.Errorf(
			"%w: resolved Artifact catalog revision is required",
			basespec.ErrInvalid,
		)
	}
	if err := basespec.ValidateSourceGeneration(r.SourceGeneration); err != nil {
		return err
	}
	if r.Artifact.RootID != r.Collection.RootID ||
		r.Artifact.CollectionID != r.Collection.ID {
		return fmt.Errorf(
			"%w: resolved Artifact belongs to another Collection",
			basespec.ErrInvalid,
		)
	}
	if r.Source.RootID != r.Collection.RootID ||
		r.Source.ID != r.Artifact.Binding.SourceID {
		return fmt.Errorf(
			"%w: resolved Source does not match Artifact binding",
			basespec.ErrInvalid,
		)
	}
	if r.Occurrence.RootID != r.Artifact.RootID ||
		r.Occurrence.CollectionID != r.Artifact.CollectionID ||
		r.Occurrence.Key.CollectionID != r.Artifact.CollectionID ||
		r.Occurrence.Key.SourceID != r.Artifact.Binding.SourceID ||
		r.Occurrence.Key.Locator != r.Artifact.Binding.Locator ||
		r.Occurrence.Key.SubresourceLocator != r.Artifact.Binding.SubresourceLocator {
		return fmt.Errorf(
			"%w: resolved occurrence does not match Artifact binding",
			basespec.ErrInvalid,
		)
	}
	if r.Occurrence.State != catalog.OccurrenceValid ||
		r.Occurrence.Kind != r.Artifact.Kind ||
		r.Occurrence.DefinitionDigest == nil ||
		r.Occurrence.SourceContentDigest == nil ||
		r.Artifact.ResolvedDefinition == nil {
		return fmt.Errorf(
			"%w: resolved Artifact has no current valid occurrence",
			basespec.ErrReferenceUnresolved,
		)
	}
	if r.Definition.Kind != r.Artifact.Kind ||
		r.Definition.Digest != *r.Artifact.ResolvedDefinition ||
		r.Definition.Digest != *r.Occurrence.DefinitionDigest {
		return fmt.Errorf(
			"%w: resolved Definition does not match Artifact state",
			basespec.ErrDigestMismatch,
		)
	}
	return nil
}

func (r ResolvedArtifact) Clone() ResolvedArtifact {
	output := r
	output.Artifact = r.Artifact.Clone()
	output.Collection = r.Collection.Clone()
	output.Definition = r.Definition.Clone()
	output.Occurrence = r.Occurrence.Clone()
	output.Source = r.Source.Clone()
	return output
}

type ArtifactResolutionStatus string

const (
	ArtifactResolutionRecordUnavailable     ArtifactResolutionStatus = "record-unavailable"
	ArtifactResolutionSourceUnavailable     ArtifactResolutionStatus = "source-unavailable"
	ArtifactResolutionOccurrenceUnavailable ArtifactResolutionStatus = "occurrence-unavailable"
	ArtifactResolutionDefinitionUnavailable ArtifactResolutionStatus = "definition-unavailable"
	ArtifactResolutionDefinitionMismatch    ArtifactResolutionStatus = "definition-mismatch"
)

// ArtifactResolutionIssue describes a generic Artifact Store link that could
// not be materialized from a Collection Catalog.
//
// It deliberately contains no Workspace-specific diagnostics or presentation
// strings. Feature consumers map the generic status into their own views.
type ArtifactResolutionIssue struct {
	Artifact artifact.Artifact
	Status   ArtifactResolutionStatus
}

func (i ArtifactResolutionIssue) Clone() ArtifactResolutionIssue {
	output := i
	output.Artifact = i.Artifact.Clone()
	return output
}

// CollectionArtifactResource joins one local Artifact record to its current
// Catalog occurrence, canonical definition, and attached Source summary.
//
// Resolved is populated only when the Catalog is current and the Artifact is
// currently available. Consumers can use it for runtime preparation without
// repeating the generic Artifact Store resource-chain lookup.
type CollectionArtifactResource struct {
	Artifact       artifact.Artifact
	Definition     definition.Definition
	Occurrence     catalog.Occurrence
	Source         source.Summary
	CatalogCurrent bool
	Resolved       *ResolvedArtifact
}

func (r CollectionArtifactResource) Clone() CollectionArtifactResource {
	output := r
	output.Artifact = r.Artifact.Clone()
	output.Definition = r.Definition.Clone()
	output.Occurrence = r.Occurrence.Clone()
	output.Source = r.Source.Clone()
	if r.Resolved != nil {
		value := r.Resolved.Clone()
		output.Resolved = &value
	}
	return output
}

// CollectionResourceInspection is the generic Artifact Store read model for
// one Collection. It separates resolved resources, unresolved local Artifact
// records, and unrecorded Catalog occurrences.
type CollectionResourceInspection struct {
	Catalog               catalog.CatalogInspection
	Resources             []CollectionArtifactResource
	UnresolvedArtifacts   []ArtifactResolutionIssue
	UnrecordedOccurrences []catalog.Occurrence
}

func (i CollectionResourceInspection) Clone() CollectionResourceInspection {
	output := i
	output.Catalog = i.Catalog.Clone()
	output.Resources = make(
		[]CollectionArtifactResource,
		len(i.Resources),
	)
	for index, value := range i.Resources {
		output.Resources[index] = value.Clone()
	}
	output.UnresolvedArtifacts = make(
		[]ArtifactResolutionIssue,
		len(i.UnresolvedArtifacts),
	)
	for index, value := range i.UnresolvedArtifacts {
		output.UnresolvedArtifacts[index] = value.Clone()
	}
	output.UnrecordedOccurrences = make(
		[]catalog.Occurrence,
		len(i.UnrecordedOccurrences),
	)
	for index, value := range i.UnrecordedOccurrences {
		output.UnrecordedOccurrences[index] = value.Clone()
	}
	return output
}

// VerifiedCollectionEntry binds a verified Collection entry to the exact
// current Catalog snapshot used to validate its Source revision and
// generation.
type VerifiedCollectionEntry struct {
	Catalog catalog.Snapshot
	Entry   VerifiedEntry
}

func (e VerifiedCollectionEntry) Validate() error {
	if err := e.Catalog.Validate(); err != nil {
		return err
	}
	if err := e.Entry.Validate(); err != nil {
		return err
	}
	if e.Catalog.RootID != e.Entry.Collection.RootID ||
		e.Catalog.CollectionID != e.Entry.Collection.CollectionID {
		return fmt.Errorf(
			"%w: verified entry Catalog belongs to another Collection",
			basespec.ErrInvalid,
		)
	}
	if e.Catalog.Revision != e.Entry.CatalogRevision {
		return fmt.Errorf(
			"%w: verified entry Catalog revision does not match snapshot",
			basespec.ErrInvalid,
		)
	}
	sourceRevision, found := e.Catalog.SourceRevisions[e.Entry.SourceID]
	if !found || sourceRevision != e.Entry.SourceRevision {
		return fmt.Errorf(
			"%w: verified entry Source revision does not match snapshot",
			basespec.ErrInvalid,
		)
	}
	sourceGeneration, found := e.Catalog.SourceGenerations[e.Entry.SourceID]
	if !found || sourceGeneration != e.Entry.SourceGeneration {
		return fmt.Errorf(
			"%w: verified entry Source generation does not match snapshot",
			basespec.ErrInvalid,
		)
	}
	return nil
}

func (e VerifiedCollectionEntry) Clone() VerifiedCollectionEntry {
	output := e
	output.Catalog = e.Catalog.Clone()
	output.Entry = e.Entry.Clone()
	return output
}

type VerifiedEntry struct {
	Collection       collection.CollectionRef `json:"collection"`
	SourceID         source.SourceID          `json:"sourceID"`
	CatalogRevision  uint64                   `json:"catalogRevision"`
	SourceRevision   uint64                   `json:"sourceRevision"`
	SourceGeneration string                   `json:"sourceGeneration"`
	Content          []byte                   `json:"content"`
	Digest           cryptoutil.Digest        `json:"digest"`
}

func (e VerifiedEntry) Validate() error {
	if err := e.Collection.Validate(); err != nil {
		return err
	}
	if err := e.SourceID.Validate(); err != nil {
		return err
	}
	if e.CatalogRevision == 0 || e.SourceRevision == 0 {
		return fmt.Errorf(
			"%w: verified entry revisions are required",
			basespec.ErrInvalid,
		)
	}
	if err := basespec.ValidateSourceGeneration(e.SourceGeneration); err != nil {
		return err
	}
	if err := cryptoutil.ValidateDigest(e.Digest); err != nil {
		return err
	}
	if cryptoutil.DigestBytes(e.Content) != e.Digest {
		return fmt.Errorf(
			"%w: verified entry content does not match digest",
			basespec.ErrDigestMismatch,
		)
	}
	return nil
}

func (e VerifiedEntry) Clone() VerifiedEntry {
	output := e
	output.Content = append([]byte(nil), e.Content...)
	return output
}
