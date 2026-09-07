package schema

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
)

type (
	// Kind is the entity-neutral kind portion of a schema key.
	//
	// The Entity field determines whether the Kind must satisfy CollectionKind
	// or ArtifactKind validation.
	Kind string

	EntityType string
	SchemaID   string
)

func (v SchemaID) Validate() error {
	return basespec.ValidateIdentifier("schema ID", string(v), basespec.MaxSchemaIDBytes)
}

const (
	EntityCollection EntityType = "collection"
	EntityArtifact   EntityType = "artifact"
)

type Key struct {
	Entity        EntityType `json:"entity"`
	Kind          Kind       `json:"kind"`
	SchemaID      SchemaID   `json:"schemaID"`
	SchemaVersion string     `json:"schemaVersion"`
}

func (k Key) Validate() error {
	switch k.Entity {
	case EntityCollection:
		if err := collection.CollectionKind(k.Kind).Validate(); err != nil {
			return err
		}

	case EntityArtifact:
		if err := artifact.ArtifactKind(k.Kind).Validate(); err != nil {
			return err
		}

	default:
		return fmt.Errorf(
			"%w: unsupported schema entity %q",
			basespec.ErrInvalid,
			k.Entity,
		)
	}

	if err := k.SchemaID.Validate(); err != nil {
		return err
	}
	return basespec.ValidateRequiredText(
		"schema version",
		k.SchemaVersion,
		basespec.MaxVersionBytes,
	)
}

func CollectionKey(
	kind collection.CollectionKind,
	schemaID SchemaID,
	schemaVersion string,
) Key {
	return Key{
		Entity:        EntityCollection,
		Kind:          Kind(kind),
		SchemaID:      schemaID,
		SchemaVersion: schemaVersion,
	}
}

func ArtifactKey(
	kind artifact.ArtifactKind,
	schemaID SchemaID,
	schemaVersion string,
) Key {
	return Key{
		Entity:        EntityArtifact,
		Kind:          Kind(kind),
		SchemaID:      schemaID,
		SchemaVersion: schemaVersion,
	}
}
