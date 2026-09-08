package domain

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

type AttachmentTopologyEntry struct {
	SourceID source.SourceID
	Role     collection.AttachmentRole
}

func ValidateAttachmentTopology(
	data CollectionData,
	attachments []AttachmentTopologyEntry,
) error {
	var (
		managedAttachmentCount int
		managedAttachmentID    source.SourceID
		builtInAttachmentCount int
	)

	for _, attachment := range attachments {
		switch attachment.Role {
		case artifactbuiltin.ManagedAttachmentRole:
			managedAttachmentCount++
			managedAttachmentID = attachment.SourceID

		case artifactbuiltin.BuiltInAttachmentRole:
			builtInAttachmentCount++
		}
	}

	if managedAttachmentCount > 1 {
		return fmt.Errorf(
			"%w: skill bundle has multiple managed attachments",
			basespec.ErrInvalid,
		)
	}
	if builtInAttachmentCount > 1 {
		return fmt.Errorf(
			"%w: skill bundle has multiple built-in attachments",
			basespec.ErrInvalid,
		)
	}

	if data.ManagedSourceID == "" {
		return nil
	}
	if managedAttachmentCount != 1 ||
		managedAttachmentID != data.ManagedSourceID {
		return fmt.Errorf(
			"%w: bundle-owned managed Source %q is not its sole managed attachment",
			basespec.ErrInvalid,
			data.ManagedSourceID,
		)
	}

	return nil
}
