package lifecycle

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

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

type CreateArtifactSourceRequest struct {
	RootID root.RootID `path:"rootID" required:"true"`
	Body   *source.Draft
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

type UpdateArtifactSourceRequest struct {
	RootID   root.RootID     `path:"rootID"   required:"true"`
	SourceID source.SourceID `path:"sourceID" required:"true"`
	Body     *source.Update
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
