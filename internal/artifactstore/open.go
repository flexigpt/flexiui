package artifactstore

import (
	"context"
	"errors"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/resource"
	rootimpl "github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/system"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
)

// Open creates one complete Artifact Store.
//
// The returned API owns the Store lifecycle and must be closed by its
// application composition owner.
func Open(
	ctx context.Context,
	config Config,
) (*API, error) {
	rootPolicy, err := rootimpl.NewSetRootPolicy(
		append([]basespec.RootID(nil), config.ProtectedRoots...),
		append([]basespec.RootID(nil), config.RetainedRoots...),
	)
	if err != nil {
		return nil, err
	}

	components, err := system.Open(
		ctx,
		system.Config{
			BaseDirectory: config.BaseDirectory,
			ArtifactProviders: append(
				[]providerapi.Provider(nil),
				config.ArtifactProviders...,
			),
			RootMutationPolicy: rootPolicy,
		},
	)
	if err != nil {
		return nil, err
	}

	api, err := newAPI(components)
	if err != nil {
		_ = components.Close()
		return nil, err
	}
	return api, nil
}

func newAPI(components *system.Components) (*API, error) {
	if components == nil ||
		components.Roots == nil ||
		components.Sources == nil {
		return nil, errors.New("artifact store components are required")
	}

	if components.ArtifactReader == nil ||
		components.CollectionReader == nil ||
		components.Refresh == nil ||
		components.SourceRuntime == nil {
		return nil, errors.New(
			"artifact store resource components are required",
		)
	}
	resources, err := resource.NewService(
		components.ArtifactReader,
		components.CollectionReader,
		components.Refresh,
		components.SourceRuntime,
	)
	if err != nil {
		return nil, err
	}
	return &API{
		components: components,
		resources:  resources,
	}, nil
}
