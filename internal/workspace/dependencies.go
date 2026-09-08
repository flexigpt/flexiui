package workspace

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/workspace/spec"
)

type Dependencies struct {
	Sources     compositionapi.SourceAPI
	Collections compositionapi.CollectionAPI
	Artifacts   compositionapi.ArtifactAPI
	Catalogs    compositionapi.CatalogAPI
	Resources   compositionapi.ResourceAPI
}

func (d Dependencies) Validate() error {
	if d.Sources == nil ||
		d.Collections == nil ||
		d.Artifacts == nil ||
		d.Catalogs == nil ||
		d.Resources == nil {
		return fmt.Errorf(
			"%w: Workspace Artifact Store dependencies are incomplete",
			spec.ErrInvalidWorkspace,
		)
	}
	return nil
}
