package support

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	skillDomain "github.com/flexigpt/flexigpt-app/internal/skill/store/domain"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
	workspaceDomainContext "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/context"
)

type defaultArtifactSupport struct {
	support workspaceDomain.ArtifactSupport
}

var defaultArtifactSupportMatrix = []defaultArtifactSupport{
	{
		support: workspaceDomainContext.ArtifactSupport(),
	},
	{
		support: workspaceDomain.ArtifactSupport{
			Kind:      artifactbuiltin.AgentSkillArtifactKind,
			SchemaID:  artifactbuiltin.AgentSkillSchemaID,
			DecoderID: artifactbuiltin.AgentSkillDecoderID,
			Validator: skillDomain.ValidateDefinition,
		},
	},
}

func DefaultArtifactSupports() []workspaceDomain.ArtifactSupport {
	output := make(
		[]workspaceDomain.ArtifactSupport,
		0,
		len(defaultArtifactSupportMatrix),
	)

	for _, value := range defaultArtifactSupportMatrix {
		output = append(output, value.support)
	}

	return output
}
