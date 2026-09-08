package consumerapi

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/skill/store/workspaceadapter"
)

type components struct {
	workspaceRootID root.RootID
	service         *Service
	query           *QueryService
	contextAdapter  *ContextAdapter
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
		artifacts,
		catalogs,
		supports...,
	)
	if err != nil {
		return nil, err
	}

	runtimePolicy := config.runtimePolicy()

	contextAdapter, err := NewContextAdapter(
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
