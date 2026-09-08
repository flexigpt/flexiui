package workspace

import (
	"fmt"

	artifactConsumerAPI "github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi"
	"github.com/flexigpt/flexigpt-app/internal/workspace/spec"
)

type Dependencies struct {
	Store *artifactConsumerAPI.API
}

func (d Dependencies) Validate() error {
	if d.Store == nil {
		return fmt.Errorf(
			"%w: Workspace Artifact Store dependencies are incomplete",
			spec.ErrInvalidWorkspace,
		)
	}
	return nil
}
