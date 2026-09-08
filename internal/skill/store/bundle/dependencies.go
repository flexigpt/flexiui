package bundle

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	artifactConsumerAPI "github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi"
)

type Dependencies struct {
	Store *artifactConsumerAPI.API
}

func (d Dependencies) Validate() error {
	if d.Store == nil {
		return fmt.Errorf(
			"%w: skill bundle Artifact Store dependencies are incomplete",
			basespec.ErrInvalid,
		)
	}
	return nil
}
