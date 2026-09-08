package bundle

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
)

type CreateSkillBundleRequest struct {
	Body *CreateBundleRequest `json:"body"`
}

type CreateSkillBundleResponse struct {
	Body *Bundle `json:"body"`
}

type GetSkillBundleRequest struct {
	Bundle collection.CollectionRef `json:"bundle"`
}

type GetSkillBundleResponse struct {
	Body *Bundle `json:"body"`
}

type ListSkillBundlesRequest struct {
	RootID root.RootID `json:"rootID"`
}

type ListSkillBundlesResponseBody struct {
	Bundles []Bundle `json:"bundles"`
}

type ListSkillBundlesResponse struct {
	Body *ListSkillBundlesResponseBody `json:"body"`
}

type UpdateSkillBundleRequest struct {
	Body *UpdateBundleRequest `json:"body"`
}

type UpdateSkillBundleResponse struct {
	Body *Bundle `json:"body"`
}

type RetireSkillBundleRequest struct {
	Bundle           collection.CollectionRef `json:"bundle"`
	ExpectedRevision uint64                   `json:"expectedRevision"`
}

type RetireSkillBundleResponse struct {
	Body *collection.Collection `json:"body"`
}

type PurgeSkillBundleRequest struct {
	Bundle           collection.CollectionRef `json:"bundle"`
	ExpectedRevision uint64                   `json:"expectedRevision"`
}

type PurgeSkillBundleResponse struct {
	Bundle collection.CollectionRef `json:"bundle"`
}

type AttachSkillBundleSourceRequestBody struct {
	Bundle                     collection.CollectionRef `json:"bundle"`
	ExpectedCollectionRevision uint64                   `json:"expectedCollectionRevision"`
	Attachment                 AttachmentDraft          `json:"attachment"`
}

type AttachSkillBundleSourceRequest struct {
	Body *AttachSkillBundleSourceRequestBody `json:"body"`
}

type AttachSkillBundleSourceResponse struct {
	Body *Bundle `json:"body"`
}

type RefreshSkillBundleRequest struct {
	Bundle collection.CollectionRef `json:"bundle"`
}

type RefreshSkillBundleResponse struct {
	Body *catalog.RefreshCollectionResult `json:"body"`
}

type CreateManagedSkillStoreRequest struct {
	Body *CreateManagedSkillRequest `json:"body"`
}

type CreateManagedSkillStoreResponse struct {
	Body *CreateManagedSkillResponse `json:"body"`
}

type GetManagedSkillDocumentRequest struct {
	Artifact artifact.ArtifactRef `json:"artifact"`
}

type GetManagedSkillDocumentResponse struct {
	Body *ManagedSkillDocument `json:"body"`
}

type AdoptSkillStoreRequest struct {
	Body *AdoptSkillRequest `json:"body"`
}

type AdoptSkillStoreResponse struct {
	Body *artifact.Artifact `json:"body"`
}

type PinSkillStoreRequest struct {
	Body *PinSkillRequest `json:"body"`
}

type PinSkillStoreResponse struct {
	Body *artifact.Artifact `json:"body"`
}

type GetSkillRequest struct {
	Artifact artifact.ArtifactRef `json:"artifact"`
}

type GetSkillResponse struct {
	Body *artifact.Artifact `json:"body"`
}

type ListBundleSkillsRequest struct {
	Bundle collection.CollectionRef `json:"bundle"`
}

type ListBundleSkillsResponseBody struct {
	Skills []artifact.Artifact `json:"skills"`
}

type ListBundleSkillsResponse struct {
	Body *ListBundleSkillsResponseBody `json:"body"`
}

type SetSkillEnabledRequestBody struct {
	Artifact         artifact.ArtifactRef `json:"artifact"`
	ExpectedRevision uint64               `json:"expectedRevision"`
	Enabled          bool                 `json:"enabled"`
}

type SetSkillEnabledRequest struct {
	Body *SetSkillEnabledRequestBody `json:"body"`
}

type SetSkillEnabledResponse struct {
	Body *artifact.Artifact `json:"body"`
}

type UnadoptSkillRequest struct {
	Artifact         artifact.ArtifactRef `json:"artifact"`
	ExpectedRevision uint64               `json:"expectedRevision"`
	Suppress         bool                 `json:"suppress"`
}

type UnadoptSkillResponse struct {
	Artifact artifact.ArtifactRef `json:"artifact"`
}

type PurgeSkillRequest struct {
	Artifact         artifact.ArtifactRef `json:"artifact"`
	ExpectedRevision uint64               `json:"expectedRevision"`
}

type PurgeSkillResponse struct {
	Artifact artifact.ArtifactRef `json:"artifact"`
}

// requireRequestBody performs only transport-shape validation.
//
// Domain semantic validation remains in bundle.API. Generic lifecycle,
// persistence, revision, and source/package validation remain in Artifact
// Store.
func requireRequestBody[T any](
	request *T,
	bodyPresent bool,
	requireBody bool,
	subject string,
) error {
	if request == nil {
		return fmt.Errorf(
			"%w: %s request is required",
			basespec.ErrInvalid,
			subject,
		)
	}
	if requireBody && !bodyPresent {
		return fmt.Errorf(
			"%w: %s request body is required",
			basespec.ErrInvalid,
			subject,
		)
	}
	return nil
}

func wrapStoreError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("skill bundle %s: %w", operation, err)
}
