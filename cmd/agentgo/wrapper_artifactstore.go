package main

import (
	"errors"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	artifactConsumerAPI "github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
)

// ArtifactStoreWrapper is the command binding surface for Artifact Store.
// It contains transport recovery only. "consumerapi.API" owns validation, lifecycle dispatch, and first-level command
// handling.
type ArtifactStoreWrapper struct {
	api *artifactConsumerAPI.API
}

func InitArtifactStoreWrapper(
	wrapper *ArtifactStoreWrapper,
	api *artifactConsumerAPI.API,
) error {
	if wrapper == nil {
		return errors.New("artifact store wrapper is required")
	}
	if api == nil {
		return errors.New("artifact store consumer API is required")
	}
	wrapper.api = api
	return nil
}

func withArtifactStore[T any](
	wrapper *ArtifactStoreWrapper,
	call func(*artifactConsumerAPI.API) (T, error),
) (T, error) {
	return middleware.WithRecoveryResp(func() (T, error) {
		var zero T
		if wrapper == nil || wrapper.api == nil {
			return zero, basespec.ErrClosed
		}
		return call(wrapper.api)
	})
}

func (w *ArtifactStoreWrapper) CreateArtifactRoot(
	request *artifactConsumerAPI.CreateArtifactRootRequest,
) (*artifactConsumerAPI.CreateArtifactRootResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.CreateArtifactRootResponse, error) {
			return api.CreateArtifactRoot(request)
		},
	)
}

func (w *ArtifactStoreWrapper) GetArtifactRoot(
	request *artifactConsumerAPI.GetArtifactRootRequest,
) (*artifactConsumerAPI.GetArtifactRootResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.GetArtifactRootResponse, error) {
			return api.GetArtifactRoot(request)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactRoots(
	request *artifactConsumerAPI.ListArtifactRootsRequest,
) (*artifactConsumerAPI.ListArtifactRootsResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.ListArtifactRootsResponse, error) {
			return api.ListArtifactRoots(request)
		},
	)
}

func (w *ArtifactStoreWrapper) UpdateArtifactRoot(
	request *artifactConsumerAPI.UpdateArtifactRootRequest,
) (*artifactConsumerAPI.UpdateArtifactRootResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.UpdateArtifactRootResponse, error) {
			return api.UpdateArtifactRoot(request)
		},
	)
}

func (w *ArtifactStoreWrapper) RetireArtifactRoot(
	request *artifactConsumerAPI.RetireArtifactRootRequest,
) (*artifactConsumerAPI.RetireArtifactRootResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.RetireArtifactRootResponse, error) {
			return api.RetireArtifactRoot(request)
		},
	)
}

func (w *ArtifactStoreWrapper) PurgeArtifactRoot(
	request *artifactConsumerAPI.PurgeArtifactRootRequest,
) (*artifactConsumerAPI.PurgeArtifactRootResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.PurgeArtifactRootResponse, error) {
			return api.PurgeArtifactRoot(request)
		},
	)
}

func (w *ArtifactStoreWrapper) CreateArtifactSource(
	request *artifactConsumerAPI.CreateArtifactSourceRequest,
) (*artifactConsumerAPI.CreateArtifactSourceResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.CreateArtifactSourceResponse, error) {
			return api.CreateArtifactSource(request)
		},
	)
}

func (w *ArtifactStoreWrapper) GetArtifactSource(
	request *artifactConsumerAPI.GetArtifactSourceRequest,
) (*artifactConsumerAPI.GetArtifactSourceResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.GetArtifactSourceResponse, error) {
			return api.GetArtifactSource(request)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactSources(
	request *artifactConsumerAPI.ListArtifactSourcesRequest,
) (*artifactConsumerAPI.ListArtifactSourcesResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.ListArtifactSourcesResponse, error) {
			return api.ListArtifactSources(request)
		},
	)
}

func (w *ArtifactStoreWrapper) UpdateArtifactSource(
	request *artifactConsumerAPI.UpdateArtifactSourceRequest,
) (*artifactConsumerAPI.UpdateArtifactSourceResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.UpdateArtifactSourceResponse, error) {
			return api.UpdateArtifactSource(request)
		},
	)
}

func (w *ArtifactStoreWrapper) RetireArtifactSource(
	request *artifactConsumerAPI.RetireArtifactSourceRequest,
) (*artifactConsumerAPI.RetireArtifactSourceResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.RetireArtifactSourceResponse, error) {
			return api.RetireArtifactSource(request)
		},
	)
}

func (w *ArtifactStoreWrapper) PurgeArtifactSource(
	request *artifactConsumerAPI.PurgeArtifactSourceRequest,
) (*artifactConsumerAPI.PurgeArtifactSourceResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.PurgeArtifactSourceResponse, error) {
			return api.PurgeArtifactSource(request)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactSourceKinds(
	request *artifactConsumerAPI.ListArtifactSourceKindsRequest,
) (*artifactConsumerAPI.ListArtifactSourceKindsResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.ListArtifactSourceKindsResponse, error) {
			return api.ListArtifactSourceKinds(request)
		},
	)
}

