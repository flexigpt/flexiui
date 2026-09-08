package consumerapi

import (
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
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

type CreateArtifactCollectionRequestBody struct {
	Draft       collection.Draft             `json:"draft"`
	Attachments []collection.AttachmentDraft `json:"attachments"`
}

type CreateArtifactCollectionRequest struct {
	RootID root.RootID `path:"rootID" required:"true"`
	Body   *CreateArtifactCollectionRequestBody
}

type CreateArtifactCollectionResponseBody struct {
	Collection  collection.Collection   `json:"collection"`
	Attachments []collection.Attachment `json:"attachments"`
}

type CreateArtifactCollectionResponse struct {
	Body *CreateArtifactCollectionResponseBody
}

type GetArtifactCollectionRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
}

type GetArtifactCollectionResponse struct {
	Body *collection.Collection
}

type ListArtifactCollectionsRequest struct {
	RootID root.RootID `path:"rootID" required:"true"`
}

type ListArtifactCollectionsResponseBody struct {
	Collections []collection.Collection `json:"collections"`
}

type ListArtifactCollectionsResponse struct {
	Body *ListArtifactCollectionsResponseBody
}

type UpdateArtifactCollectionRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
	Body       *collection.Update
}

type UpdateArtifactCollectionResponse struct {
	Body *collection.Collection
}

type RetireArtifactCollectionRequest struct {
	Collection       collection.CollectionRef `json:"collection"       required:"true"`
	ExpectedRevision uint64                   `json:"expectedRevision" required:"true"`
}

type RetireArtifactCollectionResponse struct {
	Body *collection.Collection
}

type PurgeArtifactCollectionRequest struct {
	Collection       collection.CollectionRef `json:"collection"       required:"true"`
	ExpectedRevision uint64                   `json:"expectedRevision" required:"true"`
}

type PurgeArtifactCollectionResponse struct {
	Collection collection.CollectionRef `json:"collection"`
}

type AttachArtifactCollectionSourceRequestBody struct {
	ExpectedCollectionRevision uint64                     `json:"expectedCollectionRevision"`
	Attachment                 collection.AttachmentDraft `json:"attachment"`
}

type AttachArtifactCollectionSourceRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
	Body       *AttachArtifactCollectionSourceRequestBody
}

type ArtifactCollectionAttachmentResponseBody struct {
	Collection collection.Collection `json:"collection"`
	Attachment collection.Attachment `json:"attachment"`
}

type AttachArtifactCollectionSourceResponse struct {
	Body *ArtifactCollectionAttachmentResponseBody
}

type GetArtifactCollectionAttachmentRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
	SourceID   source.SourceID          `json:"sourceID"   required:"true"`
}

type GetArtifactCollectionAttachmentResponse struct {
	Body *collection.Attachment
}

type ListArtifactCollectionAttachmentsRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
}

type ListArtifactCollectionAttachmentsResponseBody struct {
	Attachments []collection.Attachment `json:"attachments"`
}

type ListArtifactCollectionAttachmentsResponse struct {
	Body *ListArtifactCollectionAttachmentsResponseBody
}

type UpdateArtifactCollectionAttachmentRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
	SourceID   source.SourceID          `json:"sourceID"   required:"true"`
	Body       *collection.AttachmentUpdate
}

type UpdateArtifactCollectionAttachmentResponse struct {
	Body *ArtifactCollectionAttachmentResponseBody
}

type DetachArtifactCollectionSourceRequest struct {
	Collection                 collection.CollectionRef `json:"collection"                 required:"true"`
	SourceID                   source.SourceID          `json:"sourceID"                   required:"true"`
	ExpectedCollectionRevision uint64                   `json:"expectedCollectionRevision"`
	ExpectedAttachmentRevision uint64                   `json:"expectedAttachmentRevision"`
}

type DetachArtifactCollectionSourceResponse struct {
	Body *collection.Collection
}

type ReplaceArtifactCollectionAttachmentRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
	Body       *collection.AttachmentReplacement
}

type ReplaceArtifactCollectionAttachmentResponse struct {
	Body *ArtifactCollectionAttachmentResponseBody
}

type GetArtifactRecordRequest struct {
	Artifact artifact.ArtifactRef `json:"artifact" required:"true"`
}

type GetArtifactRecordResponse struct {
	Body *artifact.Artifact
}

type ListArtifactCollectionRecordsRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
}

