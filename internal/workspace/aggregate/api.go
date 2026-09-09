package aggregate

import (
	"context"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	workspaceConsumerAPI "github.com/flexigpt/flexigpt-app/internal/workspace/store/consumerapi"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
)

type WorkspaceStoreReader interface {
	GetWorkspace(
		ctx context.Context,
		request *workspaceConsumerAPI.GetWorkspaceRequest,
	) (*workspaceConsumerAPI.GetWorkspaceResponse, error)
}

type WorkspaceRuntimeReader interface {
	ComposeWorkspaceContext(
		ctx context.Context,
		request *workspaceConsumerAPI.ComposeWorkspaceContextRequest,
	) (*workspaceConsumerAPI.ComposeWorkspaceContextResponse, error)

	LoadWorkspaceSkills(
		ctx context.Context,
		request *workspaceConsumerAPI.LoadWorkspaceSkillsRequest,
	) (*workspaceConsumerAPI.LoadWorkspaceSkillsResponse, error)
}

type WorkspaceArtifactSettingsStore interface {
	SetArtifactRuntimeDisabled(
		ctx context.Context,
		workspace workspaceDomain.WorkspaceRef,
		ref artifact.ArtifactRef,
		expectedRevision uint64,
		runtimeDisabled bool,
	) (workspaceConsumerAPI.WorkspaceArtifactView, error)
}

type AggregateAPI struct {
	store    WorkspaceStoreReader
	runtime  WorkspaceRuntimeReader
	settings WorkspaceArtifactSettingsStore
}

func NewAggregateAPI(
	store WorkspaceStoreReader,
	runtime WorkspaceRuntimeReader,
	settings WorkspaceArtifactSettingsStore,
) (*AggregateAPI, error) {
	if store == nil || runtime == nil || settings == nil {
		return nil, fmt.Errorf(
			"%w: Workspace Aggregate dependencies are incomplete",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}

	return &AggregateAPI{
		store:    store,
		runtime:  runtime,
		settings: settings,
	}, nil
}

func (a *AggregateAPI) GetWorkspace(
	ctx context.Context,
	request *workspaceConsumerAPI.GetWorkspaceRequest,
) (*workspaceConsumerAPI.GetWorkspaceResponse, error) {
	return a.store.GetWorkspace(ctx, request)
}

func (a *AggregateAPI) ComposeWorkspaceContext(
	ctx context.Context,
	request *workspaceConsumerAPI.ComposeWorkspaceContextRequest,
) (*workspaceConsumerAPI.ComposeWorkspaceContextResponse, error) {
	return a.runtime.ComposeWorkspaceContext(ctx, request)
}

func (a *AggregateAPI) LoadWorkspaceSkills(
	ctx context.Context,
	request *workspaceConsumerAPI.LoadWorkspaceSkillsRequest,
) (*workspaceConsumerAPI.LoadWorkspaceSkillsResponse, error) {
	return a.runtime.LoadWorkspaceSkills(ctx, request)
}

func (a *AggregateAPI) SetWorkspaceArtifactRuntimeDisabled(
	ctx context.Context,
	request *workspaceConsumerAPI.SetWorkspaceArtifactRuntimeDisabledRequest,
) (*workspaceConsumerAPI.SetWorkspaceArtifactRuntimeDisabledResponse, error) {
	if err := workspaceConsumerAPI.RequireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Workspace Artifact runtime settings update",
	); err != nil {
		return nil, err
	}

	value, err := a.settings.SetArtifactRuntimeDisabled(
		ctx,
		request.Workspace,
		request.Artifact,
		request.Body.ExpectedRevision,
		request.Body.RuntimeDisabled,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"workspace aggregate set Artifact runtime settings: %w",
			err,
		)
	}

	return &workspaceConsumerAPI.SetWorkspaceArtifactRuntimeDisabledResponse{
		Body: &value,
	}, nil
}
