package main

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/middleware"
	"github.com/flexigpt/flexigpt-app/internal/workspace"
)

type WorkspaceRuntimeWrapper struct {
	api *workspace.RuntimeAPI
}

func (w *WorkspaceRuntimeWrapper) ListWorkspaceContexts(
	request *workspace.ListWorkspaceContextsRequest,
) (*workspace.ListWorkspaceContextsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.ListWorkspaceContextsResponse, error) {
			return w.api.ListWorkspaceContexts(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) LoadWorkspaceContexts(
	request *workspace.LoadWorkspaceContextsRequest,
) (*workspace.LoadWorkspaceContextsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.LoadWorkspaceContextsResponse, error) {
			return w.api.LoadWorkspaceContexts(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) ComposeWorkspaceContext(
	request *workspace.ComposeWorkspaceContextRequest,
) (*workspace.ComposeWorkspaceContextResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.ComposeWorkspaceContextResponse, error) {
			return w.api.ComposeWorkspaceContext(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) ListWorkspaceSkills(
	request *workspace.ListWorkspaceSkillsRequest,
) (*workspace.ListWorkspaceSkillsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.ListWorkspaceSkillsResponse, error) {
			return w.api.ListWorkspaceSkills(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) LoadWorkspaceSkills(
	request *workspace.LoadWorkspaceSkillsRequest,
) (*workspace.LoadWorkspaceSkillsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.LoadWorkspaceSkillsResponse, error) {
			return w.api.LoadWorkspaceSkills(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) close() {
	if w == nil {
		return
	}
	w.api = nil
}
