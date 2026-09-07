package resource

import (
	"context"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	artifactConsumerAPIresource "github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi/reqresp/resource"
	artifactimpl "github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/artifact"
	collectionimpl "github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/collection"
	sourceimpl "github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/source"
)

type catalogReader interface {
	CurrentCatalog(
		ctx context.Context,
		ref collection.CollectionRef,
	) (catalog.Snapshot, error)
}

type Service struct {
	artifacts   artifactimpl.Reader
	collections collectionimpl.Reader
	catalogs    catalogReader
	sources     sourceimpl.Runtime
}

func NewService(
	artifacts artifactimpl.Reader,
	collections collectionimpl.Reader,
	catalogs catalogReader,
	sources sourceimpl.Runtime,
) (*Service, error) {
	if artifacts == nil ||
		collections == nil ||
		catalogs == nil ||
		sources == nil {
		return nil, fmt.Errorf(
			"%w: Artifact resource resolver dependencies are incomplete",
			basespec.ErrInvalid,
		)
	}
	return &Service{
		artifacts:   artifacts,
		collections: collections,
		catalogs:    catalogs,
		sources:     sources,
	}, nil
}

// ResolveArtifact verifies the complete current resource chain:
//
// Artifact -> Collection -> current provider plan -> Catalog occurrence ->
// Definition -> Source revision -> optional exact source bytes.
func (s *Service) ResolveArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
	options artifactConsumerAPIresource.ResolveOptions,
) (artifactConsumerAPIresource.ResolvedArtifact, error) {
	if err := validateContext(ctx, "Artifact resolution"); err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}
	if s == nil ||
		s.artifacts == nil ||
		s.collections == nil ||
		s.catalogs == nil ||
		s.sources == nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, basespec.ErrClosed
	}
	if err := ref.Validate(); err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}

	record, err := s.artifacts.Get(ctx, ref)
	if err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}
	if record.State != artifact.StateAvailable ||
		record.ResolvedDefinition == nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, fmt.Errorf(
			"%w: Artifact %q is not currently available",
			basespec.ErrReferenceUnresolved,
			ref.ArtifactID,
		)
	}

	collectionRef := collection.CollectionRef{
		RootID:       record.RootID,
		CollectionID: record.CollectionID,
	}
	collectionValue, err := s.collections.Get(ctx, collectionRef)
	if err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}
	if collectionValue.Ref() != collectionRef {
		return artifactConsumerAPIresource.ResolvedArtifact{}, fmt.Errorf(
			"%w: Collection reader returned another Collection",
			basespec.ErrInvalid,
		)
	}

	snapshot, err := s.catalogs.CurrentCatalog(ctx, collectionRef)
	if err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}

	var occurrence *catalog.Occurrence
	for index := range snapshot.Occurrences {
		current := snapshot.Occurrences[index]
		if current.Key.SourceID != record.Binding.SourceID ||
			current.Key.Locator != record.Binding.Locator ||
			current.Key.SubresourceLocator !=
				record.Binding.SubresourceLocator {
			continue
		}
		value := current.Clone()
		occurrence = &value
		break
	}
	if occurrence == nil ||
		occurrence.State != catalog.OccurrenceValid ||
		occurrence.Kind != record.Kind ||
		occurrence.DefinitionDigest == nil ||
		occurrence.SourceContentDigest == nil ||
		*occurrence.DefinitionDigest != *record.ResolvedDefinition {
		return artifactConsumerAPIresource.ResolvedArtifact{}, fmt.Errorf(
			"%w: Artifact %q does not match its current Catalog occurrence",
			basespec.ErrCatalogStale,
			record.ID,
		)
	}

	definitionValue, err := snapshot.DefinitionForOccurrence(
		occurrence.Key,
	)
	if err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}
	if definitionValue.Kind != record.Kind ||
		definitionValue.Digest != *record.ResolvedDefinition {
		return artifactConsumerAPIresource.ResolvedArtifact{}, fmt.Errorf(
			"%w: Artifact %q definition does not match current state",
			basespec.ErrDigestMismatch,
			record.ID,
		)
	}

	sourceRevision := snapshot.SourceRevisions[record.Binding.SourceID]
	sourceGeneration := snapshot.SourceGenerations[record.Binding.SourceID]
	if sourceRevision == 0 || sourceGeneration == "" {
		return artifactConsumerAPIresource.ResolvedArtifact{}, fmt.Errorf(
			"%w: Artifact Source has no current Catalog state",
			basespec.ErrCatalogStale,
		)
	}

	sourceValue, err := s.sources.Get(
		ctx,
		record.RootID,
		record.Binding.SourceID,
	)
	if err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}
	if sourceValue.Revision != sourceRevision {
		return artifactConsumerAPIresource.ResolvedArtifact{}, fmt.Errorf(
			"%w: Artifact Source changed after Catalog publication",
			basespec.ErrCatalogStale,
		)
	}

	if options.VerifySourceContent {
		if err := sourceimpl.VerifySnapshotContentDigest(
			ctx,
			s.sources,
			sourceValue,
			record.Binding.Locator,
			sourceGeneration,
			*occurrence.SourceContentDigest,
			basespec.MaxCandidateBytes,
		); err != nil {
			return artifactConsumerAPIresource.ResolvedArtifact{}, err
		}
	}

	output := artifactConsumerAPIresource.ResolvedArtifact{
		Artifact:         record.Clone(),
		Collection:       collectionValue.Clone(),
		Definition:       definitionValue.Clone(),
		Occurrence:       occurrence.Clone(),
		Source:           sourceValue.Summary(),
		CatalogRevision:  snapshot.Revision,
		SourceGeneration: sourceGeneration,
	}
	if err := output.Validate(); err != nil {
		return artifactConsumerAPIresource.ResolvedArtifact{}, err
	}
	return output.Clone(), nil
}

