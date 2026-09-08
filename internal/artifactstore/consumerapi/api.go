package consumerapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	resourceimpl "github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/resource"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/system"
)

// API is the direct consumer-facing Artifact Store implementation.
//
// Composition constructs this value after provider registration. The command
// binding wrapper forwards request/response-shaped calls into this API.
type API struct {
	components *system.Components
	resources  *resourceimpl.Service

	closeOnce sync.Once
	closeErr  error
	closed    atomic.Bool
}

func New(components *system.Components) (*API, error) {
	if components == nil ||
		components.Roots == nil ||
		components.Sources == nil {
		return nil, errors.New("artifact store components are required")
	}
	if components.ArtifactReader == nil ||
		components.CollectionReader == nil ||
		components.Refresh == nil ||
		components.SourceRuntime == nil {
		return nil, errors.New(
			"artifact store resource components are required",
		)
	}

	resources, err := resourceimpl.NewService(
		components.ArtifactReader,
		components.CollectionReader,
		components.Refresh,
		components.SourceRuntime,
	)
	if err != nil {
		return nil, err
	}

	return &API{
		components: components,
		resources:  resources,
	}, nil
}

func (a *API) CreateRoot(
	ctx context.Context,
	draft root.RootDraft,
) (root.Root, error) {
	if err := a.check(ctx); err != nil {
		return root.Root{}, err
	}
	return a.components.Roots.Create(ctx, draft)
}

func (a *API) GetRoot(
	ctx context.Context,
	rootID root.RootID,
) (root.Root, error) {
	if err := a.check(ctx); err != nil {
		return root.Root{}, err
	}
	return a.components.Roots.Get(ctx, rootID)
}

func (a *API) ListRoots(
	ctx context.Context,
) ([]root.Root, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	return a.components.Roots.List(ctx)
}

func (a *API) UpdateRoot(
	ctx context.Context,
	rootID root.RootID,
	update root.RootUpdate,
) (root.Root, error) {
	if err := a.check(ctx); err != nil {
		return root.Root{}, err
	}
	return a.components.Roots.Update(ctx, rootID, update)
}

func (a *API) RetireRoot(
	ctx context.Context,
	rootID root.RootID,
	expectedRevision uint64,
) (root.Root, error) {
	if err := a.check(ctx); err != nil {
		return root.Root{}, err
	}
	return a.components.Roots.Retire(
		ctx,
		rootID,
		expectedRevision,
	)
}

func (a *API) PurgeRoot(
	ctx context.Context,
	rootID root.RootID,
	expectedRevision uint64,
) error {
	if err := a.check(ctx); err != nil {
		return err
	}
	return a.components.Roots.Purge(
		ctx,
		rootID,
		expectedRevision,
	)
}

func (a *API) CreateSource(
	ctx context.Context,
	rootID root.RootID,
	draft source.Draft,
) (source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, err
	}
	draft.Config = append(json.RawMessage(nil), draft.Config...)
	return a.components.Sources.Create(ctx, rootID, draft)
}

func (a *API) CreateSourceWithStatus(
	ctx context.Context,
	rootID root.RootID,
	draft source.Draft,
) (source.Summary, bool, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, false, err
	}
	draft.Config = append(json.RawMessage(nil), draft.Config...)
	return a.components.Sources.CreateWithStatus(ctx, rootID, draft)
}

func (a *API) DiscardSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
	expectedRevision uint64,
) error {
	if err := a.check(ctx); err != nil {
		return err
	}
	return a.components.Sources.Discard(
		ctx,
		rootID,
		sourceID,
		expectedRevision,
	)
}

func (a *API) GetSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
) (source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, err
	}
	return a.components.Sources.Get(ctx, rootID, sourceID)
}

func (a *API) ListSources(
	ctx context.Context,
	rootID root.RootID,
) ([]source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	return a.components.Sources.List(ctx, rootID)
}

