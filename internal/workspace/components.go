package workspace

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/skill/store/workspaceadapter"
	"github.com/flexigpt/flexigpt-app/internal/workspace/artifactadapter"
	"github.com/flexigpt/flexigpt-app/internal/workspace/contextadapter"
)

type components struct {
	workspaceRootID root.RootID
	service         *artifactadapter.Service
	query           *artifactadapter.QueryService
	contextAdapter  *contextadapter.Adapter
	skillAdapter    *workspaceadapter.Adapter
	supportedKinds  map[artifact.ArtifactKind]struct{}
}

func newComponents(
	sources compositionapi.SourceAPI,
	collections compositionapi.CollectionAPI,
	artifacts compositionapi.ArtifactAPI,
	catalogs compositionapi.CatalogAPI,
	resources compositionapi.ResourceAPI,
	config Config,
) (*components, error) {
	supports, err := config.normalizedSupports()
	if err != nil {
		return nil, err
	}

	service, err := artifactadapter.NewService(
		collections,
		sources,
		config.WorkspaceRootID,
	)
	if err != nil {
		return nil, err
	}

	query, err := artifactadapter.NewQueryService(
		service,
		artifacts,
		catalogs,
		supports...,
	)
	if err != nil {
		return nil, err
	}

	runtimePolicy := config.runtimePolicy()

	contextAdapter, err := contextadapter.NewAdapter(
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

	supportedKinds := make(
		map[artifact.ArtifactKind]struct{},
		len(supports),
	)
	for _, support := range supports {
		supportedKinds[support.Kind] = struct{}{}
	}

	return &components{
		workspaceRootID: config.WorkspaceRootID,
		service:         service,
		query:           query,
		contextAdapter:  contextAdapter,
		skillAdapter:    skillAdapter,
		supportedKinds:  supportedKinds,
	}, nil
}
