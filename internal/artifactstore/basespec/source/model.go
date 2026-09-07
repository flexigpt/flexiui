package source

import (
	"encoding/json"
	"time"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
)

type Draft struct {
	ID          basespec.SourceID   `json:"id"`
	StorageKey  basespec.StorageKey `json:"storageKey"`
	Kind        basespec.SourceKind `json:"kind"`
	DisplayName string              `json:"displayName"`
	Enabled     bool                `json:"enabled"`
	Config      json.RawMessage     `json:"config"`
}

type Update struct {
	ExpectedRevision uint64
	DisplayName      string
	Enabled          bool

	// Config is write-only replacement configuration. A nil value preserves
	// the current normalized configuration so public callers can update Source
	// metadata without reading or resending private Source configuration.
	Config json.RawMessage
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}
