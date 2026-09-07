package refresh

import "github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"

// CatalogInspection reports every Store-owned catalog freshness dimension.
type CatalogInspection struct {
	Catalog         catalog.Snapshot
	MetadataChanged bool
	PlanChanged     bool
	DecoderChanged  bool
}

func (i CatalogInspection) IsCurrent() bool {
	return !i.MetadataChanged &&
		!i.PlanChanged &&
		!i.DecoderChanged
}
