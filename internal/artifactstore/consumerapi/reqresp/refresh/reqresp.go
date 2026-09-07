package refresh

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
)

// RefreshCollectionResult is the consumer-facing result of one Store-owned Collection refresh.
type RefreshCollectionResult struct {
	Catalog          catalog.Snapshot
	CreatedArtifacts []basespec.ArtifactID
	UpdatedArtifacts []basespec.ArtifactID
	Diagnostics      []providerapi.Diagnostic
	Candidates       int
}