func (w *ArtifactStoreWrapper) CreateArtifactCollection(
	request *artifactConsumerAPI.CreateArtifactCollectionRequest,
) (*artifactConsumerAPI.CreateArtifactCollectionResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.CreateArtifactCollectionResponse, error) {
			return api.CreateArtifactCollection(request)
		},
	)
}

func (w *ArtifactStoreWrapper) GetArtifactCollection(
	request *artifactConsumerAPI.GetArtifactCollectionRequest,
) (*artifactConsumerAPI.GetArtifactCollectionResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.GetArtifactCollectionResponse, error) {
			return api.GetArtifactCollection(request)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactCollections(
	request *artifactConsumerAPI.ListArtifactCollectionsRequest,
) (*artifactConsumerAPI.ListArtifactCollectionsResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.ListArtifactCollectionsResponse, error) {
			return api.ListArtifactCollections(request)
		},
	)
}

func (w *ArtifactStoreWrapper) UpdateArtifactCollection(
	request *artifactConsumerAPI.UpdateArtifactCollectionRequest,
) (*artifactConsumerAPI.UpdateArtifactCollectionResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.UpdateArtifactCollectionResponse, error) {
			return api.UpdateArtifactCollection(request)
		},
	)
}

func (w *ArtifactStoreWrapper) RetireArtifactCollection(
	request *artifactConsumerAPI.RetireArtifactCollectionRequest,
) (*artifactConsumerAPI.RetireArtifactCollectionResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.RetireArtifactCollectionResponse, error) {
			return api.RetireArtifactCollection(request)
		},
	)
}

func (w *ArtifactStoreWrapper) PurgeArtifactCollection(
	request *artifactConsumerAPI.PurgeArtifactCollectionRequest,
) (*artifactConsumerAPI.PurgeArtifactCollectionResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.PurgeArtifactCollectionResponse, error) {
			return api.PurgeArtifactCollection(request)
		},
	)
}

func (w *ArtifactStoreWrapper) AttachArtifactCollectionSource(
	request *artifactConsumerAPI.AttachArtifactCollectionSourceRequest,
) (*artifactConsumerAPI.AttachArtifactCollectionSourceResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.AttachArtifactCollectionSourceResponse, error) {
			return api.AttachArtifactCollectionSource(request)
		},
	)
}

func (w *ArtifactStoreWrapper) GetArtifactCollectionAttachment(
	request *artifactConsumerAPI.GetArtifactCollectionAttachmentRequest,
) (*artifactConsumerAPI.GetArtifactCollectionAttachmentResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.GetArtifactCollectionAttachmentResponse, error) {
			return api.GetArtifactCollectionAttachment(request)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactCollectionAttachments(
	request *artifactConsumerAPI.ListArtifactCollectionAttachmentsRequest,
) (*artifactConsumerAPI.ListArtifactCollectionAttachmentsResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.ListArtifactCollectionAttachmentsResponse, error) {
			return api.ListArtifactCollectionAttachments(request)
		},
	)
}

func (w *ArtifactStoreWrapper) UpdateArtifactCollectionAttachment(
	request *artifactConsumerAPI.UpdateArtifactCollectionAttachmentRequest,
) (*artifactConsumerAPI.UpdateArtifactCollectionAttachmentResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.UpdateArtifactCollectionAttachmentResponse, error) {
			return api.UpdateArtifactCollectionAttachment(request)
		},
	)
}

func (w *ArtifactStoreWrapper) DetachArtifactCollectionSource(
	request *artifactConsumerAPI.DetachArtifactCollectionSourceRequest,
) (*artifactConsumerAPI.DetachArtifactCollectionSourceResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.DetachArtifactCollectionSourceResponse, error) {
			return api.DetachArtifactCollectionSource(request)
		},
	)
}

func (w *ArtifactStoreWrapper) ReplaceArtifactCollectionAttachment(
	request *artifactConsumerAPI.ReplaceArtifactCollectionAttachmentRequest,
) (*artifactConsumerAPI.ReplaceArtifactCollectionAttachmentResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.ReplaceArtifactCollectionAttachmentResponse, error) {
			return api.ReplaceArtifactCollectionAttachment(request)
		},
	)
}

func (w *ArtifactStoreWrapper) GetArtifactRecord(
	request *artifactConsumerAPI.GetArtifactRecordRequest,
) (*artifactConsumerAPI.GetArtifactRecordResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.GetArtifactRecordResponse, error) {
			return api.GetArtifactRecord(request)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactCollectionRecords(
	request *artifactConsumerAPI.ListArtifactCollectionRecordsRequest,
) (*artifactConsumerAPI.ListArtifactCollectionRecordsResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.ListArtifactCollectionRecordsResponse, error) {
			return api.ListArtifactCollectionRecords(request)
		},
	)
}

func (w *ArtifactStoreWrapper) AdoptArtifactRecord(
	request *artifactConsumerAPI.AdoptArtifactRecordRequest,
) (*artifactConsumerAPI.AdoptArtifactRecordResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.AdoptArtifactRecordResponse, error) {
			return api.AdoptArtifactRecord(request)
		},
	)
}

