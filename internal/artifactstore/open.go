package artifactstore

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
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
	config OpenConfig,
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
