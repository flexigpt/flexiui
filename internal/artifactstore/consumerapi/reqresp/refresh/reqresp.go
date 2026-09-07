package refresh

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
)

// RefreshCollectionResult is the consumer-facing result of one Store-owned Collection refresh.
type RefreshCollectionResult struct {
	Catalog          catalog.Snapshot
	CreatedArtifacts []artifact.ArtifactID
	UpdatedArtifacts []artifact.ArtifactID
	Diagnostics      []diagnostic.Diagnostic
	Candidates       int
}
