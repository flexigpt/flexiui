package domain

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/definition"
)

// DefinitionForArtifact returns the catalog definition currently resolved by
// an Artifact. It verifies that the catalog and Artifact have matching
// Collection identities and matching definition digests.
func DefinitionForArtifact(
	snapshot catalog.Snapshot,
	record artifact.Artifact,
) (definition.Definition, error) {
	if record.ResolvedDefinition == nil {
		return definition.Definition{}, fmt.Errorf(
			"%w: Skill Artifact %q has no current definition",
			basespec.ErrReferenceUnresolved,
			record.ID,
		)
	}
	if snapshot.RootID != record.RootID ||
		snapshot.CollectionID != record.CollectionID {
		return definition.Definition{}, fmt.Errorf(
			"%w: Skill catalog belongs to another Collection",
			basespec.ErrInvalid,
		)
	}

	value, err := snapshot.DefinitionForOccurrence(catalog.OccurrenceKey{
		CollectionID:       record.CollectionID,
		SourceID:           record.Binding.SourceID,
		Locator:            record.Binding.Locator,
		SubresourceLocator: record.Binding.SubresourceLocator,
	})
	if err != nil {
		return definition.Definition{}, err
	}
	if value.Digest != *record.ResolvedDefinition {
		return definition.Definition{}, fmt.Errorf(
			"%w: Skill Artifact %q catalog definition changed",
			basespec.ErrConflict,
			record.ID,
		)
	}
	return value, nil
}
