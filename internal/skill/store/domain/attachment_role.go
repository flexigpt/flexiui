package domain

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

const (
	DiscoveryPolicyRevision = "skill.bundle.discovery.v1"

	RoleExternal collection.AttachmentRole = "external"
	RoleLibrary  collection.AttachmentRole = "library"
)

func ValidateAttachmentRole(role collection.AttachmentRole) error {
	switch role {
	case artifactbuiltin.ManagedAttachmentRole,
		artifactbuiltin.BuiltInAttachmentRole,
		RoleExternal,
		RoleLibrary:
		return nil

	default:
		return fmt.Errorf(
			"%w: unsupported skill bundle attachment role %q",
			basespec.ErrInvalid,
			role,
		)
	}
}

func ValidateAttachmentRoleSourceKind(
	role collection.AttachmentRole,
	kind source.SourceKind,
) error {
	switch role {
	case artifactbuiltin.ManagedAttachmentRole,
		artifactbuiltin.BuiltInAttachmentRole:
		if kind != source.SourceKindManagedDirectory {
			return fmt.Errorf(
				"%w: skill bundle role %q requires source kind %q",
				basespec.ErrInvalid,
				role,
				source.SourceKindManagedDirectory,
			)
		}

	case RoleExternal, RoleLibrary:
		if kind != source.SourceKindFilesystemDirectory {
			return fmt.Errorf(
				"%w: skill bundle role %q requires source kind %q",
				basespec.ErrInvalid,
				role,
				source.SourceKindFilesystemDirectory,
			)
		}
	}

	return nil
}
