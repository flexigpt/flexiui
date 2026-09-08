package artifactadapter

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
)

func WorkspaceArtifactErrorDiagnostics(
	locator basespec.Locator,
	err error,
) []diagnostic.Diagnostic {
	return WorkspaceArtifactDiagnostics(
		locator,
		workspaceDomain.DiagnosticCodeArtifactInvalid,
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
