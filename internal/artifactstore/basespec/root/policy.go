package root

import "github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"

// RootPolicy identifies Roots whose ordinary mutations are prohibited. The
// application composition owns the concrete policy and the protected-root
// declaration. Artifact Store does not know which application feature owns it.
type RootPolicy interface {
	IsProtectedRoot(r basespec.RootID) bool
}

// RootDeletionPolicy is an optional lifecycle policy for Roots that must
// remain present while still allowing ordinary mutation of their descendants.
//
// A retained Root is intentionally different from a protected topology Root:
// protected Roots reject ordinary Source, Collection, and Artifact mutations;
// retained Roots reject only Root retirement and purge.
type RootDeletionPolicy interface {
	IsRootDeletionProtected(rootID basespec.RootID) bool
}
