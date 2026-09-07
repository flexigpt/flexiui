package artifactstore

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"

	artifactConsumerAPIresource "github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi/reqresp/resource"
)

func (a *API) ResolveArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
	options artifactConsumerAPIresource.ResolveOptions,
) (artifactConsumerAPIresource.ResolvedArtifact, error) {
	if err := a.check(ctx); err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}
	if a.resources == nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, basespec.ErrClosed
	}
	value, err := a.resources.ResolveArtifact(
		ctx,
		ref,
		options,
	)
	if err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}
	return value, nil
}

func (a *API) ResolveVerifiedLocalPath(
	ctx context.Context,
	resolved artifactConsumerAPIresource.ResolvedArtifact,
	localLocator basespec.Locator,
) (string, error) {
	if err := a.check(ctx); err != nil {
		return "", err
	}
	if a.resources == nil {
		return "", basespec.ErrClosed
	}
	if err := resolved.Validate(); err != nil {
		return "", err
	}
	return a.resources.ResolveVerifiedLocalPath(
		ctx,
		resolvedArtifactForStore(resolved),
		localLocator,
	)
}

func (a *API) ReadCollectionEntry(
	ctx context.Context,
	ref collection.CollectionRef,
	sourceID basespec.SourceID,
	locator basespec.Locator,
	maximumBytes int64,
) (artifactConsumerAPIresource.VerifiedEntry, error) {
	if err := a.check(ctx); err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}
	if a.resources == nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, basespec.ErrClosed
	}
	value, err := a.resources.ReadCollectionEntry(
		ctx,
		ref,
		sourceID,
		locator,
		maximumBytes,
	)
	if err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}
	return value, nil
}

func (a *API) ResolveSourceLocalPath(
	ctx context.Context,
	rootID root.RootID,
	sourceID basespec.SourceID,
	locator basespec.Locator,
) (string, error) {
	if err := a.check(ctx); err != nil {
		return "", err
	}
	if a.resources == nil {
		return "", basespec.ErrClosed
	}
	return a.resources.ResolveSourceLocalPath(
		ctx,
		rootID,
		sourceID,
		locator,
	)
}

func (a *API) SupportsLocalPath(kind basespec.SourceKind) bool {
	return a != nil &&
		a.resources != nil &&
		a.resources.SupportsLocalPath(kind)
}

func resolvedArtifactForStore(
	value artifactConsumerAPIresource.ResolvedArtifact,
) artifactConsumerAPIresource.ResolvedArtifact {
	value = value.Clone()
	return artifactConsumerAPIresource.ResolvedArtifact{
		Artifact:         value.Artifact.Clone(),
		Collection:       value.Collection.Clone(),
		Definition:       value.Definition.Clone(),
		Occurrence:       value.Occurrence.Clone(),
		Source:           value.Source.Clone(),
		CatalogRevision:  value.CatalogRevision,
		SourceGeneration: value.SourceGeneration,
	}
}