func (a *API) UpdateSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
	update source.Update,
) (source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, err
	}
	update.Config = append(json.RawMessage(nil), update.Config...)
	return a.components.Sources.Update(
		ctx,
		rootID,
		sourceID,
		update,
	)
}

func (a *API) RetireSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
	expectedRevision uint64,
) (source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, err
	}
	return a.components.Sources.Retire(
		ctx,
		rootID,
		sourceID,
		expectedRevision,
	)
}

func (a *API) PurgeSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
	expectedRevision uint64,
) error {
	if err := a.check(ctx); err != nil {
		return err
	}
	return a.components.Sources.Purge(
		ctx,
		rootID,
		sourceID,
		expectedRevision,
	)
}

func (a *API) ListSourceKinds(
	ctx context.Context,
) ([]source.SourceKind, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	return a.components.Sources.Kinds(), nil
}

func (a *API) Close() error {
	if a == nil {
		return nil
	}
	a.closeOnce.Do(func() {
		a.closed.Store(true)
		if a.components != nil {
			a.closeErr = a.components.Close()
		}
		a.resources = nil
		a.components = nil
	})
	return a.closeErr
}

func (a *API) CreateArtifactRoot(
	request *CreateArtifactRootRequest,
) (*CreateArtifactRootResponse, error) {
	if err := requireRequest(request, "create Artifact Root request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Root body"); err != nil {
		return nil, err
	}
	value, err := a.CreateRoot(context.Background(), *request.Body)
	if err != nil {
		return nil, err
	}
	return &CreateArtifactRootResponse{Body: &value}, nil
}

func (a *API) GetArtifactRoot(
	request *GetArtifactRootRequest,
) (*GetArtifactRootResponse, error) {
	if err := requireRequest(request, "get Artifact Root request"); err != nil {
		return nil, err
	}
	value, err := a.GetRoot(context.Background(), request.RootID)
	if err != nil {
		return nil, err
	}
	return &GetArtifactRootResponse{Body: &value}, nil
}

func (a *API) ListArtifactRoots(
	_ *ListArtifactRootsRequest,
) (*ListArtifactRootsResponse, error) {
	values, err := a.ListRoots(context.Background())
	if err != nil {
		return nil, err
	}
	return &ListArtifactRootsResponse{
		Body: &ListArtifactRootsResponseBody{Roots: values},
	}, nil
}

func (a *API) UpdateArtifactRoot(
	request *UpdateArtifactRootRequest,
) (*UpdateArtifactRootResponse, error) {
	if err := requireRequest(request, "update Artifact Root request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Root update body"); err != nil {
		return nil, err
	}
	value, err := a.UpdateRoot(
		context.Background(),
		request.RootID,
		*request.Body,
	)
	if err != nil {
		return nil, err
	}
	return &UpdateArtifactRootResponse{Body: &value}, nil
}

func (a *API) RetireArtifactRoot(
	request *RetireArtifactRootRequest,
) (*RetireArtifactRootResponse, error) {
	if err := requireRequest(request, "retire Artifact Root request"); err != nil {
		return nil, err
	}
	value, err := a.RetireRoot(
		context.Background(),
		request.RootID,
		request.ExpectedRevision,
	)
	if err != nil {
		return nil, err
	}
	return &RetireArtifactRootResponse{Body: &value}, nil
}

func (a *API) PurgeArtifactRoot(
	request *PurgeArtifactRootRequest,
) (*PurgeArtifactRootResponse, error) {
	if err := requireRequest(request, "purge Artifact Root request"); err != nil {
		return nil, err
	}
	if err := a.PurgeRoot(
		context.Background(),
		request.RootID,
		request.ExpectedRevision,
	); err != nil {
		return nil, err
	}
	return &PurgeArtifactRootResponse{RootID: request.RootID}, nil
}

