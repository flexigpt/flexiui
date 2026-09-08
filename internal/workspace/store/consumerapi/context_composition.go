package consumerapi

import (
	"errors"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	workspaceRuntime "github.com/flexigpt/flexigpt-app/internal/workspace/runtime"
)

type OverflowBehavior = workspaceRuntime.OverflowBehavior

const (
	OverflowTruncate = workspaceRuntime.OverflowTruncate
	OverflowExclude  = workspaceRuntime.OverflowExclude
)

const (
	DiagnosticCodeContextDocumentTruncated = workspaceRuntime.DiagnosticCodeContextDocumentTruncated
	DiagnosticCodeContextDocumentExcluded  = workspaceRuntime.DiagnosticCodeContextDocumentExcluded
	DiagnosticCodeContextBudgetExceeded    = workspaceRuntime.DiagnosticCodeContextBudgetExceeded
)

type CompositionPolicy = workspaceRuntime.CompositionPolicy

func DefaultCompositionPolicy() CompositionPolicy {
	return workspaceRuntime.DefaultCompositionPolicy()
}

type CompositionStatus = workspaceRuntime.CompositionStatus

const (
	CompositionIncluded    = workspaceRuntime.CompositionIncluded
	CompositionTruncated   = workspaceRuntime.CompositionTruncated
	CompositionExcluded    = workspaceRuntime.CompositionExcluded
	CompositionDenied      = workspaceRuntime.CompositionDenied
	CompositionUnavailable = workspaceRuntime.CompositionUnavailable
)

type CompositionDecision struct {
	Artifact      artifact.ArtifactRef `json:"artifact"`
	Status        CompositionStatus    `json:"status"`
	Code          string               `json:"code,omitempty"`
	OriginalBytes int                  `json:"originalBytes"`
	IncludedBytes int                  `json:"includedBytes"`
}

func applyCompositionPolicy(
	engine *workspaceRuntime.Engine,
	policy CompositionPolicy,
	values []ContextContribution,
	diagnostics []diagnostic.Diagnostic,
	decisions []CompositionDecision,
) (
	[]ContextContribution,
	string,
	[]diagnostic.Diagnostic,
	[]CompositionDecision,
	error,
) {
	if engine == nil {
		return nil, "", nil, nil, errors.New(
			"workspace context runtime engine is nil",
		)
	}

	runtimeValues := make(
		[]workspaceRuntime.ContextContribution,
		0,
		len(values),
	)
	byID := make(map[string]ContextContribution, len(values))

	for _, value := range values {
		id := contextContributionID(value.Artifact)
		byID[id] = value
		runtimeValues = append(
			runtimeValues,
			workspaceRuntime.ContextContribution{
				ID:              id,
				Name:            value.Name,
				Role:            string(value.Role),
				Locator:         string(value.Locator),
				Content:         value.Content,
				ConventionOrder: value.ConventionOrder,
			},
		)
	}

	result, err := engine.Compose(policy, runtimeValues)
	if err != nil {
		return nil, "", nil, nil, err
	}

	included := make([]ContextContribution, 0, len(result.Contributions))
	for _, value := range result.Contributions {
		original, found := byID[value.ID]
		if !found {
			return nil, "", nil, nil, fmt.Errorf(
				"workspace context runtime returned an unknown contribution %q",
				value.ID,
			)
		}
		original.Content = value.Content
		original.OriginalBytes = value.OriginalBytes
		original.IncludedBytes = value.IncludedBytes
		original.Truncated = value.Truncated
		included = append(included, original)
	}

	for _, value := range result.Diagnostics {
		original, found := byID[value.ID]
		if !found {
			return nil, "", nil, nil, fmt.Errorf(
				"workspace context runtime returned an unknown diagnostic contribution %q",
				value.ID,
			)
		}
		diagnostics = diagnostic.Append(
			diagnostics,
			compositionDiagnostic(
				original,
				value.Code,
				value.Message,
			),
		)
	}

	for _, value := range result.Decisions {
		original, found := byID[value.ID]
		if !found {
			return nil, "", nil, nil, fmt.Errorf(
				"workspace context runtime returned an unknown decision contribution %q",
				value.ID,
			)
		}
		decisions = append(decisions, CompositionDecision{
			Artifact:      original.Artifact,
			Status:        value.Status,
			Code:          value.Code,
			OriginalBytes: value.OriginalBytes,
			IncludedBytes: value.IncludedBytes,
		})
	}

	return included, result.Prompt, diagnostics, decisions, nil
}

func compositionDiagnostic(
	value ContextContribution,
	code string,
	message string,
) diagnostic.Diagnostic {
	return diagnostic.Diagnostic{
		Severity: diagnostic.SeverityWarning,
		Code:     code,
		Message:  message,
		Location: &diagnostic.Location{
			Locator: value.Locator,
		},
	}
}

func contextContributionID(value artifact.ArtifactRef) string {
	return string(value.RootID) + "\x00" + string(value.ArtifactID)
}
