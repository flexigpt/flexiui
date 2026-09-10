package domain

import (
	"fmt"
	"strings"

	"github.com/flexigpt/agentskills-go/document"
	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/definition"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
)

func ValidateDefinition(value definition.Definition) error {
	_, err := DocumentFromDefinition(value)
	return err
}

// DecodeSkillDocument is the shared SKILL.md parse-and-definition path used
// by discovery, managed Skill publication, and built-in Skill hydration.
func DecodeSkillDocument(
	content []byte,
	expectedName string,
) (definition.Definition, []diagnostic.Diagnostic, error) {
	doc, warnings, err := document.ParseSkillDocument(
		content,
		document.ParseSkillDocumentOptions{
			ExpectedName: expectedName,
		},
	)
	if err != nil {
		return definition.Definition{}, nil, err
	}

	value, err := definitionForDocument(doc)
	if err != nil {
		return definition.Definition{}, nil, err
	}
	canonical, err := definition.Canonicalize(value)
	if err != nil {
		return definition.Definition{}, nil, err
	}
	return canonical, warningDiagnostics(warnings), nil
}

// DocumentFromDefinition validates the canonical Artifact definition and
// returns the corresponding agentskills-go document. Consumers must use this
// projection instead of decoding the Skill body independently.
//
// Raw SKILL.md parsing remains exclusively in DecodeSkillDocument, which
// delegates to agentskills-go.ParseSkillDocument.
func DocumentFromDefinition(
	value definition.Definition,
) (document.SkillDocument, error) {
	if value.Kind != artifactbuiltin.AgentSkillArtifactKind {
		return document.SkillDocument{}, fmt.Errorf(
			"%w: Skill definition kind must be %q",
			basespec.ErrInvalid,
			artifactbuiltin.AgentSkillArtifactKind,
		)
	}
	if value.SchemaID != artifactbuiltin.AgentSkillSchemaID {
		return document.SkillDocument{}, fmt.Errorf(
			"%w: Skill definition schema must be %q",
			basespec.ErrInvalid,
			artifactbuiltin.AgentSkillSchemaID,
		)
	}
	if value.SchemaVersion != artifactbuiltin.AgentSkillSchemaVersion {
		return document.SkillDocument{}, fmt.Errorf(
			"%w: Skill definition schema version must be %q",
			basespec.ErrInvalid,
			artifactbuiltin.AgentSkillSchemaVersion,
		)
	}
	if len(value.Dependencies) != 0 {
		return document.SkillDocument{}, fmt.Errorf(
			"%w: Agent Skills cannot declare generic portable dependencies",
			basespec.ErrInvalid,
		)
	}
	doc, err := definition.DecodeBody[document.SkillDocument](value.Body)
	if err != nil {
		return document.SkillDocument{}, err
	}

	if err := document.ValidateSkillDocument(doc); err != nil {
		return document.SkillDocument{}, fmt.Errorf(
			"%w: invalid Agent Skill document: %w",
			basespec.ErrInvalid,
			err,
		)
	}
	if string(value.LogicalName) != doc.Name {
		return document.SkillDocument{}, fmt.Errorf(
			"%w: Skill logical name does not match body.name",
			basespec.ErrInvalid,
		)
	}
	if value.DisplayName != doc.DisplayName {
		return document.SkillDocument{}, fmt.Errorf(
			"%w: Skill display name does not match body.displayName",
			basespec.ErrInvalid,
		)
	}
	if value.Description != doc.Description {
		return document.SkillDocument{}, fmt.Errorf(
			"%w: Skill description does not match body.description",
			basespec.ErrInvalid,
		)
	}
	if value.Labels[artifactbuiltin.AgentSkillInsertLabelKey] != string(doc.Insert) {
		return document.SkillDocument{}, fmt.Errorf(
			"%w: Skill insert label does not match body.insert",
			basespec.ErrInvalid,
		)
	}
	return doc, nil
}

func definitionForDocument(
	doc document.SkillDocument,
) (definition.Definition, error) {
	raw, err := definition.EncodeBody(doc)
	if err != nil {
		return definition.Definition{}, err
	}

	return definition.Definition{
		Kind:          artifactbuiltin.AgentSkillArtifactKind,
		SchemaID:      artifactbuiltin.AgentSkillSchemaID,
		SchemaVersion: artifactbuiltin.AgentSkillSchemaVersion,
		LogicalName:   basespec.LogicalName(doc.Name),
		DisplayName:   doc.DisplayName,
		Description:   doc.Description,
		Labels: map[string]string{
			artifactbuiltin.AgentSkillInsertLabelKey: string(doc.Insert),
		},
		Body: raw,
	}, nil
}

func warningDiagnostics(
	warnings []string,
) []diagnostic.Diagnostic {
	output := make([]diagnostic.Diagnostic, 0, len(warnings))
	for _, warning := range warnings {
		if len(output) == diagnostic.MaxDiagnostics {
			break
		}
		warning = strings.TrimSpace(warning)
		if warning == "" {
			continue
		}
		output = append(output, diagnostic.Diagnostic{
			Severity: diagnostic.SeverityWarning,
			Code:     "agent.skill.parse-warning",
			Message:  diagnostic.BoundedMessage(warning),
		})
	}
	return output
}
