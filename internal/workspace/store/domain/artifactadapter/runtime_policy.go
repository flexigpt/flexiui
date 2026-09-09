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

type RuntimePolicyRequest struct {
	Use              workspaceRuntime.RuntimeUse
	Workspace        workspaceDomain.Workspace
	Artifact         artifact.Artifact
	DefinitionDigest cryptoutil.Digest
	SourceID         source.SourceID

	RuntimeDisabled        bool
	RuntimeSettingsInvalid bool
}

type SourceUsePolicy interface {
	Decide(
		ctx context.Context,
		request RuntimePolicyRequest,
	) workspaceRuntime.RuntimeDecision
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
) workspaceRuntime.RuntimeDecision {
	if err := ctx.Err(); err != nil {
		return workspaceRuntime.RuntimeDecision{
			Disposition: workspaceRuntime.RuntimeUnavailable,
			Code:        workspaceDomain.DiagnosticCodeRuntimeUnavailable,
			Message:     "runtime policy evaluation was cancelled",
		}
	}
	if !request.Workspace.Collection.Enabled {
		return workspaceRuntime.RuntimeDecision{
			Disposition: workspaceRuntime.RuntimeUnavailable,
			Code:        workspaceDomain.DiagnosticCodeRuntimeUnavailable,
			Message:     "the Workspace is disabled",
		}
	}
	if !request.Artifact.Enabled ||
		request.Artifact.State != artifact.StateAvailable {
		return workspaceRuntime.RuntimeDecision{
			Disposition: workspaceRuntime.RuntimeUnavailable,
			Code:        workspaceDomain.DiagnosticCodeRuntimeUnavailable,
			Message:     "the Workspace Artifact is not enabled and available",
		}
	}
	if request.RuntimeSettingsInvalid {
		return workspaceRuntime.RuntimeDecision{
			Disposition: workspaceRuntime.RuntimeUnavailable,
			Code:        workspaceDomain.DiagnosticCodeRuntimeUnavailable,
			Message:     "the Workspace Artifact has invalid local runtime policy data",
		}
	}
	if request.RuntimeDisabled {
		return workspaceRuntime.RuntimeDecision{
			Disposition: workspaceRuntime.RuntimeDenied,
			Code:        workspaceDomain.DiagnosticCodeRuntimeDenied,
			Message:     "runtime use is disabled for this Workspace Artifact",
		}
	}
	return workspaceRuntime.RuntimeDecision{
		Disposition: workspaceRuntime.RuntimeAllowed,
	}
}

func RuntimeDecisionDiagnostic(
	decision workspaceRuntime.RuntimeDecision,
	value artifact.Artifact,
) diagnostic.Diagnostic {
	severity := diagnostic.SeverityWarning
	if decision.Disposition == workspaceRuntime.RuntimeUnavailable {
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
