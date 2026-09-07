// Package wailsapi adapts the transport-independent consumer contract to
// request and response structures suitable for Wails binding generation.
package wailsapi

import (
	"context"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi/reqresp/lifecycle"
)

type API struct {
	consumer consumerapi.ConsumerAPI
}

func New(consumer consumerapi.ConsumerAPI) (*API, error) {
	if consumer == nil {
		return nil, fmt.Errorf(
			"%w: Artifact Store consumer API is required",
			basespec.ErrInvalid,
		)
	}
	return &API{consumer: consumer}, nil
}

func (a *API) CreateArtifactRoot(
	ctx context.Context,
	request *lifecycle.CreateArtifactRootRequest,
) (*lifecycle.CreateArtifactRootResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "create Artifact Root request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Root body"); err != nil {
		return nil, err
	}
	value, err := a.consumer.CreateRoot(ctx, *request.Body)
	if err != nil {
		return nil, err
	}
	return &lifecycle.CreateArtifactRootResponse{Body: &value}, nil
}

func (a *API) GetArtifactRoot(
	ctx context.Context,
	request *lifecycle.GetArtifactRootRequest,
) (*lifecycle.GetArtifactRootResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "get Artifact Root request"); err != nil {
		return nil, err
	}
	value, err := a.consumer.GetRoot(ctx, request.RootID)
	if err != nil {
		return nil, err
	}
	return &lifecycle.GetArtifactRootResponse{Body: &value}, nil
}

func (a *API) ListArtifactRoots(
	ctx context.Context,
	_ *lifecycle.ListArtifactRootsRequest,
) (*lifecycle.ListArtifactRootsResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	values, err := a.consumer.ListRoots(ctx)
	if err != nil {
		return nil, err
	}
	return &lifecycle.ListArtifactRootsResponse{
		Body: &lifecycle.ListArtifactRootsResponseBody{
			Roots: values,
		},
	}, nil
}

func (a *API) UpdateArtifactRoot(
	ctx context.Context,
	request *lifecycle.UpdateArtifactRootRequest,
) (*lifecycle.UpdateArtifactRootResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "update Artifact Root request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Root update body"); err != nil {
		return nil, err
	}
	value, err := a.consumer.UpdateRoot(
		ctx,
		request.RootID,
		*request.Body,
	)
	if err != nil {
		return nil, err
	}
	return &lifecycle.UpdateArtifactRootResponse{Body: &value}, nil
}

func (a *API) RetireArtifactRoot(
	ctx context.Context,
	request *lifecycle.RetireArtifactRootRequest,
) (*lifecycle.RetireArtifactRootResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "retire Artifact Root request"); err != nil {
		return nil, err
	}
	value, err := a.consumer.RetireRoot(
		ctx,
		request.RootID,
		request.ExpectedRevision,
	)
	if err != nil {
		return nil, err
	}
	return &lifecycle.RetireArtifactRootResponse{Body: &value}, nil
}

func (a *API) PurgeArtifactRoot(
	ctx context.Context,
	request *lifecycle.PurgeArtifactRootRequest,
) (*lifecycle.PurgeArtifactRootResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "purge Artifact Root request"); err != nil {
		return nil, err
	}
	if err := a.consumer.PurgeRoot(
		ctx,
		request.RootID,
		request.ExpectedRevision,
	); err != nil {
		return nil, err
	}
	return &lifecycle.PurgeArtifactRootResponse{
		RootID: request.RootID,
	}, nil
}

func (a *API) CreateArtifactSource(
	ctx context.Context,
	request *lifecycle.CreateArtifactSourceRequest,
) (*lifecycle.CreateArtifactSourceResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "create Artifact Source request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Source body"); err != nil {
		return nil, err
	}
	value, err := a.consumer.CreateSource(
		ctx,
		request.RootID,
		*request.Body,
	)
	if err != nil {
		return nil, err
	}
	return &lifecycle.CreateArtifactSourceResponse{Body: &value}, nil
}

func (a *API) GetArtifactSource(
	ctx context.Context,
	request *lifecycle.GetArtifactSourceRequest,
) (*lifecycle.GetArtifactSourceResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "get Artifact Source request"); err != nil {
		return nil, err
	}
	value, err := a.consumer.GetSource(
		ctx,
		request.RootID,
		request.SourceID,
	)
	if err != nil {
		return nil, err
	}
	return &lifecycle.GetArtifactSourceResponse{Body: &value}, nil
}

func (a *API) ListArtifactSources(
	ctx context.Context,
	request *lifecycle.ListArtifactSourcesRequest,
) (*lifecycle.ListArtifactSourcesResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "list Artifact Sources request"); err != nil {
		return nil, err
	}
	values, err := a.consumer.ListSources(ctx, request.RootID)
	if err != nil {
		return nil, err
	}
	return &lifecycle.ListArtifactSourcesResponse{
		Body: &lifecycle.ListArtifactSourcesResponseBody{
			Sources: values,
		},
	}, nil
}

func (a *API) UpdateArtifactSource(
	ctx context.Context,
	request *lifecycle.UpdateArtifactSourceRequest,
) (*lifecycle.UpdateArtifactSourceResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "update Artifact Source request"); err != nil {
		return nil, err
	}
	if err := requireBody(request.Body, "Artifact Source update body"); err != nil {
		return nil, err
	}
	value, err := a.consumer.UpdateSource(
		ctx,
		request.RootID,
		request.SourceID,
		*request.Body,
	)
	if err != nil {
		return nil, err
	}
	return &lifecycle.UpdateArtifactSourceResponse{Body: &value}, nil
}

func (a *API) RetireArtifactSource(
	ctx context.Context,
	request *lifecycle.RetireArtifactSourceRequest,
) (*lifecycle.RetireArtifactSourceResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "retire Artifact Source request"); err != nil {
		return nil, err
	}
	value, err := a.consumer.RetireSource(
		ctx,
		request.RootID,
		request.SourceID,
		request.ExpectedRevision,
	)
	if err != nil {
		return nil, err
	}
	return &lifecycle.RetireArtifactSourceResponse{Body: &value}, nil
}

func (a *API) PurgeArtifactSource(
	ctx context.Context,
	request *lifecycle.PurgeArtifactSourceRequest,
) (*lifecycle.PurgeArtifactSourceResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	if err := requireRequest(request, "purge Artifact Source request"); err != nil {
		return nil, err
	}
	if err := a.consumer.PurgeSource(
		ctx,
		request.RootID,
		request.SourceID,
		request.ExpectedRevision,
	); err != nil {
		return nil, err
	}
	return &lifecycle.PurgeArtifactSourceResponse{
		RootID:   request.RootID,
		SourceID: request.SourceID,
	}, nil
}

func (a *API) ListArtifactSourceKinds(
	ctx context.Context,
	_ *lifecycle.ListArtifactSourceKindsRequest,
) (*lifecycle.ListArtifactSourceKindsResponse, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	kinds, err := a.consumer.ListSourceKinds(ctx)
	if err != nil {
		return nil, err
	}
	return &lifecycle.ListArtifactSourceKindsResponse{
		Body: &lifecycle.ListArtifactSourceKindsResponseBody{
			Kinds: kinds,
		},
	}, nil
}

func (a *API) check(ctx context.Context) error {
	if a == nil || a.consumer == nil {
		return basespec.ErrClosed
	}
	if ctx == nil {
		return fmt.Errorf(
			"%w: Artifact Store Wails context is nil",
			basespec.ErrInvalid,
		)
	}
	return ctx.Err()
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
