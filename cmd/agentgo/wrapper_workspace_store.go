package main

import (
	"context"
	"errors"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
	"github.com/flexigpt/flexigpt-app/internal/workspace"
)

type WorkspaceStoreWrapper struct {
	api *workspace.StoreAPI
}

func InitWorkspaceWrappers(
	storeWrapper *WorkspaceStoreWrapper,
	runtimeWrapper *WorkspaceRuntimeWrapper,
	aggregateWrapper *WorkspaceAggregateWrapper,
	sources compositionapi.SourceAPI,
	collections compositionapi.CollectionAPI,
	artifacts compositionapi.ArtifactAPI,
	catalogs compositionapi.CatalogAPI,
	resources compositionapi.ResourceAPI,
) error {
	if storeWrapper == nil ||
		runtimeWrapper == nil ||
		aggregateWrapper == nil {
		return errors.New("workspace wrapper dependencies are incomplete")
	}

	storeAPI, err := workspace.NewStoreAPI(
		sources,
		collections,
		artifacts,
		catalogs,
		resources,
		workspace.DefaultConfig(),
	)
	if err != nil {
		return err
	}

	runtimeAPI, err := workspace.NewRuntimeAPI(
		storeAPI.ContextRuntime(),
		storeAPI.SkillAdapter(),
	)
	if err != nil {
		return err
	}

	aggregateAPI, err := workspace.NewAggregateAPI(
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
	request *workspace.GetWorkspaceRequest,
) (*workspace.GetWorkspaceResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.GetWorkspaceResponse, error) {
			return w.api.GetWorkspace(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) ListWorkspaces(
	request *workspace.ListWorkspacesRequest,
) (*workspace.ListWorkspacesResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.ListWorkspacesResponse, error) {
			return w.api.ListWorkspaces(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) GetWorkspaceCatalog(
	request *workspace.GetWorkspaceCatalogRequest,
) (*workspace.GetWorkspaceCatalogResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.GetWorkspaceCatalogResponse, error) {
			return w.api.GetWorkspaceCatalog(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) GetWorkspaceArtifact(
	request *workspace.GetWorkspaceArtifactRequest,
) (*workspace.GetWorkspaceArtifactResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.GetWorkspaceArtifactResponse, error) {
			return w.api.GetWorkspaceArtifact(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) ListWorkspaceArtifacts(
	request *workspace.ListWorkspaceArtifactsRequest,
) (*workspace.ListWorkspaceArtifactsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.ListWorkspaceArtifactsResponse, error) {
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
