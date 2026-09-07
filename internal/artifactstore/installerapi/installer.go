package installerapi

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/installerapi/topology"
)

// API is the privileged application-composition capability for protected
// Artifact Store topology installation and hydration.
//
// It is intentionally separate from consumerapi.ConsumerAPI.
type API interface {
	EnsureProtectedTopology(
		ctx context.Context,
		declaration topology.Declaration,
	) (topology.Installed, error)

	PrepareTopologyHydrations(
		ctx context.Context,
		desired []topology.Hydration,
	) (map[string]bool, error)

	CommitTopologyHydration(
		ctx context.Context,
		desired topology.Hydration,
	) error
}
