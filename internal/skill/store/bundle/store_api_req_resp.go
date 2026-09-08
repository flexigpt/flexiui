package bundle

import (
	"fmt"

	"github.com/flexigpt/agentskills-go/document"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

type CreateSkillBundleBody struct {
	RootID                  root.RootID
	CollectionID            collection.CollectionID
	ManagedSourceID         source.SourceID
	ManagedSourceStorageKey basespec.StorageKey
	DisplayName             string
	Description             string
	Enabled                 bool
	LogicalName             basespec.LogicalName
	LogicalVersion          basespec.LogicalVersion
	Labels                  map[string]string
	Attachments             []AttachmentDraft
}

type CreateSkillBundleRequest struct {
	Body *CreateSkillBundleBody `json:"body"`
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

type UpdateSkillBundleBody struct {
	Bundle           collection.CollectionRef
	ExpectedRevision uint64
	DisplayName      string
	Description      string
	Enabled          bool
}

type UpdateSkillBundleRequest struct {
	Body *UpdateSkillBundleBody `json:"body"`
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

type AttachSkillBundleSourceBody struct {
	Bundle                     collection.CollectionRef
	ExpectedCollectionRevision uint64
	Attachment                 AttachmentDraft
}

type AttachSkillBundleSourceRequest struct {
	Body *AttachSkillBundleSourceBody `json:"body"`
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

type CreateManagedSkillBody struct {
	Bundle                     collection.CollectionRef
	ExpectedCollectionRevision uint64
	ArtifactID                 artifact.ArtifactID
	SkillName                  string
	SKILLMD                    []byte
	ExpectedArtifactRevision   uint64
	Document                   *document.SkillDocument
	Files                      []source.ManagedPackageFile
	Enabled                    bool
}

type CreateManagedSkillRequest struct {
	Body *CreateManagedSkillBody `json:"body"`
}

type CreateManagedSkillResponse struct {
	Artifact artifact.Artifact
	Address  artifact.ArtifactAddress
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

type AdoptSkillBody struct {
	Bundle                  collection.CollectionRef
	Occurrence              catalog.OccurrenceKey
	ArtifactID              artifact.ArtifactID
	ExpectedCatalogRevision uint64
	Name                    string
	Enabled                 bool
}

type AdoptSkillRequest struct {
	Body *AdoptSkillBody `json:"body"`
}

type AdoptSkillResponse struct {
	Body *artifact.Artifact `json:"body"`
}

type PinSkillBody struct {
	Bundle                     collection.CollectionRef
	ExpectedCollectionRevision uint64
	ArtifactID                 artifact.ArtifactID
	Binding                    artifact.SourceBinding
	Name                       string
	Enabled                    bool
}

type PinSkillRequest struct {
	Body *PinSkillBody `json:"body"`
}

type PinSkillResponse struct {
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

type SetSkillEnabledBody struct {
	Artifact         artifact.ArtifactRef
	ExpectedRevision uint64
	Enabled          bool
}

type SetSkillEnabledRequest struct {
	Body *SetSkillEnabledBody `json:"body"`
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

func requireStoreRequest[T any](
	request *T,
	requireBody bool,
	bodyPresent bool,
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
