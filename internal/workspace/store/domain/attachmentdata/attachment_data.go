package attachmentdata

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
)

type AttachmentOperation struct {
	Role                                 collection.AttachmentRole
	CanAttach                            bool
	IsPrimary                            bool
	RequiredSourceKind                   source.SourceKind
	DefaultAuthoritative                 bool
	IncludeReadmeWhenRequested           bool
	AppliesWorkspaceDiscoveryPreferences bool
	AllowsAttachmentDiscoveryOverrides   bool
}

// attachmentOperationMatrix is the workspace attachment lifecycle and
// discovery-operation matrix.
//
// A role must be present here before it can be attached, validated, or planned.
var attachmentOperationMatrix = [...]AttachmentOperation{
	{
		Role:                                 workspaceDomain.RolePrimary,
		IsPrimary:                            true,
		RequiredSourceKind:                   source.SourceKindFilesystemDirectory,
		DefaultAuthoritative:                 true,
		IncludeReadmeWhenRequested:           true,
		AppliesWorkspaceDiscoveryPreferences: true,
	},
	{
		Role:                               workspaceDomain.RoleLibrary,
		CanAttach:                          true,
		DefaultAuthoritative:               true,
		AllowsAttachmentDiscoveryOverrides: true,
	},
	{
		Role:                               workspaceDomain.RoleAttachedPackage,
		CanAttach:                          true,
		DefaultAuthoritative:               true,
		AllowsAttachmentDiscoveryOverrides: true,
	},
	{
		Role:                               workspaceDomain.RoleOverlay,
		CanAttach:                          true,
		DefaultAuthoritative:               true,
		AllowsAttachmentDiscoveryOverrides: true,
	},
}

func ValidateAttachmentDataForRole(
	role collection.AttachmentRole,
	value workspaceDomain.AttachmentData,
) error {
	operation, supported := AttachmentOperationFor(role)
	if !supported {
		return fmt.Errorf(
			"%w: unsupported attachment role %q",
			workspaceDomain.ErrInvalidWorkspace,
			role,
		)
	}
	if !operation.AllowsAttachmentDiscoveryOverrides &&
		(value.Recursive != nil || value.Authoritative != nil) {
		return fmt.Errorf(
			"%w: attachment role %q does not allow discovery overrides",
			workspaceDomain.ErrInvalidWorkspace,
			role,
		)
	}
	return nil
}

func EncodeAttachmentData(
	value workspaceDomain.AttachmentData,
) (json.RawMessage, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	canonical, err := jsonutil.CanonicalizeObject(
		raw,
		basespec.MaxLocalDataBytes,
	)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(canonical), nil
}

func DecodeAttachmentData(
	raw json.RawMessage,
) (workspaceDomain.AttachmentData, error) {
	canonical, err := jsonutil.CanonicalizeObject(
		raw,
		basespec.MaxLocalDataBytes,
	)
	if err != nil {
		return workspaceDomain.AttachmentData{}, err
	}
	decoder := json.NewDecoder(bytes.NewReader(canonical))
	decoder.DisallowUnknownFields()
	var value workspaceDomain.AttachmentData
	if err := decoder.Decode(&value); err != nil {
		return workspaceDomain.AttachmentData{}, err
	}
	return value, nil
}

func AttachmentOperationFor(
	role collection.AttachmentRole,
) (AttachmentOperation, bool) {
	for _, operation := range attachmentOperationMatrix {
		if operation.Role == role {
			return operation, true
		}
	}
	return AttachmentOperation{}, false
}
