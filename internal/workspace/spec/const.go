package spec

import "github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"

const (
	RolePrimary         collection.AttachmentRole = "primary"
	RoleLibrary         collection.AttachmentRole = "library"
	RoleAttachedPackage collection.AttachmentRole = "attached-package"
	RoleOverlay         collection.AttachmentRole = "overlay"
)
