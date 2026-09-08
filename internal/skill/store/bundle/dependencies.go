package bundle

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
)

type Dependencies struct {
	Sources          compositionapi.SourceAPI
	Collections      compositionapi.CollectionAPI
	Artifacts        compositionapi.ArtifactAPI
	Catalogs         compositionapi.CatalogAPI
	Resources        compositionapi.ResourceAPI
	Schemas          compositionapi.SchemaAPI
	ManagedArtifacts compositionapi.ManagedArtifactAPI
	Protection       compositionapi.ProtectionAPI
}

func (d Dependencies) Validate() error {
	if d.Sources == nil ||
		d.Collections == nil ||
		d.Artifacts == nil ||
		d.Catalogs == nil ||
		d.Resources == nil ||
		d.Schemas == nil ||
		d.ManagedArtifacts == nil ||
		d.Protection == nil {
		return fmt.Errorf(
			"%w: skill bundle Artifact Store dependencies are incomplete",
			basespec.ErrInvalid,
		)
	}
	return nil
}
