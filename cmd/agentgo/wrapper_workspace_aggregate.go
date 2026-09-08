package main

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/middleware"
	"github.com/flexigpt/flexigpt-app/internal/workspace"
)

type WorkspaceAggregateWrapper struct {
	api *workspace.AggregateAPI
}

func (w *WorkspaceAggregateWrapper) SetWorkspaceArtifactRuntimeDisabled(
	request *workspace.SetWorkspaceArtifactRuntimeDisabledRequest,
) (*workspace.SetWorkspaceArtifactRuntimeDisabledResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspace.SetWorkspaceArtifactRuntimeDisabledResponse, error) {
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
