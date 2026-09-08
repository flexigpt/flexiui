package main

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/middleware"
	workspaceAggregate "github.com/flexigpt/flexigpt-app/internal/workspace/aggregate"
	workspaceConsumerAPI "github.com/flexigpt/flexigpt-app/internal/workspace/store/consumerapi"
)

type WorkspaceRuntimeWrapper struct {
	api *workspaceAggregate.RuntimeAPI
}

func (w *WorkspaceRuntimeWrapper) ListWorkspaceContexts(
	request *workspaceConsumerAPI.ListWorkspaceContextsRequest,
) (*workspaceConsumerAPI.ListWorkspaceContextsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.ListWorkspaceContextsResponse, error) {
			return w.api.ListWorkspaceContexts(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) LoadWorkspaceContexts(
	request *workspaceConsumerAPI.LoadWorkspaceContextsRequest,
) (*workspaceConsumerAPI.LoadWorkspaceContextsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.LoadWorkspaceContextsResponse, error) {
			return w.api.LoadWorkspaceContexts(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) ComposeWorkspaceContext(
	request *workspaceConsumerAPI.ComposeWorkspaceContextRequest,
) (*workspaceConsumerAPI.ComposeWorkspaceContextResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.ComposeWorkspaceContextResponse, error) {
			return w.api.ComposeWorkspaceContext(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) ListWorkspaceSkills(
	request *workspaceConsumerAPI.ListWorkspaceSkillsRequest,
) (*workspaceConsumerAPI.ListWorkspaceSkillsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.ListWorkspaceSkillsResponse, error) {
			return w.api.ListWorkspaceSkills(ctx, request)
		},
	)
}

func (w *WorkspaceRuntimeWrapper) LoadWorkspaceSkills(
	request *workspaceConsumerAPI.LoadWorkspaceSkillsRequest,
) (*workspaceConsumerAPI.LoadWorkspaceSkillsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.LoadWorkspaceSkillsResponse, error) {
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
