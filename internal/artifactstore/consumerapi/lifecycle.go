package consumerapi

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

// RootLifecycle is the consumer-facing Root lifecycle contract.
type RootLifecycle interface {
	CreateRoot(
		ctx context.Context,
		draft root.RootDraft,
	) (root.Root, error)

	GetRoot(
		ctx context.Context,
		rootID root.RootID,
	) (root.Root, error)

	ListRoots(
		ctx context.Context,
	) ([]root.Root, error)

	UpdateRoot(
		ctx context.Context,
		rootID root.RootID,
		update root.RootUpdate,
	) (root.Root, error)

	RetireRoot(
		ctx context.Context,
		rootID root.RootID,
		expectedRevision uint64,
	) (root.Root, error)

	PurgeRoot(
		ctx context.Context,
		rootID root.RootID,
		expectedRevision uint64,
	) error
}

// SourceLifecycle is the consumer-facing Source lifecycle contract.
type SourceLifecycle interface {
	CreateSource(
		ctx context.Context,
		rootID root.RootID,
		draft source.Draft,
	) (source.Summary, error)

	CreateSourceWithStatus(
		ctx context.Context,
		rootID root.RootID,
		draft source.Draft,
	) (source.Summary, bool, error)

	DiscardSource(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
		expectedRevision uint64,
	) error

	GetSource(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
	) (source.Summary, error)

	ListSources(
		ctx context.Context,
		rootID root.RootID,
	) ([]source.Summary, error)

	UpdateSource(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
		update source.Update,
	) (source.Summary, error)

	RetireSource(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
		expectedRevision uint64,
	) (source.Summary, error)

	PurgeSource(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
		expectedRevision uint64,
	) error

	ListSourceKinds(
		ctx context.Context,
	) ([]source.SourceKind, error)
}
