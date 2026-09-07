package artifact

import (
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
)

// AdoptRequest identifies one currently catalogued valid occurrence that a
// caller explicitly wants to represent as an observed Artifact.
type AdoptRequest struct {
	ArtifactID              basespec.ArtifactID
	Collection              collection.CollectionRef
	Occurrence              catalog.OccurrenceKey
	ExpectedCatalogRevision uint64
	Name                    string
	Enabled                 bool
	Data                    json.RawMessage
}

// PinRequest creates an Artifact binding that remains represented even when
// its source occurrence is missing or temporarily invalid.
type PinRequest struct {
	ArtifactID                 basespec.ArtifactID
	Collection                 collection.CollectionRef
	ExpectedCollectionRevision uint64
	Binding                    artifact.SourceBinding
	Name                       string
	Enabled                    bool
	Data                       json.RawMessage
}

// SuppressRequest records that automatic reconciliation must not create an
// observed Artifact for one attached source binding.
type SuppressRequest struct {
	Collection                 collection.CollectionRef
	ExpectedCollectionRevision uint64
	Binding                    artifact.SourceBinding
}