func (a *API) CreateArtifactSource(
	request *CreateArtifactSourceRequest,
) (*CreateArtifactSourceResponse, error) {
	if err := requireRequest(request, "create Artifact Source request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Source body"); err != nil {
		return nil, err
	}
	value, err := a.CreateSource(
		context.Background(),
		request.RootID,
		*request.Body,
	)
	if err != nil {
		return nil, err
	}
	return &CreateArtifactSourceResponse{Body: &value}, nil
}

func (a *API) GetArtifactSource(
	request *GetArtifactSourceRequest,
) (*GetArtifactSourceResponse, error) {
	if err := requireRequest(request, "get Artifact Source request"); err != nil {
		return nil, err
	}
	value, err := a.GetSource(
		context.Background(),
		request.RootID,
		request.SourceID,
	)
	if err != nil {
		return nil, err
	}
	return &GetArtifactSourceResponse{Body: &value}, nil
}

func (a *API) ListArtifactSources(
	request *ListArtifactSourcesRequest,
) (*ListArtifactSourcesResponse, error) {
	if err := requireRequest(request, "list Artifact Sources request"); err != nil {
		return nil, err
	}
	values, err := a.ListSources(context.Background(), request.RootID)
	if err != nil {
		return nil, err
	}
	return &ListArtifactSourcesResponse{
		Body: &ListArtifactSourcesResponseBody{Sources: values},
	}, nil
}

func (a *API) UpdateArtifactSource(
	request *UpdateArtifactSourceRequest,
) (*UpdateArtifactSourceResponse, error) {
	if err := requireRequest(request, "update Artifact Source request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Source update body"); err != nil {
		return nil, err
	}
	value, err := a.UpdateSource(
		context.Background(),
		request.RootID,
		request.SourceID,
		*request.Body,
	)
	if err != nil {
		return nil, err
	}
	return &UpdateArtifactSourceResponse{Body: &value}, nil
}

func (a *API) RetireArtifactSource(
	request *RetireArtifactSourceRequest,
) (*RetireArtifactSourceResponse, error) {
	if err := requireRequest(request, "retire Artifact Source request"); err != nil {
		return nil, err
	}
	value, err := a.RetireSource(
		context.Background(),
		request.RootID,
		request.SourceID,
		request.ExpectedRevision,
	)
	if err != nil {
		return nil, err
	}
	return &RetireArtifactSourceResponse{Body: &value}, nil
}

func (a *API) PurgeArtifactSource(
	request *PurgeArtifactSourceRequest,
) (*PurgeArtifactSourceResponse, error) {
	if err := requireRequest(request, "purge Artifact Source request"); err != nil {
		return nil, err
	}
	if err := a.PurgeSource(
		context.Background(),
		request.RootID,
		request.SourceID,
		request.ExpectedRevision,
	); err != nil {
		return nil, err
	}
	return &PurgeArtifactSourceResponse{
		RootID:   request.RootID,
		SourceID: request.SourceID,
	}, nil
}

func (a *API) ListArtifactSourceKinds(
	_ *ListArtifactSourceKindsRequest,
) (*ListArtifactSourceKindsResponse, error) {
	kinds, err := a.ListSourceKinds(context.Background())
	if err != nil {
		return nil, err
	}
	return &ListArtifactSourceKindsResponse{
		Body: &ListArtifactSourceKindsResponseBody{Kinds: kinds},
	}, nil
}

func (a *API) CreateArtifactCollection(
	request *CreateArtifactCollectionRequest,
) (*CreateArtifactCollectionResponse, error) {
	if err := requireRequest(request, "create Artifact Collection request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Collection body"); err != nil {
		return nil, err
	}
	value, attachments, err := a.CreateCollection(
		context.Background(),
		request.RootID,
		request.Body.Draft,
		request.Body.Attachments,
	)
	if err != nil {
		return nil, err
	}
	return &CreateArtifactCollectionResponse{
		Body: &CreateArtifactCollectionResponseBody{
			Collection:  value,
			Attachments: attachments,
		},
	}, nil
}

