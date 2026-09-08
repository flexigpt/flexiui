package catalog

import (
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
)

// AdoptRequest identifies one current catalog occurrence that receives an explicit observed Artifact record.
type AdoptRequest struct {
	ArtifactID              artifact.ArtifactID      `json:"artifactID"`
	Collection              collection.CollectionRef `json:"collection"`
	Occurrence              OccurrenceKey            `json:"occurrence"`
	ExpectedCatalogRevision uint64                   `json:"expectedCatalogRevision"`
	Name                    string                   `json:"name"`
	Enabled                 bool                     `json:"enabled"`
	Data                    json.RawMessage          `json:"data"`
}

// PinRequest creates an Artifact record that remains represented while its
// source occurrence is temporarily absent or invalid.
type PinRequest struct {
	ArtifactID                 artifact.ArtifactID      `json:"artifactID"`
	Collection                 collection.CollectionRef `json:"collection"`
	ExpectedCollectionRevision uint64                   `json:"expectedCollectionRevision"`
	Binding                    artifact.SourceBinding   `json:"binding"`
	Name                       string                   `json:"name"`
	Enabled                    bool                     `json:"enabled"`
	Data                       json.RawMessage          `json:"data"`
}

// SuppressRequest records an explicit opt-out from automatic adoption.
type SuppressRequest struct {
	Collection                 collection.CollectionRef `json:"collection"`
	ExpectedCollectionRevision uint64                   `json:"expectedCollectionRevision"`
	Binding                    artifact.SourceBinding   `json:"binding"`
}
