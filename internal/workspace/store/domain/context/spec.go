package contextdomain

import (
	"strings"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
)

type contextFileSupport struct {
	FileName         string
	Role             artifactbuiltin.WorkspaceContextRole
	DefaultDiscovery bool
	Preference       artifactbuiltin.WorkspaceContextPreference
	RuntimeOrder     int
}

var contextConventionRegistry = func() []contextFileSupport {
	input := artifactbuiltin.WorkspaceContextFileConventions()
	output := make([]contextFileSupport, 0, len(input))
	for _, value := range input {
		output = append(output, contextFileSupport{
			FileName:         string(value.FileName),
			Role:             value.Role,
			DefaultDiscovery: value.DefaultDiscovery,
			Preference:       value.Preference,
			RuntimeOrder:     value.RuntimeOrder,
		})
	}
	return output
}()

func contextConventionFor(
	locator basespec.Locator,
) (contextFileSupport, bool) {
	value := string(locator)
	for _, convention := range contextConventionRegistry {
		if strings.EqualFold(value, convention.FileName) {
			return convention, true
		}
	}
	return contextFileSupport{}, false
}

func DiscoveryProfile() workspaceDomain.DiscoveryProfile {
	var profile workspaceDomain.DiscoveryProfile

	for _, convention := range contextConventionRegistry {
		locator := basespec.Locator(convention.FileName)
		switch {
		case convention.DefaultDiscovery:
			profile.ExplicitLocators = append(
				profile.ExplicitLocators,
				locator,
			)
		case convention.Preference == artifactbuiltin.WorkspaceContextPreferenceIncludeReadme:
			profile.ReadmeLocator = locator
		}
	}

	return profile
}

func RoleFor(
	locator basespec.Locator,
) (artifactbuiltin.WorkspaceContextRole, bool) {
	convention, found := contextConventionFor(locator)
	if !found {
		return "", false
	}
	return convention.Role, true
}

func RuntimeOrder(locator basespec.Locator) int {
	if convention, found := contextConventionFor(locator); found {
		return convention.RuntimeOrder
	}
	return 10_000
}

func supportedContextRole(role artifactbuiltin.WorkspaceContextRole) bool {
	switch role {
	case artifactbuiltin.WorkspaceContextRoleAgentInstructions,
		artifactbuiltin.WorkspaceContextRoleAssistantInstructions,
		artifactbuiltin.WorkspaceContextRoleProjectReadme,
		artifactbuiltin.WorkspaceContextRoleProjectContext:
		return true
	default:
		return false
	}
}

var artifactSupport = workspaceDomain.ArtifactSupport{
	Kind:      artifactbuiltin.WorkspaceContextArtifactKind,
	SchemaID:  artifactbuiltin.WorkspaceContextSchemaID,
	DecoderID: artifactbuiltin.WorkspaceContextDecoderID,
	Validator: ValidateContextDefinition,
}

func ArtifactSupport() workspaceDomain.ArtifactSupport {
	return artifactSupport
}

type contextDefinition struct {
	Name      string                                    `json:"name"`
	Role      artifactbuiltin.WorkspaceContextRole      `json:"role"`
	MediaType artifactbuiltin.WorkspaceContextMediaType `json:"mediaType"`
	Content   string                                    `json:"content"`
}

type Definition = contextDefinition
