package bundle

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

// StoreTopology is the normalized MCP Bundle topology projected from
// Artifact Store values.
//
// Generic Artifact Store entity validation is complete before this value is
// assembled. This type validates only MCP Bundle-specific relationships.
type StoreTopology struct {
	Collection  collection.CollectionRef
	Data        CollectionData
	Attachments []StoreAttachment
	Sources     []StoreSource
}

type StoreAttachment struct {
	RootID       root.RootID
	CollectionID collection.CollectionID
	SourceID     source.SourceID
	Role         collection.AttachmentRole
}

type StoreSource struct {
	ID     source.SourceID
	RootID root.RootID
	Kind   source.SourceKind
}

// ValidateStoreTopology validates MCP Bundle-specific attachment and Source
// invariants. It deliberately does not repeat generic Artifact Store model
// validation or CollectionData decoding.
func ValidateStoreTopology(value StoreTopology) error {
	if len(value.Attachments) != 1 {
		return fmt.Errorf(
			"%w: MCP Bundle must have exactly one Source Attachment",
			basespec.ErrInvalid,
		)
	}
	if len(value.Sources) != len(value.Attachments) {
		return fmt.Errorf(
			"%w: MCP Bundle topology Source count does not match attachments",
			basespec.ErrInvalid,
		)
	}

	attachment := value.Attachments[0]
	if attachment.RootID != value.Collection.RootID ||
		attachment.CollectionID != value.Collection.CollectionID {
		return fmt.Errorf(
			"%w: MCP Bundle attachment belongs to another Collection",
			basespec.ErrInvalid,
		)
	}
	switch attachment.Role {
	case artifactbuiltin.ManagedAttachmentRole,
		artifactbuiltin.BuiltInAttachmentRole:
	default:
		return fmt.Errorf(
			"%w: unsupported MCP Attachment role %q",
			basespec.ErrInvalid,
			attachment.Role,
		)
	}

	sourcesByID := make(
		map[source.SourceID]StoreSource,
		len(value.Sources),
	)
	for index, sourceValue := range value.Sources {
		if sourceValue.RootID != value.Collection.RootID {
			return fmt.Errorf(
				"%w: MCP Bundle Source %d belongs to another Root",
				basespec.ErrInvalid,
				index,
			)
		}
		if _, duplicate := sourcesByID[sourceValue.ID]; duplicate {
			return fmt.Errorf(
				"%w: MCP Bundle has duplicate Source %q",
				basespec.ErrInvalid,
				sourceValue.ID,
			)
		}
		sourcesByID[sourceValue.ID] = sourceValue
	}

	sourceValue, found := sourcesByID[attachment.SourceID]
	if !found {
		return fmt.Errorf(
			"%w: MCP Bundle Source %q is unavailable",
			basespec.ErrAttachmentNotFound,
			attachment.SourceID,
		)
	}
	if sourceValue.Kind != source.SourceKindManagedDirectory {
		return fmt.Errorf(
			"%w: MCP Bundle requires a managed Source",
			basespec.ErrInvalid,
		)
	}
	if value.Data.ManagedSourceID != "" &&
		value.Data.ManagedSourceID != sourceValue.ID {
		return fmt.Errorf(
			"%w: MCP Bundle managed Source ownership mismatch",
			basespec.ErrInvalid,
		)
	}
	return nil
}
