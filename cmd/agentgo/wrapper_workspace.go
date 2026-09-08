package main

import (
	"context"
	"errors"

	artifactConsumerAPI "github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
	"github.com/flexigpt/flexigpt-app/internal/workspace"
)

type WorkspaceWrapper struct {
	api *workspace.API
}

func InitWorkspaceWrapper(
	wrapper *WorkspaceWrapper,
	store *artifactConsumerAPI.API,
) error {
	if wrapper == nil {
		return errors.New("workspace wrapper is nil")
	}
	if store == nil {
		return errors.New("artifact store API is nil")
	}

	api, err := workspace.New(
		workspace.Dependencies{Store: store},
		workspace.DefaultConfig(),
	)
	if err != nil {
		return err
	}
	wrapper.api = api
	return nil
}

func (w *WorkspaceWrapper) GetWorkspace(
	request *workspace.GetWorkspaceRequest,
) (*workspace.GetWorkspaceResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.GetWorkspaceResponse, error) {
		return w.api.GetWorkspace(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) ListWorkspaces(
	request *workspace.ListWorkspacesRequest,
) (*workspace.ListWorkspacesResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.ListWorkspacesResponse, error) {
		return w.api.ListWorkspaces(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) GetWorkspaceCatalog(
	request *workspace.GetWorkspaceCatalogRequest,
) (*workspace.GetWorkspaceCatalogResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.GetWorkspaceCatalogResponse, error) {
		return w.api.GetWorkspaceCatalog(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) GetWorkspaceArtifact(
	request *workspace.GetWorkspaceArtifactRequest,
) (*workspace.GetWorkspaceArtifactResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.GetWorkspaceArtifactResponse, error) {
		return w.api.GetWorkspaceArtifact(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) ListWorkspaceArtifacts(
	request *workspace.ListWorkspaceArtifactsRequest,
) (*workspace.ListWorkspaceArtifactsResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.ListWorkspaceArtifactsResponse, error) {
		return w.api.ListWorkspaceArtifacts(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) ListWorkspaceContexts(
	request *workspace.ListWorkspaceContextsRequest,
) (*workspace.ListWorkspaceContextsResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.ListWorkspaceContextsResponse, error) {
		return w.api.ListWorkspaceContexts(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) LoadWorkspaceContexts(
	request *workspace.LoadWorkspaceContextsRequest,
) (*workspace.LoadWorkspaceContextsResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.LoadWorkspaceContextsResponse, error) {
		return w.api.LoadWorkspaceContexts(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) ComposeWorkspaceContext(
	request *workspace.ComposeWorkspaceContextRequest,
) (*workspace.ComposeWorkspaceContextResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.ComposeWorkspaceContextResponse, error) {
		return w.api.ComposeWorkspaceContext(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) ListWorkspaceSkills(
	request *workspace.ListWorkspaceSkillsRequest,
) (*workspace.ListWorkspaceSkillsResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.ListWorkspaceSkillsResponse, error) {
		return w.api.ListWorkspaceSkills(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) LoadWorkspaceSkills(
	request *workspace.LoadWorkspaceSkillsRequest,
) (*workspace.LoadWorkspaceSkillsResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.LoadWorkspaceSkillsResponse, error) {
		return w.api.LoadWorkspaceSkills(context.Background(), request)
	})
}

func (w *WorkspaceWrapper) SetWorkspaceArtifactRuntimeDisabled(
	request *workspace.SetWorkspaceArtifactRuntimeDisabledRequest,
) (*workspace.SetWorkspaceArtifactRuntimeDisabledResponse, error) {
	return middleware.WithRecoveryResp(func() (*workspace.SetWorkspaceArtifactRuntimeDisabledResponse, error) {
		ctx := context.Background()
		response, err := w.api.SetWorkspaceArtifactRuntimeDisabled(ctx, request)
		if err != nil {
			return nil, err
		}

		return response, nil
	})
}

func (w *WorkspaceWrapper) close() {
	if w == nil {
		return
	}
	api := w.api
	w.api = nil
	if api != nil {
		_ = api.Close()
	}
}
