package consumerapi

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	workspaceRuntime "github.com/flexigpt/flexigpt-app/internal/workspace/runtime"
	"github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/artifactadapter"
)

type Config struct {
	WorkspaceRootID    root.RootID
	ContextComposition workspaceRuntime.CompositionPolicy
	SourceUsePolicy    artifactadapter.SourceUsePolicy
}

func (c Config) normalized() Config {
	output := c
	if output.WorkspaceRootID == "" {
		output.WorkspaceRootID = artifactbuiltin.WorkspaceRootID
	}
	return output
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
		ContextComposition: workspaceRuntime.DefaultCompositionPolicy(),
		SourceUsePolicy:    artifactadapter.NewArtifactRuntimePolicy(),
	}
}
