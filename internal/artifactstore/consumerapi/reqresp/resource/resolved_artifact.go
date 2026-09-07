package resource

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
)

// ResolveOptions controls Store-owned source verification.
//
// Catalog and definition verification always occur. VerifySourceContent also
// confirms the exact current source bytes before a resolved resource is
// returned.
type ResolveOptions struct {
	VerifySourceContent bool
}

// ResolvedArtifact is the consumer-facing Store-verified resource chain:
//
// Artifact -> Collection -> current Catalog -> Definition -> Source.
//
// Source configuration and source snapshots remain private to Artifact Store.
type ResolvedArtifact struct {
	Artifact         artifact.Artifact
	Collection       collection.Collection
	Definition       providerapi.Definition
	Occurrence       catalog.Occurrence
	Source           source.Summary
	CatalogRevision  uint64
	SourceGeneration string
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
			"%w: resolved Source does not match the Artifact binding",
			basespec.ErrInvalid,
		)
	}
	if r.Occurrence.RootID != r.Artifact.RootID ||
		r.Occurrence.CollectionID != r.Artifact.CollectionID ||
		r.Occurrence.Key.CollectionID != r.Artifact.CollectionID ||
		r.Occurrence.Key.SourceID != r.Artifact.Binding.SourceID ||
		r.Occurrence.Key.Locator != r.Artifact.Binding.Locator ||
		r.Occurrence.Key.SubresourceLocator !=
			r.Artifact.Binding.SubresourceLocator {
		return fmt.Errorf(
			"%w: resolved occurrence does not match the Artifact binding",
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
			"%w: resolved definition does not match Artifact state",
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
