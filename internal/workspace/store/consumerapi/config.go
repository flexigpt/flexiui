package consumerapi

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	workspaceRuntime "github.com/flexigpt/flexigpt-app/internal/workspace/runtime"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
	"github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/artifactadapter"
	"github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/support"
)

type Config struct {
	WorkspaceRootID    root.RootID
	Supports           []workspaceDomain.ArtifactSupport
	ContextComposition workspaceRuntime.CompositionPolicy
	SourceUsePolicy    artifactadapter.SourceUsePolicy
}

func (c Config) normalized() Config {
	output := c
	if len(output.Supports) == 0 {
		output.Supports = support.DefaultArtifactSupports()
	}
	if output.WorkspaceRootID == "" {
		output.WorkspaceRootID = artifactbuiltin.WorkspaceRootID
	}
	return output
}

func (c Config) normalizedSupports() ([]workspaceDomain.ArtifactSupport, error) {
	if len(c.Supports) == 0 {
		return nil, fmt.Errorf(
			"%w: workspace artifact support is required",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}

	output := make([]workspaceDomain.ArtifactSupport, 0, len(c.Supports))
	seenKinds := make(map[artifact.ArtifactKind]struct{}, len(c.Supports))

	for _, support := range c.Supports {
		if err := support.Validate(); err != nil {
			return nil, err
		}
		if _, duplicate := seenKinds[support.Kind]; duplicate {
			return nil, fmt.Errorf(
				"%w: duplicate workspace artifact kind %q",
				workspaceDomain.ErrInvalidWorkspace,
				support.Kind,
			)
		}
		seenKinds[support.Kind] = struct{}{}
		output = append(output, support)
	}
	return output, nil
}

func (c Config) runtimePolicy() artifactadapter.SourceUsePolicy {
	if c.SourceUsePolicy != nil {
		return c.SourceUsePolicy
	}

	return artifactadapter.NewArtifactRuntimePolicy()
}

func (c Config) contextCompositionPolicy() workspaceRuntime.CompositionPolicy {
	return c.ContextComposition.Normalized()
}

func DefaultConfig() Config {
	return Config{
		WorkspaceRootID:    artifactbuiltin.WorkspaceRootID,
		Supports:           support.DefaultArtifactSupports(),
		ContextComposition: workspaceRuntime.DefaultCompositionPolicy(),
		SourceUsePolicy:    artifactadapter.NewArtifactRuntimePolicy(),
	}
}
