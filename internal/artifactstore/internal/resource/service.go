package resource

import (
	"context"
	"errors"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/definition"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/resource"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	artifactimpl "github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/artifact"
	collectionimpl "github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/collection"
	sourceimpl "github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/source"
)

type catalogReader interface {
	CurrentCatalog(
		ctx context.Context,
		ref collection.CollectionRef,
	) (catalog.Snapshot, error)

	InspectCollectionCatalog(
		ctx context.Context,
		ref collection.CollectionRef,
	) (catalog.CatalogInspection, error)
}

type catalogArtifactKey struct {
	Occurrence catalog.OccurrenceKey
	Kind       artifact.ArtifactKind
}

type catalogArtifactMaterial struct {
	occurrence       catalog.Occurrence
	definition       definition.Definition
	source           source.Source
	sourceGeneration string
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

// InspectCollectionResources returns the generic resource linkage for one
// Collection without applying Workspace, Skill, MCP, or runtime policy.
func (s *Service) InspectCollectionResources(
	ctx context.Context,
	ref collection.CollectionRef,
) (resource.CollectionResourceInspection, error) {
	if err := validateContext(ctx, "Collection resource inspection"); err != nil {
		return resource.CollectionResourceInspection{}, err
	}
	if s == nil ||
		s.artifacts == nil ||
		s.collections == nil ||
		s.catalogs == nil ||
		s.sources == nil {
		return resource.CollectionResourceInspection{}, basespec.ErrClosed
	}
	if err := ref.Validate(); err != nil {
		return resource.CollectionResourceInspection{}, err
	}

	collectionValue, err := s.collections.Get(ctx, ref)
	if err != nil {
		return resource.CollectionResourceInspection{}, err
	}
	if collectionValue.Ref() != ref {
		return resource.CollectionResourceInspection{}, fmt.Errorf(
			"%w: Collection reader returned another Collection",
			basespec.ErrInvalid,
		)
	}

	inspection, err := s.catalogs.InspectCollectionCatalog(ctx, ref)
	if err != nil {
		return resource.CollectionResourceInspection{}, err
	}
	snapshot := inspection.Catalog
	if err := snapshot.Validate(); err != nil {
		return resource.CollectionResourceInspection{}, fmt.Errorf(
			"%w: catalog inspection returned an invalid Catalog: %w",
			basespec.ErrInvalid,
			err,
		)
	}

	records, err := s.artifacts.ListByCollection(ctx, ref)
	if err != nil {
		return resource.CollectionResourceInspection{}, err
	}

	occurrences := catalogOccurrenceIndex(snapshot.Occurrences)
	recorded := make(map[catalogArtifactKey]struct{}, len(records))
	output := resource.CollectionResourceInspection{
		Catalog:             inspection.Clone(),
		Resources:           make([]resource.CollectionArtifactResource, 0, len(records)),
		UnresolvedArtifacts: make([]resource.ArtifactResolutionIssue, 0),
	}

	for _, record := range records {
		if record.RootID != ref.RootID ||
			record.CollectionID != ref.CollectionID {
			return resource.CollectionResourceInspection{}, fmt.Errorf(
				"%w: Artifact belongs to another Collection",
				basespec.ErrInvalid,
			)
		}

		recorded[catalogArtifactKeyFor(record)] = struct{}{}
		if record.ResolvedDefinition == nil {
			output.UnresolvedArtifacts = append(
				output.UnresolvedArtifacts,
				resource.ArtifactResolutionIssue{
					Artifact: record.Clone(),
					Status:   resource.ArtifactResolutionRecordUnavailable,
				},
			)
			continue
		}

		material, err := s.catalogArtifactMaterialFor(
			ctx,
			record,
			collectionValue,
			snapshot,
			occurrences,
			false,
		)
		if err != nil {
			status, expected := artifactResolutionStatusFor(err)
			if !expected {
				return resource.CollectionResourceInspection{}, err
			}
			output.UnresolvedArtifacts = append(
				output.UnresolvedArtifacts,
				resource.ArtifactResolutionIssue{
					Artifact: record.Clone(),
					Status:   status,
				},
			)
			continue
		}

		current := inspection.IsCurrent() &&
			material.source.Revision ==
				snapshot.SourceRevisions[record.Binding.SourceID]

		linked := resource.CollectionArtifactResource{
			Artifact:       record.Clone(),
			Definition:     material.definition.Clone(),
			Occurrence:     material.occurrence.Clone(),
			Source:         material.source.Summary(),
			CatalogCurrent: current,
		}
		if record.State == artifact.StateAvailable && current {
			resolved, err := resolvedArtifactFromMaterial(
				record,
				collectionValue,
				snapshot,
				material,
			)
			if err != nil {
				return resource.CollectionResourceInspection{}, err
			}
			linked.Resolved = &resolved
		}
		output.Resources = append(output.Resources, linked)
	}

	for _, occurrence := range snapshot.Occurrences {
		if occurrence.Kind == "" {
			continue
		}
		key := catalogArtifactKey{
			Occurrence: occurrence.Key,
			Kind:       occurrence.Kind,
		}
		if _, found := recorded[key]; found {
			continue
		}
		output.UnrecordedOccurrences = append(
			output.UnrecordedOccurrences,
			occurrence.Clone(),
		)
	}

	return output.Clone(), nil
}

// ResolveVerifiedLocalPath verifies the current source occurrence and then
// resolves a native path through the selected Source adapter.
func (s *Service) ResolveVerifiedLocalPath(
	ctx context.Context,
	resolved resource.ResolvedArtifact,
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
) (resource.VerifiedEntry, error) {
	value, err := s.ReadCollectionEntryWithCatalog(
		ctx,
		ref,
		sourceID,
		locator,
		maximumBytes,
	)
	if err != nil {
		return resource.VerifiedEntry{}, err
	}
	return value.Entry.Clone(), nil
}

// ReadCollectionEntryWithCatalog reads one exact current Collection entry and
// returns the exact Catalog snapshot used to validate that entry.
func (s *Service) ReadCollectionEntryWithCatalog(
	ctx context.Context,
	ref collection.CollectionRef,
	sourceID source.SourceID,
	locator basespec.Locator,
	maximumBytes int64,
) (resource.VerifiedCollectionEntry, error) {
	if err := validateContext(ctx, "Collection source read"); err != nil {
		return resource.VerifiedCollectionEntry{}, err
	}
	if s == nil || s.catalogs == nil || s.sources == nil {
		return resource.VerifiedCollectionEntry{}, basespec.ErrClosed
	}
	if err := ref.Validate(); err != nil {
		return resource.VerifiedCollectionEntry{}, err
	}
	if err := sourceID.Validate(); err != nil {
		return resource.VerifiedCollectionEntry{}, err
	}
	if err := locator.Validate(false); err != nil {
		return resource.VerifiedCollectionEntry{}, err
	}

	snapshot, err := s.catalogs.CurrentCatalog(ctx, ref)
	if err != nil {
		return resource.VerifiedCollectionEntry{}, err
	}
	sourceRevision := snapshot.SourceRevisions[sourceID]
	sourceGeneration := snapshot.SourceGenerations[sourceID]
	if sourceRevision == 0 || sourceGeneration == "" {
		return resource.VerifiedCollectionEntry{}, fmt.Errorf(
			"%w: Source %q has no current Collection Catalog state",
			basespec.ErrCatalogStale,
			sourceID,
		)
	}

	sourceValue, err := s.sources.Get(ctx, ref.RootID, sourceID)
	if err != nil {
		return resource.VerifiedCollectionEntry{}, err
	}
	if sourceValue.Revision != sourceRevision {
		return resource.VerifiedCollectionEntry{}, fmt.Errorf(
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
		return resource.VerifiedCollectionEntry{}, err
	}

	entry := resource.VerifiedEntry{
		Collection:       ref,
		SourceID:         sourceID,
		CatalogRevision:  snapshot.Revision,
		SourceRevision:   sourceRevision,
		SourceGeneration: sourceGeneration,
		Content:          append([]byte(nil), content...),
		Digest:           digest,
	}
	output := resource.VerifiedCollectionEntry{
		Catalog: snapshot.Clone(),
		Entry:   entry,
	}
	if err := output.Validate(); err != nil {
		return resource.VerifiedCollectionEntry{}, err
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

// ResolveArtifact verifies the complete current resource chain:
//
// Artifact -> Collection -> current provider plan -> Catalog occurrence ->
// Definition -> Source revision -> optional exact source bytes.
func (s *Service) ResolveArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
	options resource.ResolveOptions,
) (resource.ResolvedArtifact, error) {
	if err := validateContext(ctx, "Artifact resolution"); err != nil {
		return resource.ResolvedArtifact{}, err
	}
	if s == nil ||
		s.artifacts == nil ||
		s.collections == nil ||
		s.catalogs == nil ||
		s.sources == nil {
		return resource.ResolvedArtifact{}, basespec.ErrClosed
	}
	if err := ref.Validate(); err != nil {
		return resource.ResolvedArtifact{}, err
	}

	record, err := s.artifacts.Get(ctx, ref)
	if err != nil {
		return resource.ResolvedArtifact{}, err
	}
	if record.State != artifact.StateAvailable ||
		record.ResolvedDefinition == nil {
		return resource.ResolvedArtifact{}, fmt.Errorf(
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
		return resource.ResolvedArtifact{}, err
	}
	if collectionValue.Ref() != collectionRef {
		return resource.ResolvedArtifact{}, fmt.Errorf(
			"%w: Collection reader returned another Collection",
			basespec.ErrInvalid,
		)
	}
	snapshot, err := s.catalogs.CurrentCatalog(ctx, collectionRef)
	if err != nil {
		return resource.ResolvedArtifact{}, err
	}
	material, err := s.catalogArtifactMaterialFor(
		ctx,
		record,
		collectionValue,
		snapshot,
		catalogOccurrenceIndex(snapshot.Occurrences),
		true,
	)
	if err != nil {
		return resource.ResolvedArtifact{}, err
	}
	output, err := resolvedArtifactFromMaterial(
		record,
		collectionValue,
		snapshot,
		material,
	)
	if err != nil {
		return resource.ResolvedArtifact{}, err
	}
	if options.VerifySourceContent {
		if err := sourceimpl.VerifySnapshotContentDigest(
			ctx,
			s.sources,
			material.source,
			record.Binding.Locator,
			material.sourceGeneration,
			*material.occurrence.SourceContentDigest,
			basespec.MaxCandidateBytes,
		); err != nil {
			return resource.ResolvedArtifact{}, err
		}
	}
	return output.Clone(), nil
}

func (s *Service) catalogArtifactMaterialFor(
	ctx context.Context,
	record artifact.Artifact,
	collectionValue collection.Collection,
	snapshot catalog.Snapshot,
	occurrences map[catalog.OccurrenceKey]catalog.Occurrence,
	requireCurrentSource bool,
) (catalogArtifactMaterial, error) {
	if record.RootID != collectionValue.RootID ||
		record.CollectionID != collectionValue.ID {
		return catalogArtifactMaterial{}, fmt.Errorf(
			"%w: Artifact belongs to another Collection",
			basespec.ErrInvalid,
		)
	}
	if record.ResolvedDefinition == nil {
		return catalogArtifactMaterial{}, fmt.Errorf(
			"%w: Artifact %q has no resolved definition",
			basespec.ErrReferenceUnresolved,
			record.ID,
		)
	}

	key := catalog.OccurrenceKey{
		CollectionID:       record.CollectionID,
		SourceID:           record.Binding.SourceID,
		Locator:            record.Binding.Locator,
		SubresourceLocator: record.Binding.SubresourceLocator,
	}
	occurrence, found := occurrences[key]
	if !found ||
		occurrence.State != catalog.OccurrenceValid ||
		occurrence.Kind != record.Kind ||
		occurrence.DefinitionDigest == nil ||
		occurrence.SourceContentDigest == nil ||
		*occurrence.DefinitionDigest != *record.ResolvedDefinition {
		return catalogArtifactMaterial{}, fmt.Errorf(
			"%w: Artifact %q does not match its current Catalog occurrence",
			basespec.ErrCatalogStale,
			record.ID,
		)
	}

	definitionValue, err := snapshot.DefinitionForOccurrence(key)
	if err != nil {
		return catalogArtifactMaterial{}, err
	}
	if definitionValue.Kind != record.Kind ||
		definitionValue.Digest != *record.ResolvedDefinition {
		return catalogArtifactMaterial{}, fmt.Errorf(
			"%w: Artifact %q definition does not match current state",
			basespec.ErrDigestMismatch,
			record.ID,
		)
	}

	sourceRevision := snapshot.SourceRevisions[record.Binding.SourceID]
	sourceGeneration := snapshot.SourceGenerations[record.Binding.SourceID]
	if sourceRevision == 0 || sourceGeneration == "" {
		return catalogArtifactMaterial{}, fmt.Errorf(
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
		return catalogArtifactMaterial{}, err
	}
	if requireCurrentSource && sourceValue.Revision != sourceRevision {
		return catalogArtifactMaterial{}, fmt.Errorf(
			"%w: Artifact Source changed after Catalog publication",
			basespec.ErrCatalogStale,
		)
	}

	return catalogArtifactMaterial{
		occurrence:       occurrence.Clone(),
		definition:       definitionValue.Clone(),
		source:           sourceValue.Clone(),
		sourceGeneration: sourceGeneration,
	}, nil
}

func resolvedArtifactFromMaterial(
	record artifact.Artifact,
	collectionValue collection.Collection,
	snapshot catalog.Snapshot,
	material catalogArtifactMaterial,
) (resource.ResolvedArtifact, error) {
	if record.State != artifact.StateAvailable {
		return resource.ResolvedArtifact{}, fmt.Errorf(
			"%w: Artifact %q is not currently available",
			basespec.ErrReferenceUnresolved,
			record.ID,
		)
	}

	output := resource.ResolvedArtifact{
		Artifact:         record.Clone(),
		Collection:       collectionValue.Clone(),
		Definition:       material.definition.Clone(),
		Occurrence:       material.occurrence.Clone(),
		Source:           material.source.Summary(),
		CatalogRevision:  snapshot.Revision,
		SourceGeneration: material.sourceGeneration,
	}
	if err := output.Validate(); err != nil {
		return resource.ResolvedArtifact{}, fmt.Errorf(
			"%w: resolved Artifact material is invalid: %w",
			basespec.ErrInvalid,
			err,
		)
	}
	return output.Clone(), nil
}

func catalogOccurrenceIndex(
	values []catalog.Occurrence,
) map[catalog.OccurrenceKey]catalog.Occurrence {
	output := make(map[catalog.OccurrenceKey]catalog.Occurrence, len(values))
	for _, value := range values {
		output[value.Key] = value.Clone()
	}
	return output
}

func catalogArtifactKeyFor(
	record artifact.Artifact,
) catalogArtifactKey {
	return catalogArtifactKey{
		Occurrence: catalog.OccurrenceKey{
			CollectionID:       record.CollectionID,
			SourceID:           record.Binding.SourceID,
			Locator:            record.Binding.Locator,
			SubresourceLocator: record.Binding.SubresourceLocator,
		},
		Kind: record.Kind,
	}
}

func artifactResolutionStatusFor(
	err error,
) (resource.ArtifactResolutionStatus, bool) {
	switch {
	case errors.Is(err, basespec.ErrSourceNotFound),
		errors.Is(err, basespec.ErrSourceUnavailable):
		return resource.ArtifactResolutionSourceUnavailable, true

	case errors.Is(err, basespec.ErrDefinitionNotFound):
		return resource.ArtifactResolutionDefinitionUnavailable, true

	case errors.Is(err, basespec.ErrDigestMismatch):
		return resource.ArtifactResolutionDefinitionMismatch, true

	case errors.Is(err, basespec.ErrCatalogStale):
		return resource.ArtifactResolutionOccurrenceUnavailable, true

	case errors.Is(err, basespec.ErrReferenceUnresolved):
		return resource.ArtifactResolutionRecordUnavailable, true

	default:
		return "", false
	}
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
