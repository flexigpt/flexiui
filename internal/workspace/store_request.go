package workspace

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/workspace/spec"
)

// requireRequestBody performs transport-shape checks only.
//
// Domain validation belongs to Workspace services. Store validation and
// lifecycle enforcement belong to Artifact Store.
func requireRequestBody[T any](
	request *T,
	bodyPresent bool,
	requireBody bool,
	subject string,
) error {
	if request == nil {
		return fmt.Errorf(
			"%w: %s request is required",
			spec.ErrInvalidWorkspace,
			subject,
		)
	}
	if requireBody && !bodyPresent {
		return fmt.Errorf(
			"%w: %s request body is required",
			spec.ErrInvalidWorkspace,
			subject,
		)
	}
	return nil
}