func (a *API) GetArtifactCollection(
	request *GetArtifactCollectionRequest,
) (*GetArtifactCollectionResponse, error) {
	if err := requireRequest(request, "get Artifact Collection request"); err != nil {
		return nil, err
	}
	value, err := a.GetCollection(context.Background(), request.Collection)
	if err != nil {
		return nil, err
	}
	return &GetArtifactCollectionResponse{Body: &value}, nil
}

func (a *API) ListArtifactCollections(
	request *ListArtifactCollectionsRequest,
) (*ListArtifactCollectionsResponse, error) {
	if err := requireRequest(request, "list Artifact Collections request"); err != nil {
		return nil, err
	}
	values, err := a.ListCollections(context.Background(), request.RootID)
	if err != nil {
		return nil, err
	}
	return &ListArtifactCollectionsResponse{
		Body: &ListArtifactCollectionsResponseBody{Collections: values},
	}, nil
}

func (a *API) UpdateArtifactCollection(
	request *UpdateArtifactCollectionRequest,
) (*UpdateArtifactCollectionResponse, error) {
	if err := requireRequest(request, "update Artifact Collection request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Collection update body"); err != nil {
		return nil, err
	}
	value, err := a.UpdateCollection(
		context.Background(),
		request.Collection,
		*request.Body,
	)
	if err != nil {
		return nil, err
	}
	return &UpdateArtifactCollectionResponse{Body: &value}, nil
}

func (a *API) RetireArtifactCollection(
	request *RetireArtifactCollectionRequest,
) (*RetireArtifactCollectionResponse, error) {
	if err := requireRequest(request, "retire Artifact Collection request"); err != nil {
		return nil, err
	}
	value, err := a.RetireCollection(
		context.Background(),
		request.Collection,
		request.ExpectedRevision,
	)
	if err != nil {
		return nil, err
	}
	return &RetireArtifactCollectionResponse{Body: &value}, nil
}

func (a *API) PurgeArtifactCollection(
	request *PurgeArtifactCollectionRequest,
) (*PurgeArtifactCollectionResponse, error) {
	if err := requireRequest(request, "purge Artifact Collection request"); err != nil {
		return nil, err
	}
	if err := a.PurgeCollection(
		context.Background(),
		request.Collection,
		request.ExpectedRevision,
	); err != nil {
		return nil, err
	}
	return &PurgeArtifactCollectionResponse{
		Collection: request.Collection,
	}, nil
}

func (a *API) AttachArtifactCollectionSource(
	request *AttachArtifactCollectionSourceRequest,
) (*AttachArtifactCollectionSourceResponse, error) {
	if err := requireRequest(request, "attach Artifact Collection Source request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Collection attachment body"); err != nil {
		return nil, err
	}
	value, attachment, err := a.AttachCollectionSource(
		context.Background(),
		request.Collection,
		request.Body.ExpectedCollectionRevision,
		request.Body.Attachment,
	)
	if err != nil {
		return nil, err
	}
	return &AttachArtifactCollectionSourceResponse{
		Body: &ArtifactCollectionAttachmentResponseBody{
			Collection: value,
			Attachment: attachment,
		},
	}, nil
}

func (a *API) GetArtifactCollectionAttachment(
	request *GetArtifactCollectionAttachmentRequest,
) (*GetArtifactCollectionAttachmentResponse, error) {
	if err := requireRequest(request, "get Artifact Collection attachment request"); err != nil {
		return nil, err
	}
	value, err := a.GetCollectionAttachment(
		context.Background(),
		request.Collection,
		request.SourceID,
	)
	if err != nil {
		return nil, err
	}
	return &GetArtifactCollectionAttachmentResponse{Body: &value}, nil
}

