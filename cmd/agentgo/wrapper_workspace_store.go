package main

import (
	"context"
	"errors"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
	workspaceAggregate "github.com/flexigpt/flexigpt-app/internal/workspace/aggregate"
	workspaceConsumerAPI "github.com/flexigpt/flexigpt-app/internal/workspace/store/consumerapi"
)

type WorkspaceStoreWrapper struct {
	api *workspaceConsumerAPI.StoreAPI
}

func InitWorkspaceWrappers(
	storeWrapper *WorkspaceStoreWrapper,
	runtimeWrapper *WorkspaceRuntimeWrapper,
	aggregateWrapper *WorkspaceAggregateWrapper,
	sources compositionapi.SourceAPI,
	collections compositionapi.CollectionAPI,
	artifacts compositionapi.ArtifactAPI,
	resources compositionapi.ResourceAPI,
) error {
	if storeWrapper == nil ||
		runtimeWrapper == nil ||
		aggregateWrapper == nil {
		return errors.New("workspace wrapper dependencies are incomplete")
	}

	storeAPI, err := workspaceConsumerAPI.NewStoreAPI(
		sources,
		collections,
		artifacts,
		resources,
		workspaceConsumerAPI.DefaultConfig(),
	)
	if err != nil {
		return err
	}

	runtimeAPI, err := workspaceAggregate.NewRuntimeAPI(
		storeAPI.ContextService(),
		storeAPI.SkillAdapter(),
	)
	if err != nil {
		return err
	}

	aggregateAPI, err := workspaceAggregate.NewAggregateAPI(
		storeAPI,
		runtimeAPI,
		storeAPI,
	)
	if err != nil {
		return err
	}

	storeWrapper.api = storeAPI
	runtimeWrapper.api = runtimeAPI
	aggregateWrapper.api = aggregateAPI
	return nil
}

func (w *WorkspaceStoreWrapper) GetWorkspace(
	request *workspaceConsumerAPI.GetWorkspaceRequest,
) (*workspaceConsumerAPI.GetWorkspaceResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.GetWorkspaceResponse, error) {
			return w.api.GetWorkspace(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) ListWorkspaces(
	request *workspaceConsumerAPI.ListWorkspacesRequest,
) (*workspaceConsumerAPI.ListWorkspacesResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.ListWorkspacesResponse, error) {
			return w.api.ListWorkspaces(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) GetWorkspaceCatalog(
	request *workspaceConsumerAPI.GetWorkspaceCatalogRequest,
) (*workspaceConsumerAPI.GetWorkspaceCatalogResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.GetWorkspaceCatalogResponse, error) {
			return w.api.GetWorkspaceCatalog(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) GetWorkspaceArtifact(
	request *workspaceConsumerAPI.GetWorkspaceArtifactRequest,
) (*workspaceConsumerAPI.GetWorkspaceArtifactResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.GetWorkspaceArtifactResponse, error) {
			return w.api.GetWorkspaceArtifact(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) ListWorkspaceArtifacts(
	request *workspaceConsumerAPI.ListWorkspaceArtifactsRequest,
) (*workspaceConsumerAPI.ListWorkspaceArtifactsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.ListWorkspaceArtifactsResponse, error) {
			return w.api.ListWorkspaceArtifacts(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) close() {
	if w == nil {
		return
	}
	w.api = nil
}
