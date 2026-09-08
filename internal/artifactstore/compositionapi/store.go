package compositionapi

import (
	"context"
	"fmt"
	"sync"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/installerapi/topology"
	rootimpl "github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/system"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
)

// Store owns one composed Artifact Store instance.
//
// Consumer returns the application-facing contract. Privileged topology
// operations and Close remain on this composition owner.
type Store struct {
	mu         sync.RWMutex
	components *system.Components
	consumer   *consumerapi.API
}

// Open constructs Artifact Store from already initialized provider plugins.
func Open(
	ctx context.Context,
	config Config,
) (*Store, error) {
	if ctx == nil {
		return nil, fmt.Errorf(
			"%w: Artifact Store composition context is nil",
			basespec.ErrInvalid,
		)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	retainedRootIDs := make(
		[]root.RootID,
		0,
		len(config.RetainedRoots),
	)
	for index, draft := range config.RetainedRoots {
		if err := draft.ID.Validate(); err != nil {
			return nil, fmt.Errorf(
				"retained Root declaration %d: %w",
				index,
				err,
			)
		}
		retainedRootIDs = append(retainedRootIDs, draft.ID)
	}

	rootPolicy, err := rootimpl.NewSetRootPolicy(
		append([]root.RootID(nil), config.ProtectedRootIDs...),
		retainedRootIDs,
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
				config.Providers...,
			),
			RootMutationPolicy: rootPolicy,
		},
	)
	if err != nil {
		return nil, err
	}

	consumer, err := consumerapi.New(components)
	if err != nil {
		_ = components.Close()
		return nil, err
	}

	output := &Store{
		consumer:   consumer,
		components: components,
	}
	for _, draft := range config.RetainedRoots {
		if _, err := consumer.CreateRoot(ctx, draft); err != nil {
			_ = output.Close()
			return nil, fmt.Errorf(
				"ensure retained application Root %q: %w",
				draft.ID,
				err,
			)
		}
	}
	return output, nil
}

// Consumer returns the direct Artifact Store consumer implementation.
func (s *Store) Consumer() *consumerapi.API {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.consumer
}

func (s *Store) EnsureProtectedTopology(
	ctx context.Context,
	declaration topology.Declaration,
) (topology.Installed, error) {
	if s == nil {
		return topology.Installed{}, basespec.ErrClosed
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.components == nil {
		return topology.Installed{}, basespec.ErrClosed
	}
	return s.components.EnsureProtectedTopology(ctx, declaration)
}

func (s *Store) PrepareTopologyHydrations(
	ctx context.Context,
	desired []topology.Hydration,
) (map[string]bool, error) {
	if s == nil {
		return nil, basespec.ErrClosed
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.components == nil {
		return nil, basespec.ErrClosed
	}
	return s.components.PrepareTopologyHydrations(ctx, desired)
}

func (s *Store) CommitTopologyHydration(
	ctx context.Context,
	desired topology.Hydration,
) error {
	if s == nil {
		return basespec.ErrClosed
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.components == nil {
		return basespec.ErrClosed
	}
	return s.components.CommitTopologyHydration(ctx, desired)
}

// Close releases all Artifact Store resources owned by this composition.
func (s *Store) Close() error {
	if s == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	consumer := s.consumer
	s.consumer = nil
	s.components = nil
	if consumer == nil {
		return nil
	}
	return consumer.Close()
}
