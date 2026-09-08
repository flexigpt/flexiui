package runtime

import (
	"errors"
	"fmt"
	"strings"
)

type RuntimeUse string

const (
	RuntimeUseContextPrompt RuntimeUse = "context-prompt"
	RuntimeUseSkill         RuntimeUse = "skill"
)

type RuntimeDisposition string

const (
	RuntimeAllowed     RuntimeDisposition = "allowed"
	RuntimeDenied      RuntimeDisposition = "denied"
	RuntimeUnavailable RuntimeDisposition = "unavailable"
)

type RuntimeDecision struct {
	Disposition RuntimeDisposition
	Code        string
	Message     string
}

func (d RuntimeDecision) Validate() error {
	switch d.Disposition {
	case RuntimeAllowed:
		if d.Code != "" || d.Message != "" {
			return errors.New("allowed runtime decision cannot contain denial details")
		}
		return nil

	case RuntimeDenied, RuntimeUnavailable:
		if strings.TrimSpace(d.Code) == "" || len(d.Code) > 256 {
			return errors.New("runtime policy diagnostic code is invalid")
		}
		if strings.TrimSpace(d.Message) == "" || len(d.Message) > 4096 {
			return errors.New("runtime policy diagnostic message is invalid")
		}
		return nil

	default:
		return fmt.Errorf(
			"unsupported runtime disposition %q",
			d.Disposition,
		)
	}
}
