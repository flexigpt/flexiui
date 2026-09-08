package main

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/middleware"
	workspaceAggregate "github.com/flexigpt/flexigpt-app/internal/workspace/aggregate"
	workspaceConsumerAPI "github.com/flexigpt/flexigpt-app/internal/workspace/store/consumerapi"
)

type WorkspaceAggregateWrapper struct {
	api *workspaceAggregate.AggregateAPI
}

func (w *WorkspaceAggregateWrapper) SetWorkspaceArtifactRuntimeDisabled(
	request *workspaceConsumerAPI.SetWorkspaceArtifactRuntimeDisabledRequest,
) (*workspaceConsumerAPI.SetWorkspaceArtifactRuntimeDisabledResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.SetWorkspaceArtifactRuntimeDisabledResponse, error) {
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
