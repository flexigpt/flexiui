package main

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/middleware"
	workspaceAggregate "github.com/flexigpt/flexigpt-app/internal/workspace/aggregate"
)

type WorkspaceRuntimeWrapper struct {
	api *workspaceAggregate.RuntimeAPI
}

func (w *WorkspaceRuntimeWrapper) ListWorkspaceContexts(
	request *workspaceAggregate.ListWorkspaceContextsRequest,
) (*workspaceAggregate.ListWorkspaceContextsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceAggregate.ListWorkspaceContextsResponse, error) {
			return w.api.ListWorkspaceContexts(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) LoadWorkspaceContexts(
	request *workspaceAggregate.LoadWorkspaceContextsRequest,
) (*workspaceAggregate.LoadWorkspaceContextsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceAggregate.LoadWorkspaceContextsResponse, error) {
			return w.api.LoadWorkspaceContexts(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) ComposeWorkspaceContext(
	request *workspaceAggregate.ComposeWorkspaceContextRequest,
) (*workspaceAggregate.ComposeWorkspaceContextResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceAggregate.ComposeWorkspaceContextResponse, error) {
			return w.api.ComposeWorkspaceContext(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) ListWorkspaceSkills(
	request *workspaceAggregate.ListWorkspaceSkillsRequest,
) (*workspaceAggregate.ListWorkspaceSkillsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceAggregate.ListWorkspaceSkillsResponse, error) {
			return w.api.ListWorkspaceSkills(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) LoadWorkspaceSkills(
	request *workspaceAggregate.LoadWorkspaceSkillsRequest,
) (*workspaceAggregate.LoadWorkspaceSkillsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceAggregate.LoadWorkspaceSkillsResponse, error) {
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
