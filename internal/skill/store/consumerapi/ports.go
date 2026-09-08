package consumerapi

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
)

type BundleReader interface {
	GetBundle(
		ctx context.Context,
		ref collection.CollectionRef,
	) (Bundle, error)
}

type BuiltinStore interface {
	ListBundles(
		ctx context.Context,
		rootID root.RootID,
	) ([]Bundle, error)

	ListSkills(
		ctx context.Context,
		ref collection.CollectionRef,
	) ([]artifact.Artifact, error)

	EnsureBuiltInBundleTopology(
		ctx context.Context,
		request BuiltInBundleTopology,
	) (Bundle, error)

	InstallBuiltInCollection(
		ctx context.Context,
		request BuiltInCollectionInstallRequest,
	) ([]CreateManagedSkillResponse, error)

	EnsureBuiltInBundleCurrent(
		ctx context.Context,
		ref collection.CollectionRef,
	) error
}
