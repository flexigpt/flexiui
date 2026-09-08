package providerapi

import (
	"context"
	"path"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
	skillDomain "github.com/flexigpt/flexigpt-app/internal/skill/store/domain"
)

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (*Decoder) ID() basespec.DecoderID {
	return artifactbuiltin.AgentSkillDecoderID
}

func (*Decoder) Revision() string {
	return artifactbuiltin.AgentSkillSchemaVersion
}

func (*Decoder) Recognize(
	_ context.Context,
	candidate providerapi.Candidate,
) providerapi.Recognition {
	if candidate.RequestsDecoder(artifactbuiltin.AgentSkillDecoderID) &&
		basespec.Locator(path.Base(string(candidate.Locator))) ==
			artifactbuiltin.AgentSkillDefinitionFileName {
		return providerapi.RecognitionPreferred
	}
	return providerapi.RecognitionNone
}

func (*Decoder) Decode(
	_ context.Context,
	candidate providerapi.Candidate,
) ([]providerapi.Decoded, []diagnostic.Diagnostic) {
	if !candidate.RequestsDecoder(artifactbuiltin.AgentSkillDecoderID) ||
		basespec.Locator(path.Base(string(candidate.Locator))) !=
			artifactbuiltin.AgentSkillDefinitionFileName {
		return nil, nil
	}

	parent := path.Dir(string(candidate.Locator))
	if parent == "." || parent == "/" || parent == "" {
		return nil, nil
	}

	value, warnings, err := skillDomain.DecodeSkillDocument(
		candidate.Content,
		path.Base(parent),
	)
	if err != nil {
		return nil, []diagnostic.Diagnostic{{
			Severity: diagnostic.SeverityError,
			Code:     "agent.skill.invalid",
			Message:  diagnostic.BoundedMessage(err.Error()),
			Location: &diagnostic.Location{
				Locator: candidate.Locator,
			},
		}}
	}

	for index := range warnings {
		warnings[index].Location = &diagnostic.Location{
			Locator: candidate.Locator,
		}
	}

	return []providerapi.Decoded{{Definition: value}}, warnings
}
