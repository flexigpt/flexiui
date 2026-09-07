package schema

import (
	"encoding/json"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
)

type EntityType string

const (
	EntityCollection EntityType = "collection"
	EntityArtifact   EntityType = "artifact"
)

// Kind is the entity-neutral kind portion of a schema key.
//
// The Entity field determines whether the Kind must satisfy CollectionKind
// or ArtifactKind validation.
type Kind string

type Key struct {
	Entity        EntityType        `json:"entity"`
	Kind          Kind              `json:"kind"`
	SchemaID      basespec.SchemaID `json:"schemaID"`
	SchemaVersion string            `json:"schemaVersion"`
}

func CollectionKey(
	kind basespec.CollectionKind,
	schemaID basespec.SchemaID,
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
	kind basespec.ArtifactKind,
	schemaID basespec.SchemaID,
	schemaVersion string,
) Key {
	return Key{
		Entity:        EntityArtifact,
		Kind:          Kind(kind),
		SchemaID:      schemaID,
		SchemaVersion: schemaVersion,
	}
}

func (k Key) Validate() error {
	switch k.Entity {
	case EntityCollection:
		if err := basespec.ValidateCollectionKind(
			basespec.CollectionKind(k.Kind),
		); err != nil {
			return err
		}

	case EntityArtifact:
		if err := basespec.ValidateArtifactKind(
			basespec.ArtifactKind(k.Kind),
		); err != nil {
			return err
		}

	default:
		return fmt.Errorf(
			"%w: unsupported schema entity %q",
			basespec.ErrInvalid,
			k.Entity,
		)
	}

	if err := basespec.ValidateSchemaID(k.SchemaID); err != nil {
		return err
	}
	return basespec.ValidateRequiredText(
		"schema version",
		k.SchemaVersion,
		basespec.MaxVersionBytes,
	)
}

// ParsedDocument is a canonical document accepted by a registered schema
// codec. It is a Store domain value, not a provider-private type.
type ParsedDocument struct {
	Key    Key
	Digest cryptoutil.Digest
	Raw    json.RawMessage
}

func (d ParsedDocument) Validate() error {
	if err := d.Key.Validate(); err != nil {
		return err
	}
	if err := cryptoutil.ValidateDigest(d.Digest); err != nil {
		return err
	}
	if len(d.Raw) == 0 {
		return fmt.Errorf(
			"%w: parsed document is empty",
			basespec.ErrInvalid,
		)
	}
	return nil
}

func (d ParsedDocument) Clone() ParsedDocument {
	output := d
	output.Raw = append(json.RawMessage(nil), d.Raw...)
	return output
}
