package catalog

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
)

// RefreshCollectionResult is the Store-owned result of one Collection refresh.
type RefreshCollectionResult struct {
	Catalog          Snapshot                `json:"catalog"`
	CreatedArtifacts []artifact.ArtifactID   `json:"createdArtifacts"`
	UpdatedArtifacts []artifact.ArtifactID   `json:"updatedArtifacts"`
	Diagnostics      []diagnostic.Diagnostic `json:"diagnostics,omitempty"`
	Candidates       int                     `json:"candidates"`
}

// CatalogInspection reports Store-owned catalog freshness dimensions.
type CatalogInspection struct {
	Catalog         Snapshot `json:"catalog"`
	MetadataChanged bool     `json:"metadataChanged"`
	PlanChanged     bool     `json:"planChanged"`
	DecoderChanged  bool     `json:"decoderChanged"`
}

func (i CatalogInspection) IsCurrent() bool {
	return !i.MetadataChanged &&
		!i.PlanChanged &&
		!i.DecoderChanged
}
