package compositionapi

import (
	"context"
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/resource"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/schema"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

type RootAPI interface {
	Create(
		ctx context.Context,
		draft root.RootDraft,
	) (root.Root, error)

	Get(
		ctx context.Context,
		id root.RootID,
	) (root.Root, error)

	List(ctx context.Context) ([]root.Root, error)

	Update(
		ctx context.Context,
		id root.RootID,
		update root.RootUpdate,
	) (root.Root, error)

	Retire(
		ctx context.Context,
		id root.RootID,
		expectedRevision uint64,
	) (root.Root, error)

	Purge(
		ctx context.Context,
		id root.RootID,
		expectedRevision uint64,
	) error
}

type SourceAPI interface {
	Create(
		ctx context.Context,
		rootID root.RootID,
		draft source.Draft,
	) (source.Summary, error)

	CreateWithStatus(
		ctx context.Context,
		rootID root.RootID,
		draft source.Draft,
	) (source.Summary, bool, error)

	Discard(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
		expectedRevision uint64,
	) error

	Get(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
	) (source.Summary, error)

	List(
		ctx context.Context,
		rootID root.RootID,
	) ([]source.Summary, error)

	Update(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
		update source.Update,
	) (source.Summary, error)

	Retire(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
		expectedRevision uint64,
	) (source.Summary, error)

	Purge(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
		expectedRevision uint64,
	) error

	Kinds() []source.SourceKind
}

type CollectionAPI interface {
	Create(
		ctx context.Context,
		rootID root.RootID,
		draft collection.Draft,
		attachments []collection.AttachmentDraft,
	) (collection.Collection, []collection.Attachment, error)

	Get(
		ctx context.Context,
		ref collection.CollectionRef,
	) (collection.Collection, error)

	GetRetired(
		ctx context.Context,
		ref collection.CollectionRef,
	) (collection.Collection, error)

	ListByRoot(
		ctx context.Context,
		rootID root.RootID,
	) ([]collection.Collection, error)

	Update(
		ctx context.Context,
		ref collection.CollectionRef,
		update collection.Update,
	) (collection.Collection, error)

	Retire(
		ctx context.Context,
		ref collection.CollectionRef,
		expectedRevision uint64,
	) (collection.Collection, error)

	Purge(
		ctx context.Context,
		ref collection.CollectionRef,
		expectedRevision uint64,
	) error

	Attach(
		ctx context.Context,
		ref collection.CollectionRef,
		expectedCollectionRevision uint64,
		draft collection.AttachmentDraft,
	) (collection.Collection, collection.Attachment, error)

	GetAttachment(
		ctx context.Context,
		ref collection.CollectionRef,
		sourceID source.SourceID,
	) (collection.Attachment, error)

	ListAttachments(
		ctx context.Context,
		ref collection.CollectionRef,
	) ([]collection.Attachment, error)

	UpdateAttachment(
		ctx context.Context,
		ref collection.CollectionRef,
		sourceID source.SourceID,
		update collection.AttachmentUpdate,
	) (collection.Collection, collection.Attachment, error)

	Detach(
		ctx context.Context,
		ref collection.CollectionRef,
		sourceID source.SourceID,
		expectedCollectionRevision uint64,
		expectedAttachmentRevision uint64,
	) (collection.Collection, error)

	ReplaceAttachment(
		ctx context.Context,
		ref collection.CollectionRef,
		replacement collection.AttachmentReplacement,
	) (collection.Collection, collection.Attachment, error)
}

type ArtifactAPI interface {
	Get(
		ctx context.Context,
		ref artifact.ArtifactRef,
	) (artifact.Artifact, error)

	ListByCollection(
		ctx context.Context,
		ref collection.CollectionRef,
	) ([]artifact.Artifact, error)

	Adopt(
		ctx context.Context,
		request catalog.AdoptRequest,
	) (artifact.Artifact, error)

	Pin(
		ctx context.Context,
		request catalog.PinRequest,
	) (artifact.Artifact, error)

	SetEnabled(
		ctx context.Context,
		ref artifact.ArtifactRef,
		expectedRevision uint64,
		enabled bool,
	) (artifact.Artifact, error)

	SetName(
		ctx context.Context,
		ref artifact.ArtifactRef,
		expectedRevision uint64,
		name string,
	) (artifact.Artifact, error)

	UpdateData(
		ctx context.Context,
		ref artifact.ArtifactRef,
		expectedRevision uint64,
		data json.RawMessage,
	) (artifact.Artifact, error)

	Unadopt(
		ctx context.Context,
		ref artifact.ArtifactRef,
		expectedRevision uint64,
		suppress bool,
	) error

	Purge(
		ctx context.Context,
		ref artifact.ArtifactRef,
		expectedRevision uint64,
	) error

	PurgeAndSuppress(
		ctx context.Context,
		ref artifact.ArtifactRef,
		expectedRevision uint64,
	) error

	ListSuppressions(
		ctx context.Context,
		ref collection.CollectionRef,
	) ([]artifact.Suppression, error)

	Suppress(
		ctx context.Context,
		request catalog.SuppressRequest,
	) (artifact.Suppression, error)

	Unsuppress(
		ctx context.Context,
		ref collection.CollectionRef,
		binding artifact.SourceBinding,
		expectedRevision uint64,
	) error
}

type CatalogAPI interface {
	RefreshCollection(
		ctx context.Context,
		ref collection.CollectionRef,
	) (catalog.RefreshCollectionResult, error)

	CurrentCatalog(
		ctx context.Context,
		ref collection.CollectionRef,
	) (catalog.Snapshot, error)

	InspectCollectionCatalog(
		ctx context.Context,
		ref collection.CollectionRef,
	) (catalog.CatalogInspection, error)
}

type ResourceAPI interface {
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
		sourceID source.SourceID,
		locator basespec.Locator,
		maximumBytes int64,
	) (resource.VerifiedEntry, error)

	ResolveSourceLocalPath(
		ctx context.Context,
		rootID root.RootID,
		sourceID source.SourceID,
		locator basespec.Locator,
	) (string, error)

	SupportsLocalPath(kind source.SourceKind) bool
}

type SchemaAPI interface {
	CanonicalizeExpected(
		ctx context.Context,
		expected schema.Key,
		raw []byte,
	) (schema.ParsedDocument, error)
}

type ManagedArtifactAPI interface {
	Publish(
		ctx context.Context,
		request artifact.PublishArtifactRequest,
	) (artifact.PublishArtifactResult, error)

	PublishCollection(
		ctx context.Context,
		request collection.PublishCollectionRequest,
	) (collection.PublishCollectionResult, error)

	Remove(
		ctx context.Context,
		request artifact.RemoveArtifactRequest,
	) error
}

type ProtectionAPI interface {
	IsProtectedRoot(rootID root.RootID) bool
	RequirePrivilegedInstaller(ctx context.Context) error
}
