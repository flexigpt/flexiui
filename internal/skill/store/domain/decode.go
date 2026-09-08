package domain

import (
	"strings"

	"github.com/flexigpt/agentskills-go/document"
	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/definition"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
)

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

func definitionForDocument(
	doc document.SkillDocument,
) (definition.Definition, error) {
	arguments := make([]Argument, 0, len(doc.Arguments))
	for _, argument := range doc.Arguments {
		arguments = append(arguments, Argument{
			Name:        argument.Name,
			Description: argument.Description,
			Default:     argument.Default,
		})
	}

	raw, err := definition.EncodeBody(Body{
		Name:           doc.Name,
		DisplayName:    doc.DisplayName,
		Description:    doc.Description,
		Insert:         string(doc.Insert),
		Arguments:      arguments,
		Tags:           append([]string(nil), doc.Tags...),
		MarkdownBody:   doc.MarkdownBody,
		RawFrontmatter: doc.RawFrontmatter,
	})
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
