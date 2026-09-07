package collection

import (
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
)

type Draft struct {
	ID          basespec.CollectionID   `json:"id"`
	Kind        basespec.CollectionKind `json:"kind"`
	DisplayName string                  `json:"displayName"`
	Description string                  `json:"description,omitempty"`
	Enabled     bool                    `json:"enabled"`
	Data        json.RawMessage         `json:"data"`
}

type Update struct {
	ExpectedRevision uint64          `json:"expectedRevision"`
	DisplayName      string          `json:"displayName"`
	Description      string          `json:"description,omitempty"`
	Enabled          bool            `json:"enabled"`
	Data             json.RawMessage `json:"data"`
}