func (a *API) ListArtifactCollectionAttachments(
	request *ListArtifactCollectionAttachmentsRequest,
) (*ListArtifactCollectionAttachmentsResponse, error) {
	if err := requireRequest(request, "list Artifact Collection attachments request"); err != nil {
		return nil, err
	}
	values, err := a.ListCollectionAttachments(
		context.Background(),
		request.Collection,
	)
	if err != nil {
		return nil, err
	}
	return &ListArtifactCollectionAttachmentsResponse{
		Body: &ListArtifactCollectionAttachmentsResponseBody{
			Attachments: values,
		},
	}, nil
}

func (a *API) UpdateArtifactCollectionAttachment(
	request *UpdateArtifactCollectionAttachmentRequest,
) (*UpdateArtifactCollectionAttachmentResponse, error) {
	if err := requireRequest(request, "update Artifact Collection attachment request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Collection attachment update body"); err != nil {
		return nil, err
	}
	value, attachment, err := a.UpdateCollectionAttachment(
		context.Background(),
		request.Collection,
		request.SourceID,
		*request.Body,
	)
	if err != nil {
		return nil, err
	}
	return &UpdateArtifactCollectionAttachmentResponse{
		Body: &ArtifactCollectionAttachmentResponseBody{
			Collection: value,
			Attachment: attachment,
		},
	}, nil
}

func (a *API) DetachArtifactCollectionSource(
	request *DetachArtifactCollectionSourceRequest,
) (*DetachArtifactCollectionSourceResponse, error) {
	if err := requireRequest(request, "detach Artifact Collection Source request"); err != nil {
		return nil, err
	}
	value, err := a.DetachCollectionSource(
		context.Background(),
		request.Collection,
		request.SourceID,
		request.ExpectedCollectionRevision,
		request.ExpectedAttachmentRevision,
	)
	if err != nil {
		return nil, err
	}
	return &DetachArtifactCollectionSourceResponse{Body: &value}, nil
}

func (a *API) ReplaceArtifactCollectionAttachment(
	request *ReplaceArtifactCollectionAttachmentRequest,
) (*ReplaceArtifactCollectionAttachmentResponse, error) {
	if err := requireRequest(request, "replace Artifact Collection attachment request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Collection attachment replacement body"); err != nil {
		return nil, err
	}
	value, attachment, err := a.ReplaceCollectionAttachment(
		context.Background(),
		request.Collection,
		*request.Body,
	)
	if err != nil {
		return nil, err
	}
	return &ReplaceArtifactCollectionAttachmentResponse{
		Body: &ArtifactCollectionAttachmentResponseBody{
			Collection: value,
			Attachment: attachment,
		},
	}, nil
}

func (a *API) GetArtifactRecord(
	request *GetArtifactRecordRequest,
) (*GetArtifactRecordResponse, error) {
	if err := requireRequest(request, "get Artifact record request"); err != nil {
		return nil, err
	}
	value, err := a.GetArtifact(context.Background(), request.Artifact)
	if err != nil {
		return nil, err
	}
	return &GetArtifactRecordResponse{Body: &value}, nil
}

func (a *API) ListArtifactCollectionRecords(
	request *ListArtifactCollectionRecordsRequest,
) (*ListArtifactCollectionRecordsResponse, error) {
	if err := requireRequest(request, "list Artifact Collection records request"); err != nil {
		return nil, err
	}
	values, err := a.ListCollectionArtifacts(
		context.Background(),
		request.Collection,
	)
	if err != nil {
		return nil, err
	}
	return &ListArtifactCollectionRecordsResponse{
		Body: &ListArtifactCollectionRecordsResponseBody{
			Artifacts: values,
		},
	}, nil
}

func (a *API) AdoptArtifactRecord(
	request *AdoptArtifactRecordRequest,
) (*AdoptArtifactRecordResponse, error) {
	if err := requireRequest(request, "adopt Artifact request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact adoption body"); err != nil {
		return nil, err
	}
	value, err := a.AdoptArtifact(context.Background(), *request.Body)
	if err != nil {
		return nil, err
	}
	return &AdoptArtifactRecordResponse{Body: &value}, nil
}

