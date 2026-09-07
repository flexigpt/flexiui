package artifact

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
)

type State string

const (
	StateAvailable    State = "available"
	StateMissing      State = "missing"
	StateInvalid      State = "invalid"
	StateIncompatible State = "incompatible"
)

func (s State) Validate(
	resolvedDefinition *cryptoutil.Digest,
) error {
	if resolvedDefinition != nil {
		if err := cryptoutil.ValidateDigest(*resolvedDefinition); err != nil {
			return err
		}
	}

	switch s {
	case StateAvailable, StateIncompatible:
		if resolvedDefinition == nil {
			return fmt.Errorf(
				"%w: artifact state %q requires a resolved definition",
				basespec.ErrInvalid,
				s,
			)
		}

	case StateMissing, StateInvalid:
		if resolvedDefinition != nil {
			return fmt.Errorf(
				"%w: artifact state %q cannot retain a resolved definition",
				basespec.ErrInvalid,
				s,
			)
		}

	default:
		return fmt.Errorf(
			"%w: invalid artifact state %q",
			basespec.ErrInvalid,
			s,
		)
	}
	return nil
}
