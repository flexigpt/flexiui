package basespec

import (
	"context"
	"fmt"
)

type privilegedInstallerContextKey struct{}

// WithPrivilegedInstaller grants the narrow application-composition capability
// used by trusted protected-topology installers and update paths. It must never
// be used by transport wrappers.
func WithPrivilegedInstaller(ctx context.Context) context.Context {
	return context.WithValue(ctx, privilegedInstallerContextKey{}, true)
}

func RequirePrivilegedInstaller(ctx context.Context) error {
	if ctx == nil || !IsPrivilegedInstaller(ctx) {
		return fmt.Errorf(
			"%w: protected topology installation requires trusted installer access",
			ErrProtected,
		)
	}
	return nil
}

func IsPrivilegedInstaller(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	value, _ := ctx.Value(privilegedInstallerContextKey{}).(bool)
	return value
}
