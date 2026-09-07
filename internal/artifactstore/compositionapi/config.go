package compositionapi

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
)

// Config contains application-composition inputs for one Artifact Store.
//
// Providers are fully constructed before Open is called. RetainedRoots are
// both lifecycle-policy declarations and initial generic Root declarations.
type Config struct {
	BaseDirectory string

	Providers []providerapi.Provider

	ProtectedRootIDs []root.RootID
	RetainedRoots    []root.RootDraft
}
