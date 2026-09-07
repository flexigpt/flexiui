package artifactstore

import (
	"errors"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/resource"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/system"
)

// NewAPI creates the consumer implementation over composed Store components.
//
// Application composition calls this after provider registration and system
// construction have completed.
func NewAPI(components *system.Components) (*API, error) {
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
