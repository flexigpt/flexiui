package consumerapi

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi/reqresp/resource"
)

// ResourceResolver is the consumer-facing Artifact Store capability for
// current resource resolution and controlled Source access.
type ResourceResolver interface {
	ResolveArtifact(
		ctx context.Context,
		ref artifact.ArtifactRef,
		options resource.ResolveOptions,
	) (resource.ResolvedArtifact, error)

	ResolveVerifiedLocalPath(
		ctx context.Context,
		resolved resource.ResolvedArtifact,
		localLocator basespec.Locator,
	) (string, error)

	ReadCollectionEntry(
		ctx context.Context,
		ref collection.CollectionRef,
		sourceID basespec.SourceID,
		locator basespec.Locator,
		maximumBytes int64,
	) (resource.VerifiedEntry, error)

	ResolveSourceLocalPath(
		ctx context.Context,
		rootID root.RootID,
		sourceID basespec.SourceID,
		locator basespec.Locator,
	) (string, error)

	SupportsLocalPath(kind basespec.SourceKind) bool
}