type ListArtifactCollectionRecordsResponseBody struct {
	Artifacts []artifact.Artifact `json:"artifacts"`
}

type ListArtifactCollectionRecordsResponse struct {
	Body *ListArtifactCollectionRecordsResponseBody
}

type AdoptArtifactRecordRequest struct {
	Body *catalog.AdoptRequest
}

type AdoptArtifactRecordResponse struct {
	Body *artifact.Artifact
}

type PinArtifactRecordRequest struct {
	Body *catalog.PinRequest
}

type PinArtifactRecordResponse struct {
	Body *artifact.Artifact
}

type SetArtifactRecordEnabledRequestBody struct {
	ExpectedRevision uint64 `json:"expectedRevision"`
	Enabled          bool   `json:"enabled"`
}

type SetArtifactRecordEnabledRequest struct {
	Artifact artifact.ArtifactRef `json:"artifact" required:"true"`
	Body     *SetArtifactRecordEnabledRequestBody
}

type SetArtifactRecordEnabledResponse struct {
	Body *artifact.Artifact
}

type SetArtifactRecordNameRequestBody struct {
	ExpectedRevision uint64 `json:"expectedRevision"`
	Name             string `json:"name"`
}

type SetArtifactRecordNameRequest struct {
	Artifact artifact.ArtifactRef `json:"artifact" required:"true"`
	Body     *SetArtifactRecordNameRequestBody
}

type SetArtifactRecordNameResponse struct {
	Body *artifact.Artifact
}

type UpdateArtifactRecordDataRequestBody struct {
	ExpectedRevision uint64          `json:"expectedRevision"`
	Data             json.RawMessage `json:"data"`
}

type UpdateArtifactRecordDataRequest struct {
	Artifact artifact.ArtifactRef `json:"artifact" required:"true"`
	Body     *UpdateArtifactRecordDataRequestBody
}

type UpdateArtifactRecordDataResponse struct {
	Body *artifact.Artifact
}

type UnadoptArtifactRecordRequest struct {
	Artifact         artifact.ArtifactRef `json:"artifact"         required:"true"`
	ExpectedRevision uint64               `json:"expectedRevision"`
	Suppress         bool                 `json:"suppress"`
}

type UnadoptArtifactRecordResponse struct {
	Artifact artifact.ArtifactRef `json:"artifact"`
}

type PurgeArtifactRecordRequest struct {
	Artifact         artifact.ArtifactRef `json:"artifact"         required:"true"`
	ExpectedRevision uint64               `json:"expectedRevision"`
}

type PurgeArtifactRecordResponse struct {
	Artifact artifact.ArtifactRef `json:"artifact"`
}

type PurgeAndSuppressArtifactRecordRequest struct {
	Artifact         artifact.ArtifactRef `json:"artifact"         required:"true"`
	ExpectedRevision uint64               `json:"expectedRevision"`
}

type PurgeAndSuppressArtifactRecordResponse struct {
	Artifact artifact.ArtifactRef `json:"artifact"`
}

type ListArtifactCollectionSuppressionsRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
}

type ListArtifactCollectionSuppressionsResponseBody struct {
	Suppressions []artifact.Suppression `json:"suppressions"`
}

type ListArtifactCollectionSuppressionsResponse struct {
	Body *ListArtifactCollectionSuppressionsResponseBody
}

type SuppressArtifactBindingRequest struct {
	Body *catalog.SuppressRequest
}

type SuppressArtifactBindingResponse struct {
	Body *artifact.Suppression
}

type UnsuppressArtifactBindingRequest struct {
	Collection       collection.CollectionRef `json:"collection"       required:"true"`
	Binding          artifact.SourceBinding   `json:"binding"          required:"true"`
	ExpectedRevision uint64                   `json:"expectedRevision"`
}

type UnsuppressArtifactBindingResponse struct {
	Collection collection.CollectionRef `json:"collection"`
	Binding    artifact.SourceBinding   `json:"binding"`
}

type RefreshArtifactCollectionRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
}

type RefreshArtifactCollectionResponse struct {
	Body *catalog.RefreshCollectionResult
}

type GetArtifactCollectionCatalogRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
}

type GetArtifactCollectionCatalogResponse struct {
	Body *catalog.Snapshot
}

type InspectArtifactCollectionCatalogRequest struct {
	Collection collection.CollectionRef `json:"collection" required:"true"`
}

type InspectArtifactCollectionCatalogResponse struct {
	Body *catalog.CatalogInspection
}
