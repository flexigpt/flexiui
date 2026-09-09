package contextdomain

import (
	"fmt"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/definition"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
)

const maxWorkspaceContextContentBytes = 2 << 20

// ContextFromDefinition validates a canonical Workspace Context Artifact
// definition and returns its typed body in a single decode.
func ContextFromDefinition(
	value definition.Definition,
) (Definition, error) {
	if value.Kind != artifactbuiltin.WorkspaceContextArtifactKind {
		return Definition{}, fmt.Errorf(
			"%w: Context definition kind must be %q",
			workspaceDomain.ErrInvalidWorkspace,
			artifactbuiltin.WorkspaceContextArtifactKind,
		)
	}
	if value.SchemaID != artifactbuiltin.WorkspaceContextSchemaID {
		return Definition{}, fmt.Errorf(
			"%w: Context definition schema must be %q",
			workspaceDomain.ErrInvalidWorkspace,
			artifactbuiltin.WorkspaceContextSchemaID,
		)
	}
	if value.SchemaVersion != artifactbuiltin.WorkspaceContextSchemaVersion {
		return Definition{}, fmt.Errorf(
			"%w: Context definition schema version must be %q",
			workspaceDomain.ErrInvalidWorkspace,
			artifactbuiltin.WorkspaceContextSchemaVersion,
		)
	}
	if len(value.Dependencies) != 0 {
		return Definition{}, fmt.Errorf(
			"%w: Context definitions cannot declare dependencies",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}

	body, err := definition.DecodeBody[Definition](value.Body)
	if err != nil {
		return Definition{}, err
	}
	if err := basespec.ValidateRequiredText(
		"Context name",
		body.Name,
		basespec.MaxDisplayNameBytes,
	); err != nil {
		return Definition{}, err
	}
	if !supportedContextRole(body.Role) {
		return Definition{}, fmt.Errorf(
			"%w: unsupported Context role %q",
			workspaceDomain.ErrInvalidWorkspace,
			body.Role,
		)
	}
	if body.MediaType != artifactbuiltin.WorkspaceContextMediaTypeMarkdown {
		return Definition{}, fmt.Errorf(
			"%w: unsupported Context media type %q",
			workspaceDomain.ErrInvalidWorkspace,
			body.MediaType,
		)
	}
	if !utf8.ValidString(body.Content) {
		return Definition{}, fmt.Errorf(
			"%w: Context content must contain valid UTF-8",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}
	if strings.ContainsRune(body.Content, 0) {
		return Definition{}, fmt.Errorf(
			"%w: Context content contains a NUL byte",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}
	if strings.TrimSpace(body.Content) == "" {
		return Definition{}, fmt.Errorf(
			"%w: Context content is empty",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}
	if len(body.Content) > maxWorkspaceContextContentBytes {
		return Definition{}, fmt.Errorf(
			"%w: Context content exceeds %d bytes",
			workspaceDomain.ErrInvalidWorkspace,
			maxWorkspaceContextContentBytes,
		)
	}
	if value.DisplayName != body.Name {
		return Definition{}, fmt.Errorf(
			"%w: Context display name does not match body.name",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}
	if value.LogicalName != LogicalName(body.Name) {
		return Definition{}, fmt.Errorf(
			"%w: Context logical name does not match body.name",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}
	if value.Labels[artifactbuiltin.WorkspaceContextRoleLabelKey] != string(body.Role) {
		return Definition{}, fmt.Errorf(
			"%w: Context role label does not match body.role",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}
	return body, nil
}

func LogicalName(name string) basespec.LogicalName {
	contextVal := "context"
	base := strings.ToLower(strings.TrimSuffix(name, path.Ext(name)))
	parts := strings.FieldsFunc(base, func(character rune) bool {
		return (character < 'a' || character > 'z') &&
			(character < '0' || character > '9')
	})

	value := strings.Join(parts, "-")
	if value == "" {
		value = contextVal
	}
	if value[0] >= '0' && value[0] <= '9' {
		value = "context-" + value
	}
	if len(value) > basespec.MaxLogicalNameBytes {
		value = value[:basespec.MaxLogicalNameBytes]
		value = strings.Trim(value, ".-")
	}
	if value == "" {
		value = contextVal
	}
	return basespec.LogicalName(value)
}
