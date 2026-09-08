package main

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/middleware"
	workspaceAggregate "github.com/flexigpt/flexigpt-app/internal/workspace/aggregate"
)

type WorkspaceAggregateWrapper struct {
	api *workspaceAggregate.AggregateAPI
}

func (w *WorkspaceAggregateWrapper) SetWorkspaceArtifactRuntimeDisabled(
	request *workspaceAggregate.SetWorkspaceArtifactRuntimeDisabledRequest,
) (*workspaceAggregate.SetWorkspaceArtifactRuntimeDisabledResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceAggregate.SetWorkspaceArtifactRuntimeDisabledResponse, error) {
			return w.api.SetWorkspaceArtifactRuntimeDisabled(ctx, request)
		},
	)
}

func (w *WorkspaceAggregateWrapper) close() {
	if w == nil {
		return
	}
	w.api = nil
}
