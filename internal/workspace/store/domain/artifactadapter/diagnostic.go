package artifactadapter

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
)

const (
	DiagnosticCodeArtifactInvalid           = workspaceDomain.DiagnosticCodeArtifactInvalid
	DiagnosticCodeContextInvalidContent     = workspaceDomain.DiagnosticCodeContextInvalidContent
	DiagnosticCodeContextInvalidUTF8        = workspaceDomain.DiagnosticCodeContextInvalidUTF8
	DiagnosticCodeArtifactSchemaUnsupported = workspaceDomain.DiagnosticCodeArtifactSchemaUnsupported
	DiagnosticCodeArtifactKindMismatch      = workspaceDomain.DiagnosticCodeArtifactKindMismatch
	DiagnosticCodeProjectionInvalid         = workspaceDomain.DiagnosticCodeProjectionInvalid
	DiagnosticCodeArtifactUnavailable       = workspaceDomain.DiagnosticCodeArtifactUnavailable
	DiagnosticCodeArtifactUnresolved        = workspaceDomain.DiagnosticCodeArtifactUnresolved
	DiagnosticCodeRuntimeDenied             = workspaceDomain.DiagnosticCodeRuntimeDenied
	DiagnosticCodeRuntimeUnavailable        = workspaceDomain.DiagnosticCodeRuntimeUnavailable
	DiagnosticCodeCatalogStale              = workspaceDomain.DiagnosticCodeCatalogStale
	DiagnosticCodeCatalogDecoderStale       = workspaceDomain.DiagnosticCodeCatalogDecoderStale
	DiagnosticCodeCatalogPlanStale          = workspaceDomain.DiagnosticCodeCatalogPlanStale
)

func WorkspaceArtifactErrorDiagnostics(
	locator basespec.Locator,
	err error,
) []diagnostic.Diagnostic {
	return WorkspaceArtifactDiagnostics(
		locator,
		DiagnosticCodeArtifactInvalid,
		err.Error(),
	)
}

func WorkspaceArtifactDiagnostics(
	locator basespec.Locator,
	code string,
	message string,
) []diagnostic.Diagnostic {
	return []diagnostic.Diagnostic{{
		Severity: diagnostic.SeverityError,
		Code:     code,
		Message:  diagnostic.BoundedMessage(message),
		Location: &diagnostic.Location{
			Locator: locator,
		},
	}}
}
