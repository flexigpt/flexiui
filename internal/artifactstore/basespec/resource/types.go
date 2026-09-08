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
