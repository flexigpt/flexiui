package collection

import (
	"encoding/json"
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