func (a *API) PinArtifactRecord(
	request *PinArtifactRecordRequest,
) (*PinArtifactRecordResponse, error) {
	if err := requireRequest(request, "pin Artifact request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact pin body"); err != nil {
		return nil, err
	}
	value, err := a.PinArtifact(context.Background(), *request.Body)
	if err != nil {
		return nil, err
	}
	return &PinArtifactRecordResponse{Body: &value}, nil
}

func (a *API) SetArtifactRecordEnabled(
	request *SetArtifactRecordEnabledRequest,
) (*SetArtifactRecordEnabledResponse, error) {
	if err := requireRequest(request, "set Artifact enabled request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact enabled body"); err != nil {
		return nil, err
	}
	value, err := a.SetArtifactEnabled(
		context.Background(),
		request.Artifact,
		request.Body.ExpectedRevision,
		request.Body.Enabled,
	)
	if err != nil {
		return nil, err
	}
	return &SetArtifactRecordEnabledResponse{Body: &value}, nil
}

func (a *API) SetArtifactRecordName(
	request *SetArtifactRecordNameRequest,
) (*SetArtifactRecordNameResponse, error) {
	if err := requireRequest(request, "set Artifact name request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact name body"); err != nil {
		return nil, err
	}
	value, err := a.SetArtifactName(
		context.Background(),
		request.Artifact,
		request.Body.ExpectedRevision,
		request.Body.Name,
	)
	if err != nil {
		return nil, err
	}
	return &SetArtifactRecordNameResponse{Body: &value}, nil
}

func (a *API) UpdateArtifactRecordData(
	request *UpdateArtifactRecordDataRequest,
) (*UpdateArtifactRecordDataResponse, error) {
	if err := requireRequest(request, "update Artifact data request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact data body"); err != nil {
		return nil, err
	}
	value, err := a.UpdateArtifactData(
		context.Background(),
		request.Artifact,
		request.Body.ExpectedRevision,
		request.Body.Data,
	)
	if err != nil {
		return nil, err
	}
	return &UpdateArtifactRecordDataResponse{Body: &value}, nil
}

func (a *API) UnadoptArtifactRecord(
	request *UnadoptArtifactRecordRequest,
) (*UnadoptArtifactRecordResponse, error) {
	if err := requireRequest(request, "unadopt Artifact request"); err != nil {
		return nil, err
	}
	if err := a.UnadoptArtifact(
		context.Background(),
		request.Artifact,
		request.ExpectedRevision,
		request.Suppress,
	); err != nil {
		return nil, err
	}
	return &UnadoptArtifactRecordResponse{Artifact: request.Artifact}, nil
}

func (a *API) PurgeArtifactRecord(
	request *PurgeArtifactRecordRequest,
) (*PurgeArtifactRecordResponse, error) {
	if err := requireRequest(request, "purge Artifact request"); err != nil {
		return nil, err
	}
	if err := a.PurgeArtifact(
		context.Background(),
		request.Artifact,
		request.ExpectedRevision,
	); err != nil {
		return nil, err
	}
	return &PurgeArtifactRecordResponse{Artifact: request.Artifact}, nil
}

func (a *API) PurgeAndSuppressArtifactRecord(
	request *PurgeAndSuppressArtifactRecordRequest,
) (*PurgeAndSuppressArtifactRecordResponse, error) {
	if err := requireRequest(request, "purge and suppress Artifact request"); err != nil {
		return nil, err
	}
	if err := a.PurgeAndSuppressArtifact(
		context.Background(),
		request.Artifact,
		request.ExpectedRevision,
	); err != nil {
		return nil, err
	}
	return &PurgeAndSuppressArtifactRecordResponse{
		Artifact: request.Artifact,
	}, nil
}

