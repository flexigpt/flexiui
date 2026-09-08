package consumerapi

import "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/artifactadapter"

const (
	DiagnosticCodeArtifactInvalid           = artifactadapter.DiagnosticCodeArtifactInvalid
	DiagnosticCodeContextInvalidContent     = artifactadapter.DiagnosticCodeContextInvalidContent
	DiagnosticCodeContextInvalidUTF8        = artifactadapter.DiagnosticCodeContextInvalidUTF8
	DiagnosticCodeArtifactSchemaUnsupported = artifactadapter.DiagnosticCodeArtifactSchemaUnsupported
	DiagnosticCodeArtifactKindMismatch      = artifactadapter.DiagnosticCodeArtifactKindMismatch
	DiagnosticCodeProjectionInvalid         = artifactadapter.DiagnosticCodeProjectionInvalid
	DiagnosticCodeArtifactUnavailable       = artifactadapter.DiagnosticCodeArtifactUnavailable
	DiagnosticCodeArtifactUnresolved        = artifactadapter.DiagnosticCodeArtifactUnresolved
	DiagnosticCodeRuntimeDenied             = artifactadapter.DiagnosticCodeRuntimeDenied
	DiagnosticCodeRuntimeUnavailable        = artifactadapter.DiagnosticCodeRuntimeUnavailable
	DiagnosticCodeCatalogStale              = artifactadapter.DiagnosticCodeCatalogStale
	DiagnosticCodeCatalogDecoderStale       = artifactadapter.DiagnosticCodeCatalogDecoderStale
	DiagnosticCodeCatalogPlanStale          = artifactadapter.DiagnosticCodeCatalogPlanStale
)
