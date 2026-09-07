package main

import (
	"context"
	"errors"

	artifactConsumerAPI "github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi"
	artifactConsumerAPIlifecycle "github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi/reqresp/lifecycle"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi/wailsapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
)

type ArtifactStoreWrapper struct {
	api      *wailsapi.API
	consumer artifactConsumerAPI.ConsumerAPI
}

func InitArtifactStoreWrapper(
	wrapper *ArtifactStoreWrapper,
	consumer artifactConsumerAPI.ConsumerAPI,
) error {
	if wrapper == nil {
		return errors.New("artifact store wrapper is required")
	}
	if consumer == nil {
		return errors.New("artifact store consumer API is required")
	}

	api, err := wailsapi.New(consumer)
	if err != nil {
		return err
	}
	wrapper.api = api
	wrapper.consumer = consumer
	return nil
}

func (w *ArtifactStoreWrapper) CreateArtifactRoot(
	request *artifactConsumerAPIlifecycle.CreateArtifactRootRequest,
) (*artifactConsumerAPIlifecycle.CreateArtifactRootResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.CreateArtifactRootResponse, error) {
			return w.api.CreateArtifactRoot(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) GetArtifactRoot(
	request *artifactConsumerAPIlifecycle.GetArtifactRootRequest,
) (*artifactConsumerAPIlifecycle.GetArtifactRootResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.GetArtifactRootResponse, error) {
			return w.api.GetArtifactRoot(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactRoots(
	request *artifactConsumerAPIlifecycle.ListArtifactRootsRequest,
) (*artifactConsumerAPIlifecycle.ListArtifactRootsResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.ListArtifactRootsResponse, error) {
			return w.api.ListArtifactRoots(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) UpdateArtifactRoot(
	request *artifactConsumerAPIlifecycle.UpdateArtifactRootRequest,
) (*artifactConsumerAPIlifecycle.UpdateArtifactRootResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.UpdateArtifactRootResponse, error) {
			return w.api.UpdateArtifactRoot(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) RetireArtifactRoot(
	request *artifactConsumerAPIlifecycle.RetireArtifactRootRequest,
) (*artifactConsumerAPIlifecycle.RetireArtifactRootResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.RetireArtifactRootResponse, error) {
			return w.api.RetireArtifactRoot(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) PurgeArtifactRoot(
	request *artifactConsumerAPIlifecycle.PurgeArtifactRootRequest,
) (*artifactConsumerAPIlifecycle.PurgeArtifactRootResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.PurgeArtifactRootResponse, error) {
			return w.api.PurgeArtifactRoot(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) CreateArtifactSource(
	request *artifactConsumerAPIlifecycle.CreateArtifactSourceRequest,
) (*artifactConsumerAPIlifecycle.CreateArtifactSourceResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.CreateArtifactSourceResponse, error) {
			return w.api.CreateArtifactSource(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) GetArtifactSource(
	request *artifactConsumerAPIlifecycle.GetArtifactSourceRequest,
) (*artifactConsumerAPIlifecycle.GetArtifactSourceResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.GetArtifactSourceResponse, error) {
			return w.api.GetArtifactSource(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactSources(
	request *artifactConsumerAPIlifecycle.ListArtifactSourcesRequest,
) (*artifactConsumerAPIlifecycle.ListArtifactSourcesResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.ListArtifactSourcesResponse, error) {
			return w.api.ListArtifactSources(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) UpdateArtifactSource(
	request *artifactConsumerAPIlifecycle.UpdateArtifactSourceRequest,
) (*artifactConsumerAPIlifecycle.UpdateArtifactSourceResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.UpdateArtifactSourceResponse, error) {
			return w.api.UpdateArtifactSource(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) RetireArtifactSource(
	request *artifactConsumerAPIlifecycle.RetireArtifactSourceRequest,
) (*artifactConsumerAPIlifecycle.RetireArtifactSourceResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.RetireArtifactSourceResponse, error) {
			return w.api.RetireArtifactSource(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) PurgeArtifactSource(
	request *artifactConsumerAPIlifecycle.PurgeArtifactSourceRequest,
) (*artifactConsumerAPIlifecycle.PurgeArtifactSourceResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.PurgeArtifactSourceResponse, error) {
			return w.api.PurgeArtifactSource(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) ListArtifactSourceKinds(
	request *artifactConsumerAPIlifecycle.ListArtifactSourceKindsRequest,
) (*artifactConsumerAPIlifecycle.ListArtifactSourceKindsResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*artifactConsumerAPIlifecycle.ListArtifactSourceKindsResponse, error) {
			return w.api.ListArtifactSourceKinds(
				context.Background(),
				request,
			)
		},
	)
}

func (w *ArtifactStoreWrapper) Store() artifactConsumerAPI.ConsumerAPI {
	if w == nil {
		return nil
	}
	return w.consumer
}

func (w *ArtifactStoreWrapper) close() {
	if w == nil {
		return
	}
	w.api = nil
	w.consumer = nil
}
