package artifact

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
)

type PublishArtifactRequest struct {
	Artifact           Artifact                         `json:"artifact"`
	ExpectedDefinition cryptoutil.Digest                `json:"expectedDefinition"`
	Package            source.ManagedPackagePublication `json:"package"`
	AllowProtected     bool                             `json:"allowProtected"`
}

type PublishArtifactResult struct {
	Artifact   Artifact       `json:"artifact"`
	Source     source.Summary `json:"source"`
	Generation string         `json:"generation"`
	Refreshed  bool           `json:"refreshed"`
}

type RemoveArtifactRequest struct {
	Artifact       Artifact                     `json:"artifact"`
	Package        source.ManagedPackageAddress `json:"package"`
	AllowProtected bool                         `json:"allowProtected"`
}