// ResolveVerifiedLocalPath verifies the current source occurrence and then
// resolves a native path through the selected Source adapter.
func (s *Service) ResolveVerifiedLocalPath(
	ctx context.Context,
	resolved artifactConsumerAPIresource.ResolvedArtifact,
	localLocator basespec.Locator,
) (string, error) {
	if err := validateContext(ctx, "verified local-path resolution"); err != nil {
		return "", err
	}
	if s == nil || s.sources == nil {
		return "", basespec.ErrClosed
	}
	if err := resolved.Validate(); err != nil {
		return "", err
	}
	if err := localLocator.Validate(true); err != nil {
		return "", err
	}

	sourceValue, err := s.sources.Get(
		ctx,
		resolved.Source.RootID,
		resolved.Source.ID,
	)
	if err != nil {
		return "", err
	}
	if sourceValue.Revision != resolved.Source.Revision ||
		sourceValue.Kind != resolved.Source.Kind ||
		sourceValue.StorageKey != resolved.Source.StorageKey {
		return "", fmt.Errorf(
			"%w: Source changed after Artifact resolution",
			basespec.ErrCatalogStale,
		)
	}

	return sourceimpl.ResolveVerifiedLocalPath(
		ctx,
		s.sources,
		sourceValue,
		resolved.Occurrence.Key.Locator,
		localLocator,
		resolved.SourceGeneration,
		*resolved.Occurrence.SourceContentDigest,
		basespec.MaxCandidateBytes,
	)
}

// ReadCollectionEntry reads one entry from an exact current Collection Source.
// It is used for source-owned aggregate documents such as an MCP Bundle.
func (s *Service) ReadCollectionEntry(
	ctx context.Context,
	ref collection.CollectionRef,
	sourceID source.SourceID,
	locator basespec.Locator,
	maximumBytes int64,
) (artifactConsumerAPIresource.VerifiedEntry, error) {
	if err := validateContext(ctx, "Collection source read"); err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}
	if s == nil || s.catalogs == nil || s.sources == nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, basespec.ErrClosed
	}
	if err := ref.Validate(); err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}
	if err := sourceID.Validate(); err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}
	if err := locator.Validate(false); err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}

	snapshot, err := s.catalogs.CurrentCatalog(ctx, ref)
	if err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}
	sourceRevision := snapshot.SourceRevisions[sourceID]
	sourceGeneration := snapshot.SourceGenerations[sourceID]
	if sourceRevision == 0 || sourceGeneration == "" {
		return artifactConsumerAPIresource.VerifiedEntry{}, fmt.Errorf(
			"%w: Source %q has no current Collection Catalog state",
			basespec.ErrCatalogStale,
			sourceID,
		)
	}

	sourceValue, err := s.sources.Get(ctx, ref.RootID, sourceID)
	if err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}
	if sourceValue.Revision != sourceRevision {
		return artifactConsumerAPIresource.VerifiedEntry{}, fmt.Errorf(
			"%w: Collection Source changed after Catalog publication",
			basespec.ErrCatalogStale,
		)
	}

	content, digest, err := sourceimpl.ReadVerifiedSnapshotEntry(
		ctx,
		s.sources,
		sourceValue,
		locator,
		sourceGeneration,
		maximumBytes,
	)
	if err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}

	output := artifactConsumerAPIresource.VerifiedEntry{
		Collection:       ref,
		SourceID:         sourceID,
		CatalogRevision:  snapshot.Revision,
		SourceRevision:   sourceRevision,
		SourceGeneration: sourceGeneration,
		Content:          append([]byte(nil), content...),
		Digest:           digest,
	}
	if err := output.Validate(); err != nil {
		return artifactConsumerAPIresource.VerifiedEntry{}, err
	}
	return output.Clone(), nil
}

// ResolveSourceLocalPath is a presentation-only path projection. It keeps
// Source configuration inside Artifact Store while allowing a trusted local
// UI to display an attached filesystem Source path.
func (s *Service) ResolveSourceLocalPath(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
	locator basespec.Locator,
) (string, error) {
	if err := validateContext(ctx, "Source local-path resolution"); err != nil {
		return "", err
	}
	if s == nil || s.sources == nil {
		return "", basespec.ErrClosed
	}
	if err := rootID.Validate(); err != nil {
		return "", err
	}
	if err := sourceID.Validate(); err != nil {
		return "", err
	}
	if err := locator.Validate(true); err != nil {
		return "", err
	}

	localPaths, supported := s.sources.(sourceimpl.LocalPathRuntime)
	if !supported {
		return "", fmt.Errorf(
			"%w: Artifact Store Source runtime has no local-path capability",
			basespec.ErrUnsupported,
		)
	}

	sourceValue, err := s.sources.Get(ctx, rootID, sourceID)
	if err != nil {
		return "", err
	}
	if !localPaths.SupportsLocalPath(sourceValue.Kind) {
		return "", fmt.Errorf(
			"%w: Source kind %q has no local path",
			basespec.ErrUnsupported,
			sourceValue.Kind,
		)
	}
	return localPaths.ResolveLocalPath(ctx, sourceValue, locator)
}

func (s *Service) SupportsLocalPath(kind source.SourceKind) bool {
	if s == nil || s.sources == nil {
		return false
	}
	localPaths, supported := s.sources.(sourceimpl.LocalPathRuntime)
	return supported && localPaths.SupportsLocalPath(kind)
}

func validateContext(ctx context.Context, operation string) error {
	if ctx == nil {
		return fmt.Errorf(
			"%w: %s context is nil",
			basespec.ErrInvalid,
			operation,
		)
	}
	return ctx.Err()
}
