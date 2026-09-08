package artifactadapter

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
	workspaceRuntime "github.com/flexigpt/flexigpt-app/internal/workspace/runtime"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
)

type RuntimeUse = workspaceRuntime.RuntimeUse

const (
	RuntimeUseContextPrompt = workspaceRuntime.RuntimeUseContextPrompt
	RuntimeUseSkill         = workspaceRuntime.RuntimeUseSkill
)

type RuntimeDisposition = workspaceRuntime.RuntimeDisposition

const (
	RuntimeAllowed     = workspaceRuntime.RuntimeAllowed
	RuntimeDenied      = workspaceRuntime.RuntimeDenied
	RuntimeUnavailable = workspaceRuntime.RuntimeUnavailable
)

type RuntimePolicyRequest struct {
	Use              RuntimeUse
	Workspace        workspaceDomain.Workspace
	Artifact         artifact.Artifact
	DefinitionDigest cryptoutil.Digest
	SourceID         source.SourceID
}

type RuntimeDecision = workspaceRuntime.RuntimeDecision

type SourceUsePolicy interface {
	Decide(
		ctx context.Context,
		request RuntimePolicyRequest,
	) RuntimeDecision
}

// ArtifactRuntimePolicy is the default local Workspace trust boundary.
//
// Discovery and management remain available. Runtime use is enabled by default
// unless the Artifact-local RuntimeDisabled flag explicitly disables it.
type ArtifactRuntimePolicy struct{}

func NewArtifactRuntimePolicy() *ArtifactRuntimePolicy {
	return &ArtifactRuntimePolicy{}
}

func (*ArtifactRuntimePolicy) Decide(
	ctx context.Context,
	request RuntimePolicyRequest,
) RuntimeDecision {
	if err := ctx.Err(); err != nil {
		return RuntimeDecision{
			Disposition: RuntimeUnavailable,
			Code:        DiagnosticCodeRuntimeUnavailable,
			Message:     "runtime policy evaluation was cancelled",
		}
	}
	if !request.Workspace.Collection.Enabled {
		return RuntimeDecision{
			Disposition: RuntimeUnavailable,
			Code:        DiagnosticCodeRuntimeUnavailable,
			Message:     "the Workspace is disabled",
		}
	}
	if !request.Artifact.Enabled ||
		request.Artifact.State != artifact.StateAvailable {
		return RuntimeDecision{
			Disposition: RuntimeUnavailable,
			Code:        DiagnosticCodeRuntimeUnavailable,
			Message:     "the Workspace Artifact is not enabled and available",
		}
	}
	disabled, err := ArtifactRuntimeDisabled(request.Artifact)
	if err != nil {
		return RuntimeDecision{
			Disposition: RuntimeUnavailable,
			Code:        DiagnosticCodeRuntimeUnavailable,
			Message:     "the Workspace Artifact has invalid local runtime policy data",
		}
	}
	if disabled {
		return RuntimeDecision{
			Disposition: RuntimeDenied,
			Code:        DiagnosticCodeRuntimeDenied,
			Message:     "runtime use is disabled for this Workspace Artifact",
		}
	}
	return RuntimeDecision{Disposition: RuntimeAllowed}
}

func RuntimeDecisionDiagnostic(
	decision RuntimeDecision,
	value artifact.Artifact,
) diagnostic.Diagnostic {
	severity := diagnostic.SeverityWarning
	if decision.Disposition == RuntimeUnavailable {
		severity = diagnostic.SeverityError
	}
	return diagnostic.Diagnostic{
		Severity: severity,
		Code:     decision.Code,
		Message:  decision.Message,
		Location: &diagnostic.Location{
			Locator:            value.Binding.Locator,
			SubresourceLocator: value.Binding.SubresourceLocator,
		},
	}
}