func (w *ArtifactStoreWrapper) PinArtifactRecord(
	request *artifactConsumerAPI.PinArtifactRecordRequest,
) (*artifactConsumerAPI.PinArtifactRecordResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.PinArtifactRecordResponse, error) {
			return api.PinArtifactRecord(request)
		},
	)
}

func (w *ArtifactStoreWrapper) SetArtifactRecordEnabled(
	request *artifactConsumerAPI.SetArtifactRecordEnabledRequest,
) (*artifactConsumerAPI.SetArtifactRecordEnabledResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.SetArtifactRecordEnabledResponse, error) {
			return api.SetArtifactRecordEnabled(request)
		},
	)
}

func (w *ArtifactStoreWrapper) SetArtifactRecordName(
	request *artifactConsumerAPI.SetArtifactRecordNameRequest,
) (*artifactConsumerAPI.SetArtifactRecordNameResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.SetArtifactRecordNameResponse, error) {
			return api.SetArtifactRecordName(request)
		},
	)
}

func (w *ArtifactStoreWrapper) UpdateArtifactRecordData(
	request *artifactConsumerAPI.UpdateArtifactRecordDataRequest,
) (*artifactConsumerAPI.UpdateArtifactRecordDataResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.UpdateArtifactRecordDataResponse, error) {
			return api.UpdateArtifactRecordData(request)
		},
	)
}

func (w *ArtifactStoreWrapper) UnadoptArtifactRecord(
	request *artifactConsumerAPI.UnadoptArtifactRecordRequest,
) (*artifactConsumerAPI.UnadoptArtifactRecordResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.UnadoptArtifactRecordResponse, error) {
			return api.UnadoptArtifactRecord(request)
		},
	)
}

func (w *ArtifactStoreWrapper) PurgeArtifactRecord(
	request *artifactConsumerAPI.PurgeArtifactRecordRequest,
) (*artifactConsumerAPI.PurgeArtifactRecordResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.PurgeArtifactRecordResponse, error) {
			return api.PurgeArtifactRecord(request)
		},
	)
}

func (w *ArtifactStoreWrapper) PurgeAndSuppressArtifactRecord(
	request *artifactConsumerAPI.PurgeAndSuppressArtifactRecordRequest,
) (*artifactConsumerAPI.PurgeAndSuppressArtifactRecordResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.PurgeAndSuppressArtifactRecordResponse, error) {
			return api.PurgeAndSuppressArtifactRecord(request)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactCollectionSuppressions(
	request *artifactConsumerAPI.ListArtifactCollectionSuppressionsRequest,
) (*artifactConsumerAPI.ListArtifactCollectionSuppressionsResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.ListArtifactCollectionSuppressionsResponse, error) {
			return api.ListArtifactCollectionSuppressions(request)
		},
	)
}

func (w *ArtifactStoreWrapper) SuppressArtifactBinding(
	request *artifactConsumerAPI.SuppressArtifactBindingRequest,
) (*artifactConsumerAPI.SuppressArtifactBindingResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.SuppressArtifactBindingResponse, error) {
			return api.SuppressArtifactBinding(request)
		},
	)
}

func (w *ArtifactStoreWrapper) UnsuppressArtifactBinding(
	request *artifactConsumerAPI.UnsuppressArtifactBindingRequest,
) (*artifactConsumerAPI.UnsuppressArtifactBindingResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.UnsuppressArtifactBindingResponse, error) {
			return api.UnsuppressArtifactBinding(request)
		},
	)
}

func (w *ArtifactStoreWrapper) RefreshArtifactCollection(
	request *artifactConsumerAPI.RefreshArtifactCollectionRequest,
) (*artifactConsumerAPI.RefreshArtifactCollectionResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.RefreshArtifactCollectionResponse, error) {
			return api.RefreshArtifactCollection(request)
		},
	)
}

func (w *ArtifactStoreWrapper) GetArtifactCollectionCatalog(
	request *artifactConsumerAPI.GetArtifactCollectionCatalogRequest,
) (*artifactConsumerAPI.GetArtifactCollectionCatalogResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.GetArtifactCollectionCatalogResponse, error) {
			return api.GetArtifactCollectionCatalog(request)
		},
	)
}

func (w *ArtifactStoreWrapper) InspectArtifactCollectionCatalog(
	request *artifactConsumerAPI.InspectArtifactCollectionCatalogRequest,
) (*artifactConsumerAPI.InspectArtifactCollectionCatalogResponse, error) {
	return withArtifactStore(
		w,
		func(api *artifactConsumerAPI.API) (*artifactConsumerAPI.InspectArtifactCollectionCatalogResponse, error) {
			return api.InspectArtifactCollectionCatalog(request)
		},
	)
}

func (w *ArtifactStoreWrapper) Store() *artifactConsumerAPI.API {
	if w == nil {
		return nil
	}
	return w.api
}

func (w *ArtifactStoreWrapper) close() {
	if w != nil {
		w.api = nil
	}
}
