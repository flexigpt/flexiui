package consumerapi

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/resource"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

func (a *API) ResolveArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
	options resource.ResolveOptions,
) (resource.ResolvedArtifact, error) {
	if err := a.check(ctx); err != nil {
		return resource.ResolvedArtifact{}, err
	}
	if a.resources == nil {
		return resource.ResolvedArtifact{}, basespec.ErrClosed
	}
	value, err := a.resources.ResolveArtifact(
		ctx,
		ref,
		options,
	)
	if err != nil {
		return resource.ResolvedArtifact{}, err
	}
	return value, nil
}

func (a *API) ResolveVerifiedLocalPath(
	ctx context.Context,
	resolved resource.ResolvedArtifact,
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
	sourceID source.SourceID,
	locator basespec.Locator,
	maximumBytes int64,
) (resource.VerifiedEntry, error) {
	if err := a.check(ctx); err != nil {
		return resource.VerifiedEntry{}, err
	}
	if a.resources == nil {
		return resource.VerifiedEntry{}, basespec.ErrClosed
	}
	value, err := a.resources.ReadCollectionEntry(
		ctx,
		ref,
		sourceID,
		locator,
		maximumBytes,
	)
	if err != nil {
		return resource.VerifiedEntry{}, err
	}
	return value, nil
}

func (a *API) ResolveSourceLocalPath(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
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

func (a *API) SupportsLocalPath(kind source.SourceKind) bool {
	return a != nil &&
		a.resources != nil &&
		a.resources.SupportsLocalPath(kind)
}

func resolvedArtifactForStore(
	value resource.ResolvedArtifact,
) resource.ResolvedArtifact {
	value = value.Clone()
	return resource.ResolvedArtifact{
		Artifact:         value.Artifact.Clone(),
		Collection:       value.Collection.Clone(),
		Definition:       value.Definition.Clone(),
		Occurrence:       value.Occurrence.Clone(),
		Source:           value.Source.Clone(),
		CatalogRevision:  value.CatalogRevision,
		SourceGeneration: value.SourceGeneration,
	}
}
