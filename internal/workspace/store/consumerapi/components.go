package consumerapi

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/skill/store/workspaceadapter"
)

type components struct {
	workspaceRootID root.RootID
	service         *Service
	query           *QueryService
	contextService  ContextService
	skillAdapter    *workspaceadapter.Adapter
}

func newComponents(
	sources compositionapi.SourceAPI,
	collections compositionapi.CollectionAPI,
	artifacts compositionapi.ArtifactAPI,
	resources compositionapi.ResourceAPI,
	config Config,
) (*components, error) {
	service, err := NewService(
		collections,
		sources,
		config.WorkspaceRootID,
	)
	if err != nil {
		return nil, err
	}

	query, err := NewQueryService(
		service,
		resources,
	)
	if err != nil {
		return nil, err
	}

	runtimePolicy := config.runtimePolicy()

	contextService, err := newContextService(
		query,
		runtimePolicy,
		config.contextCompositionPolicy(),
	)
	if err != nil {
		return nil, err
	}

	skillAdapter, err := workspaceadapter.NewAdapter(
		query,
		runtimePolicy,
		resources,
	)
	if err != nil {
		return nil, err
	}

	return &components{
		workspaceRootID: config.WorkspaceRootID,
		service:         service,
		query:           query,
		contextService:  contextService,
		skillAdapter:    skillAdapter,
	}, nil
}
