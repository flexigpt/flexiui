package artifact

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
	"github.com/flexigpt/flexigpt-app/internal/uuidutil"
)

type (
	ArtifactID   string
	ArtifactKind string
	AdoptionMode string
)

const (
	AdoptionObserved AdoptionMode = "observed"
	AdoptionPinned   AdoptionMode = "pinned"
)

func (v ArtifactID) Validate() error {
	err := uuidutil.ValidateUUIDv7(string(v))
	if err != nil {
		return fmt.Errorf("artifact ID: %w", err)
	}
	return nil
}

func (v ArtifactKind) Validate() error {
	return basespec.ValidateIdentifier("artifact kind", string(v), basespec.MaxKindBytes)
}

type ArtifactRef struct {
	RootID     root.RootID `json:"rootID"`
	ArtifactID ArtifactID  `json:"artifactID"`
}

func (r ArtifactRef) Validate() error {
	if err := r.RootID.Validate(); err != nil {
		return err
	}
	return r.ArtifactID.Validate()
}

type ArtifactAddress struct {
	RootID       root.RootID             `json:"rootID"`
	CollectionID collection.CollectionID `json:"collectionID"`
	ArtifactID   ArtifactID              `json:"artifactID"`
	Kind         ArtifactKind            `json:"kind"`
}

func (a ArtifactAddress) Validate() error {
	if err := a.RootID.Validate(); err != nil {
		return err
	}
	if err := a.CollectionID.Validate(); err != nil {
		return err
	}
	if err := a.ArtifactID.Validate(); err != nil {
		return err
	}
	return a.Kind.Validate()
}

func (a ArtifactAddress) CollectionRef() collection.CollectionRef {
	return collection.CollectionRef{
		RootID:       a.RootID,
		CollectionID: a.CollectionID,
	}
}

type Artifact struct {
	ID                 ArtifactID              `json:"id"`
	RootID             root.RootID             `json:"rootID"`
	CollectionID       collection.CollectionID `json:"collectionID"`
	Binding            SourceBinding           `json:"binding"`
	Kind               ArtifactKind            `json:"kind"`
	Name               string                  `json:"name"`
	Enabled            bool                    `json:"enabled"`
	Adoption           AdoptionMode            `json:"adoption"`
	ResolvedDefinition *cryptoutil.Digest      `json:"resolvedDefinition,omitempty"`
	Data               json.RawMessage         `json:"-"`

	State       State                   `json:"state"`
	Diagnostics []diagnostic.Diagnostic `json:"diagnostics,omitempty"`
	Revision    uint64                  `json:"revision"`
	CreatedAt   time.Time               `json:"createdAt"`
	ModifiedAt  time.Time               `json:"modifiedAt"`
}

func (a Artifact) Ref() ArtifactRef {
	return ArtifactRef{
		RootID:     a.RootID,
		ArtifactID: a.ID,
	}
}

func (a Artifact) Address() ArtifactAddress {
	return ArtifactAddress{
		RootID:       a.RootID,
		CollectionID: a.CollectionID,
		ArtifactID:   a.ID,
		Kind:         a.Kind,
	}
}

func (a Artifact) Validate() error {
	if err := a.ID.Validate(); err != nil {
		return err
	}
	if err := a.RootID.Validate(); err != nil {
		return err
	}
	if err := a.CollectionID.Validate(); err != nil {
		return err
	}
	if err := a.Binding.Validate(); err != nil {
		return err
	}
	if err := a.Kind.Validate(); err != nil {
		return err
	}
	if a.Binding.ExpectedKind != a.Kind {
		return fmt.Errorf(
			"%w: artifact binding expected kind does not match artifact kind",
			basespec.ErrInvalid,
		)
	}
	if err := basespec.ValidateRequiredText(
		"artifact name",
		a.Name,
		basespec.MaxDisplayNameBytes,
	); err != nil {
		return err
	}
	switch a.Adoption {
	case AdoptionObserved, AdoptionPinned:
	default:
		return fmt.Errorf(
			"%w: invalid artifact adoption mode %q",
			basespec.ErrInvalid,
			a.Adoption,
		)
	}
	if err := a.State.Validate(a.ResolvedDefinition); err != nil {
		return err
	}
	if _, err := jsonutil.CanonicalizeObject(
		a.Data,
		basespec.MaxLocalDataBytes,
	); err != nil {
		return fmt.Errorf("%w: artifact data: %w", basespec.ErrInvalid, err)
	}
	if err := diagnostic.Validate(a.Diagnostics); err != nil {
		return err
	}
	if a.Revision == 0 {
		return fmt.Errorf(
			"%w: artifact revision must be positive",
			basespec.ErrInvalid,
		)
	}
	if a.CreatedAt.IsZero() || a.ModifiedAt.IsZero() {
		return fmt.Errorf(
			"%w: artifact timestamps are required",
			basespec.ErrInvalid,
		)
	}
	if a.ModifiedAt.Before(a.CreatedAt) {
		return fmt.Errorf(
			"%w: artifact modified time precedes creation",
			basespec.ErrInvalid,
		)
	}
	return nil
}

func (a Artifact) Clone() Artifact {
	output := a
	output.Data = append(json.RawMessage(nil), a.Data...)
	output.Diagnostics = diagnostic.Clone(a.Diagnostics)
	if a.ResolvedDefinition != nil {
		value := *a.ResolvedDefinition
		output.ResolvedDefinition = &value
	}
	return output
}