func (a *API) ListArtifactCollectionSuppressions(
	request *ListArtifactCollectionSuppressionsRequest,
) (*ListArtifactCollectionSuppressionsResponse, error) {
	if err := requireRequest(request, "list Artifact suppressions request"); err != nil {
		return nil, err
	}
	values, err := a.ListCollectionSuppressions(
		context.Background(),
		request.Collection,
	)
	if err != nil {
		return nil, err
	}
	return &ListArtifactCollectionSuppressionsResponse{
		Body: &ListArtifactCollectionSuppressionsResponseBody{
			Suppressions: values,
		},
	}, nil
}

func (a *API) SuppressArtifactBinding(
	request *SuppressArtifactBindingRequest,
) (*SuppressArtifactBindingResponse, error) {
	if err := requireRequest(request, "suppress Artifact binding request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact suppression body"); err != nil {
		return nil, err
	}
	value, err := a.SuppressBinding(context.Background(), *request.Body)
	if err != nil {
		return nil, err
	}
	return &SuppressArtifactBindingResponse{Body: &value}, nil
}

func (a *API) UnsuppressArtifactBinding(
	request *UnsuppressArtifactBindingRequest,
) (*UnsuppressArtifactBindingResponse, error) {
	if err := requireRequest(request, "unsuppress Artifact binding request"); err != nil {
		return nil, err
	}
	if err := a.UnsuppressBinding(
		context.Background(),
		request.Collection,
		request.Binding,
		request.ExpectedRevision,
	); err != nil {
		return nil, err
	}
	return &UnsuppressArtifactBindingResponse{
		Collection: request.Collection,
		Binding:    request.Binding,
	}, nil
}

func (a *API) RefreshArtifactCollection(
	request *RefreshArtifactCollectionRequest,
) (*RefreshArtifactCollectionResponse, error) {
	if err := requireRequest(request, "refresh Artifact Collection request"); err != nil {
		return nil, err
	}
	value, err := a.RefreshCollection(
		context.Background(),
		request.Collection,
	)
	if err != nil {
		return nil, err
	}
	return &RefreshArtifactCollectionResponse{Body: &value}, nil
}

func (a *API) GetArtifactCollectionCatalog(
	request *GetArtifactCollectionCatalogRequest,
) (*GetArtifactCollectionCatalogResponse, error) {
	if err := requireRequest(request, "get Artifact Collection catalog request"); err != nil {
		return nil, err
	}
	value, err := a.CurrentCollectionCatalog(
		context.Background(),
		request.Collection,
	)
	if err != nil {
		return nil, err
	}
	return &GetArtifactCollectionCatalogResponse{Body: &value}, nil
}

func (a *API) InspectArtifactCollectionCatalog(
	request *InspectArtifactCollectionCatalogRequest,
) (*InspectArtifactCollectionCatalogResponse, error) {
	if err := requireRequest(request, "inspect Artifact Collection catalog request"); err != nil {
		return nil, err
	}
	value, err := a.InspectCollectionCatalog(
		context.Background(),
		request.Collection,
	)
	if err != nil {
		return nil, err
	}
	return &InspectArtifactCollectionCatalogResponse{Body: &value}, nil
}

func requireRequest[T any](value *T, subject string) error {
	if value != nil {
		return nil
	}
	return fmt.Errorf("%w: %s is required", basespec.ErrInvalid, subject)
}

func requireBody[T any](value *T, subject string) error {
	if value != nil {
		return nil
	}
	return fmt.Errorf("%w: %s is required", basespec.ErrInvalid, subject)
}

func (a *API) check(ctx context.Context) error {
	if a == nil ||
		a.closed.Load() ||
		a.components == nil ||
		a.components.Roots == nil ||
		a.components.Sources == nil {
		return basespec.ErrClosed
	}
	if ctx == nil {
		return fmt.Errorf(
			"%w: artifact store API context is nil",
			basespec.ErrInvalid,
		)
	}
	return ctx.Err()
}
