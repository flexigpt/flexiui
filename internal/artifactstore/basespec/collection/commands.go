package collection

import (
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

type Draft struct {
	ID          CollectionID    `json:"id"`
	Kind        CollectionKind  `json:"kind"`
	DisplayName string          `json:"displayName"`
	Description string          `json:"description,omitempty"`
	Enabled     bool            `json:"enabled"`
	Data        json.RawMessage `json:"data"`
}

type Update struct {
	ExpectedRevision uint64          `json:"expectedRevision"`
	DisplayName      string          `json:"displayName"`
	Description      string          `json:"description,omitempty"`
	Enabled          bool            `json:"enabled"`
	Data             json.RawMessage `json:"data"`
}

type PublishCollectionRequest struct {
	Collection     CollectionRef                    `json:"collection"`
	SourceID       source.SourceID                  `json:"sourceID"`
	Package        source.ManagedPackagePublication `json:"package"`
	AllowProtected bool                             `json:"allowProtected"`
	ForceRefresh   bool                             `json:"forceRefresh"`
}

type PublishCollectionResult struct {
	Source     source.Summary `json:"source"`
	Generation string         `json:"generation"`
	Refreshed  bool           `json:"refreshed"`
}
