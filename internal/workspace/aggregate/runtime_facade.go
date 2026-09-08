package aggregate

import (
	"context"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/skill/store/workspaceadapter"
	workspaceConsumerAPI "github.com/flexigpt/flexigpt-app/internal/workspace/store/consumerapi"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
)

type WorkspaceContextRuntime interface {
	List(
		ctx context.Context,
		workspace collection.CollectionRef,
	) ([]ContextDocument, error)

	Load(
		ctx context.Context,
		workspace collection.CollectionRef,
		artifactRefs []artifact.ArtifactRef,
	) (ContextInspection, error)

	Compose(
		ctx context.Context,
		workspace collection.CollectionRef,
		artifactRefs []artifact.ArtifactRef,
	) (ContextLoadPlan, error)
}

type RuntimeAPI struct {
	contexts WorkspaceContextRuntime
	skills   *workspaceadapter.Adapter
}

func NewRuntimeAPI(
	contexts WorkspaceContextRuntime,
	skills *workspaceadapter.Adapter,
) (*RuntimeAPI, error) {
	if contexts == nil || skills == nil {
		return nil, fmt.Errorf(
			"%w: Workspace Runtime dependencies are incomplete",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}

	return &RuntimeAPI{
		contexts: contexts,
		skills:   skills,
	}, nil
}

func (a *RuntimeAPI) SkillAdapter() *workspaceadapter.Adapter {
	if a == nil {
		return nil
	}
	return a.skills
}

func (a *RuntimeAPI) ListWorkspaceContexts(
	ctx context.Context,
	request *ListWorkspaceContextsRequest,
) (*ListWorkspaceContextsResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Workspace Context list",
	); err != nil {
		return nil, err
	}

	values, err := a.contexts.List(ctx, request.Workspace)
	if err != nil {
		return nil, wrapRuntimeError("list Workspace Contexts", err)
	}

	output := make([]WorkspaceContextView, 0, len(values))
	for _, value := range values {
		output = append(output, workspaceConsumerAPI.ContextViewOf(value))
	}

	return &ListWorkspaceContextsResponse{
		Body: &ListWorkspaceContextsResponseBody{
			Contexts: output,
		},
	}, nil
}

func (a *RuntimeAPI) LoadWorkspaceContexts(
	ctx context.Context,
	request *LoadWorkspaceContextsRequest,
) (*LoadWorkspaceContextsResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Workspace Context load",
	); err != nil {
		return nil, err
	}

	value, err := a.contexts.Load(
		ctx,
		request.Workspace,
		request.Body.Artifacts,
	)
	if err != nil {
		return nil, wrapRuntimeError("load Workspace Contexts", err)
	}

	output := WorkspaceContextInspectionView{
		Workspace:       value.Workspace,
		CatalogRevision: value.CatalogRevision,
		Diagnostics:     diagnostic.Clone(value.Diagnostics),
		Contributions: make(
			[]WorkspaceContextContribution,
			0,
			len(value.Contributions),
		),
	}
	for _, contribution := range value.Contributions {
		output.Contributions = append(
			output.Contributions,
			workspaceConsumerAPI.ContextContributionViewOf(contribution),
		)
	}

	return &LoadWorkspaceContextsResponse{Body: &output}, nil
}

func (a *RuntimeAPI) ComposeWorkspaceContext(
	ctx context.Context,
	request *ComposeWorkspaceContextRequest,
) (*ComposeWorkspaceContextResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Workspace Context composition",
	); err != nil {
		return nil, err
	}

	value, err := a.contexts.Compose(
		ctx,
		request.Workspace,
		request.Body.Artifacts,
	)
	if err != nil {
		return nil, wrapRuntimeError("compose Workspace Context", err)
	}

	output := workspaceConsumerAPI.ContextLoadPlanViewOf(value)
	return &ComposeWorkspaceContextResponse{Body: &output}, nil
}

func (a *RuntimeAPI) ListWorkspaceSkills(
	ctx context.Context,
	request *ListWorkspaceSkillsRequest,
) (*ListWorkspaceSkillsResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Workspace Skill list",
	); err != nil {
		return nil, err
	}

	values, err := a.skills.List(ctx, request.Workspace)
	if err != nil {
		return nil, wrapRuntimeError("list Workspace Skills", err)
	}

	output := make([]WorkspaceSkillView, 0, len(values))
	for _, value := range values {
		output = append(output, workspaceConsumerAPI.WorkspaceSkillViewOf(value))
	}

	return &ListWorkspaceSkillsResponse{
		Body: &ListWorkspaceSkillsResponseBody{
			Skills: output,
		},
	}, nil
}

func (a *RuntimeAPI) LoadWorkspaceSkills(
	ctx context.Context,
	request *LoadWorkspaceSkillsRequest,
) (*LoadWorkspaceSkillsResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Workspace Skill load",
	); err != nil {
		return nil, err
	}

	value, err := a.skills.Load(
		ctx,
		request.Workspace,
		request.Body.Artifacts,
	)
	if err != nil {
		return nil, wrapRuntimeError("load Workspace Skills", err)
	}

	output := workspaceConsumerAPI.WorkspaceSkillLoadViewOf(value)
	return &LoadWorkspaceSkillsResponse{Body: &output}, nil
}

func wrapRuntimeError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("workspace runtime %s: %w", operation, err)
}
