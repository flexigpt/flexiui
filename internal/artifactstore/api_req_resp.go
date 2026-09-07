package artifactstore

import (
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
)

// Config contains the application-composition inputs required to open one
// Artifact Store.
//
// Store implementation dependencies, source snapshots, source adapters,
// metadata repositories, SQLite handles, clocks, and automatic Artifact ID
// providers remain private to Artifact Store.
type Config struct {
	BaseDirectory string

	ArtifactProviders []providerapi.Provider

	ProtectedRoots []root.RootID
	RetainedRoots  []root.RootID
}

type CreateArtifactRootRequest struct {
	Body *root.RootDraft
}

type CreateArtifactRootResponse struct {
	Body *root.Root
}

type GetArtifactRootRequest struct {
	RootID root.RootID `path:"rootID" required:"true"`
}

type GetArtifactRootResponse struct {
	Body *root.Root
}

type ListArtifactRootsRequest struct{}

type ListArtifactRootsResponseBody struct {
	Roots []root.Root `json:"roots"`
}

type ListArtifactRootsResponse struct {
	Body *ListArtifactRootsResponseBody
}

type UpdateArtifactRootRequest struct {
	RootID root.RootID `path:"rootID" required:"true"`
	Body   *root.RootUpdate
}

type UpdateArtifactRootResponse struct {
	Body *root.Root
}

type RetireArtifactRootRequest struct {
	RootID           root.RootID `path:"rootID" required:"true"`
	ExpectedRevision uint64      `              required:"true" json:"expectedRevision"`
}

type RetireArtifactRootResponse struct {
	Body *root.Root
}

type PurgeArtifactRootRequest struct {
	RootID           root.RootID `path:"rootID" required:"true"`
	ExpectedRevision uint64      `              required:"true" json:"expectedRevision"`
}

type PurgeArtifactRootResponse struct {
	RootID root.RootID `json:"rootID"`
}

// ArtifactSourceDraft is write-only. Source configuration can contain local
// filesystem paths or provider credentials and is not returned by the API.
type ArtifactSourceDraft struct {
	ID          source.SourceID     `json:"id"          required:"true"`
	StorageKey  basespec.StorageKey `json:"storageKey"  required:"true"`
	Kind        source.SourceKind   `json:"kind"        required:"true"`
	DisplayName string              `json:"displayName" required:"true"`
	Enabled     bool                `json:"enabled"`
	Config      json.RawMessage     `json:"config"`
}

type CreateArtifactSourceRequest struct {
	RootID root.RootID `path:"rootID" required:"true"`
	Body   *ArtifactSourceDraft
}

type CreateArtifactSourceResponse struct {
	Body *source.Summary
}

type GetArtifactSourceRequest struct {
	RootID   root.RootID     `path:"rootID"   required:"true"`
	SourceID source.SourceID `path:"sourceID" required:"true"`
}

type GetArtifactSourceResponse struct {
	Body *source.Summary
}

type ListArtifactSourcesRequest struct {
	RootID root.RootID `path:"rootID" required:"true"`
}

type ListArtifactSourcesResponseBody struct {
	Sources []source.Summary `json:"sources"`
}

type ListArtifactSourcesResponse struct {
	Body *ListArtifactSourcesResponseBody
}

type UpdateArtifactSourceRequestBody struct {
	ExpectedRevision uint64          `json:"expectedRevision" required:"true"`
	DisplayName      string          `json:"displayName"      required:"true"`
	Enabled          bool            `json:"enabled"`
	Config           json.RawMessage `json:"config,omitempty"`
}

type UpdateArtifactSourceRequest struct {
	RootID   root.RootID     `path:"rootID"   required:"true"`
	SourceID source.SourceID `path:"sourceID" required:"true"`
	Body     *UpdateArtifactSourceRequestBody
}

type UpdateArtifactSourceResponse struct {
	Body *source.Summary
}

type RetireArtifactSourceRequest struct {
	RootID           root.RootID     `path:"rootID"   required:"true"`
	SourceID         source.SourceID `path:"sourceID" required:"true"`
	ExpectedRevision uint64          `                required:"true" json:"expectedRevision"`
}

type RetireArtifactSourceResponse struct {
	Body *source.Summary
}

type PurgeArtifactSourceRequest struct {
	RootID           root.RootID     `path:"rootID"   required:"true"`
	SourceID         source.SourceID `path:"sourceID" required:"true"`
	ExpectedRevision uint64          `                required:"true" json:"expectedRevision"`
}

type PurgeArtifactSourceResponse struct {
	RootID   root.RootID     `json:"rootID"`
	SourceID source.SourceID `json:"sourceID"`
}

type ListArtifactSourceKindsRequest struct{}

type ListArtifactSourceKindsResponseBody struct {
	Kinds []source.SourceKind `json:"kinds"`
}

type ListArtifactSourceKindsResponse struct {
	Body *ListArtifactSourceKindsResponseBody
}
