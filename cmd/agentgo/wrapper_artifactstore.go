package main

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
	mcpSchemaadapter "github.com/flexigpt/flexigpt-app/internal/mcp/store/schemaadapter"
	skillBundle "github.com/flexigpt/flexigpt-app/internal/skill/store/bundle"
	"github.com/flexigpt/flexigpt-app/internal/workspace"
)

func composeArtifactStore(
	ctx context.Context,
	baseDirectory string,
) (*compositionapi.Store, error) {
	if err := artifactbuiltin.ValidateApplicationTopology(); err != nil {
		return nil, err
	}

	workspaceConfig := workspace.DefaultConfig()
	workspaceProvider, err := workspace.NewProvider(
		workspaceConfig.ProviderConfig(),
	)
	if err != nil {
		return nil, err
	}

	skillProvider, err := skillBundle.NewProvider()
	if err != nil {
		return nil, err
	}

	mcpProvider, err := mcpSchemaadapter.NewProvider()
	if err != nil {
		return nil, err
	}

	providers := []providerapi.Provider{
		workspaceProvider,
		skillProvider,
		mcpProvider,
	}

	return compositionapi.Open(
		ctx,
		compositionapi.Config{
			BaseDirectory:    baseDirectory,
			Providers:        providers,
			ProtectedRootIDs: artifactbuiltin.ProtectedRootIDs(),
			RetainedRoots:    artifactbuiltin.RetainedRootDrafts(),
		},
	)
}
