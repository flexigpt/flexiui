package consumerapi

import (
	"context"
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/schema"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/installerapi"
)

func (a *API) IsProtectedRoot(rootID root.RootID) bool {
	return a != nil &&
		a.components != nil &&
		a.components.RootMutationPolicy() != nil &&
		a.components.RootMutationPolicy().IsProtectedRoot(rootID)
}

func (a *API) RequirePrivilegedInstaller(ctx context.Context) error {
	if err := a.check(ctx); err != nil {
		return err
	}
	return installerapi.RequirePrivileged(ctx)
}

func (a *API) CreateCollection(
	ctx context.Context,
	rootID root.RootID,
	draft collection.Draft,
	attachments []collection.AttachmentDraft,
) (collection.Collection, []collection.Attachment, error) {
	return a.components.Collections.Create(
		ctx,
		rootID,
		draft,
		cloneAttachmentDrafts(attachments),
	)
}

func (a *API) GetCollection(
	ctx context.Context,
	ref collection.CollectionRef,
) (collection.Collection, error) {
	return a.components.Collections.Get(ctx, ref)
}

func (a *API) GetRetiredCollection(
	ctx context.Context,
	ref collection.CollectionRef,
) (collection.Collection, error) {
	return a.components.Collections.GetRetired(ctx, ref)
}

func (a *API) ListCollections(
	ctx context.Context,
	rootID root.RootID,
) ([]collection.Collection, error) {
	return a.components.Collections.ListByRoot(ctx, rootID)
}

func (a *API) UpdateCollection(
	ctx context.Context,
	ref collection.CollectionRef,
	update collection.Update,
) (collection.Collection, error) {
	update.Data = append(json.RawMessage(nil), update.Data...)
	return a.components.Collections.Update(ctx, ref, update)
}

func (a *API) RetireCollection(
	ctx context.Context,
	ref collection.CollectionRef,
	expectedRevision uint64,
) (collection.Collection, error) {
	return a.components.Collections.Retire(
		ctx,
		ref,
		expectedRevision,
	)
}

func (a *API) PurgeCollection(
	ctx context.Context,
	ref collection.CollectionRef,
	expectedRevision uint64,
) error {
	return a.components.Collections.Purge(
		ctx,
		ref,
		expectedRevision,
	)
}

func (a *API) AttachCollectionSource(
	ctx context.Context,
	ref collection.CollectionRef,
	expectedCollectionRevision uint64,
	draft collection.AttachmentDraft,
) (collection.Collection, collection.Attachment, error) {
	draft.Data = append(json.RawMessage(nil), draft.Data...)
	return a.components.Collections.Attach(
		ctx,
		ref,
		expectedCollectionRevision,
		draft,
	)
}

func (a *API) GetCollectionAttachment(
	ctx context.Context,
	ref collection.CollectionRef,
	sourceID source.SourceID,
) (collection.Attachment, error) {
	return a.components.Collections.GetAttachment(
		ctx,
		ref,
		sourceID,
	)
}

func (a *API) ListCollectionAttachments(
	ctx context.Context,
	ref collection.CollectionRef,
) ([]collection.Attachment, error) {
	return a.components.Collections.ListAttachments(ctx, ref)
}

func (a *API) UpdateCollectionAttachment(
	ctx context.Context,
	ref collection.CollectionRef,
	sourceID source.SourceID,
	update collection.AttachmentUpdate,
) (collection.Collection, collection.Attachment, error) {
	update.Data = append(json.RawMessage(nil), update.Data...)
	return a.components.Collections.UpdateAttachment(
		ctx,
		ref,
		sourceID,
		update,
	)
}

func (a *API) DetachCollectionSource(
	ctx context.Context,
	ref collection.CollectionRef,
	sourceID source.SourceID,
	expectedCollectionRevision uint64,
	expectedAttachmentRevision uint64,
) (collection.Collection, error) {
	return a.components.Collections.Detach(
		ctx,
		ref,
		sourceID,
		expectedCollectionRevision,
		expectedAttachmentRevision,
	)
}

func (a *API) ReplaceCollectionAttachment(
	ctx context.Context,
	ref collection.CollectionRef,
	replacement collection.AttachmentReplacement,
) (collection.Collection, collection.Attachment, error) {
	replacement.Replacement.Data = append(
		json.RawMessage(nil),
		replacement.Replacement.Data...,
	)
	return a.components.Collections.ReplaceAttachment(
		ctx,
		ref,
		replacement,
	)
}

func (a *API) GetArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
) (artifact.Artifact, error) {
	return a.components.Artifacts.Get(ctx, ref)
}

func (a *API) ListCollectionArtifacts(
	ctx context.Context,
	ref collection.CollectionRef,
) ([]artifact.Artifact, error) {
	return a.components.Artifacts.ListByCollection(ctx, ref)
}

func (a *API) AdoptArtifact(
	ctx context.Context,
	request catalog.AdoptRequest,
) (artifact.Artifact, error) {
	request.Data = append(json.RawMessage(nil), request.Data...)
	return a.components.Artifacts.Adopt(ctx, request)
}

func (a *API) PinArtifact(
	ctx context.Context,
	request catalog.PinRequest,
) (artifact.Artifact, error) {
	request.Data = append(json.RawMessage(nil), request.Data...)
	return a.components.Artifacts.Pin(ctx, request)
}

func (a *API) SetArtifactEnabled(
	ctx context.Context,
	ref artifact.ArtifactRef,
	expectedRevision uint64,
	enabled bool,
) (artifact.Artifact, error) {
	return a.components.Artifacts.SetEnabled(
		ctx,
		ref,
		expectedRevision,
		enabled,
	)
}

func (a *API) SetArtifactName(
	ctx context.Context,
	ref artifact.ArtifactRef,
	expectedRevision uint64,
	name string,
) (artifact.Artifact, error) {
	return a.components.Artifacts.SetName(
		ctx,
		ref,
		expectedRevision,
		name,
	)
}

func (a *API) UpdateArtifactData(
	ctx context.Context,
	ref artifact.ArtifactRef,
	expectedRevision uint64,
	data json.RawMessage,
) (artifact.Artifact, error) {
	return a.components.Artifacts.UpdateData(
		ctx,
		ref,
		expectedRevision,
		append(json.RawMessage(nil), data...),
	)
}

func (a *API) UnadoptArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
	expectedRevision uint64,
	suppress bool,
) error {
	return a.components.Artifacts.Unadopt(
		ctx,
		ref,
		expectedRevision,
		suppress,
	)
}

func (a *API) PurgeArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
	expectedRevision uint64,
) error {
	return a.components.Artifacts.Purge(
		ctx,
		ref,
		expectedRevision,
	)
}

func (a *API) PurgeAndSuppressArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
	expectedRevision uint64,
) error {
	return a.components.Artifacts.PurgeAndSuppress(
		ctx,
		ref,
		expectedRevision,
	)
}

func (a *API) ListCollectionSuppressions(
	ctx context.Context,
	ref collection.CollectionRef,
) ([]artifact.Suppression, error) {
	return a.components.Artifacts.ListSuppressions(ctx, ref)
}

func (a *API) SuppressBinding(
	ctx context.Context,
	request catalog.SuppressRequest,
) (artifact.Suppression, error) {
	return a.components.Artifacts.Suppress(ctx, request)
}

func (a *API) UnsuppressBinding(
	ctx context.Context,
	ref collection.CollectionRef,
	binding artifact.SourceBinding,
	expectedRevision uint64,
) error {
	return a.components.Artifacts.Unsuppress(
		ctx,
		ref,
		binding,
		expectedRevision,
	)
}

func (a *API) RefreshCollection(
	ctx context.Context,
	ref collection.CollectionRef,
) (catalog.RefreshCollectionResult, error) {
	return a.components.Refresh.RefreshCollection(ctx, ref)
}

func (a *API) CurrentCollectionCatalog(
	ctx context.Context,
	ref collection.CollectionRef,
) (catalog.Snapshot, error) {
	return a.components.Refresh.CurrentCatalog(ctx, ref)
}

func (a *API) InspectCollectionCatalog(
	ctx context.Context,
	ref collection.CollectionRef,
) (catalog.CatalogInspection, error) {
	return a.components.Refresh.InspectCollectionCatalog(ctx, ref)
}

func (a *API) CanonicalizeExpected(
	ctx context.Context,
	expected schema.Key,
	raw []byte,
) (schema.ParsedDocument, error) {
	if a.components.ShareableSchemas == nil {
		return schema.ParsedDocument{}, basespec.ErrClosed
	}
	return a.components.ShareableSchemas.CanonicalizeExpected(
		ctx,
		expected,
		append([]byte(nil), raw...),
	)
}

func (a *API) PublishManagedArtifact(
	ctx context.Context,
	request artifact.PublishArtifactRequest,
) (artifact.PublishArtifactResult, error) {
	return a.components.ManagedArtifacts.Publish(ctx, request)
}

func (a *API) PublishManagedCollection(
	ctx context.Context,
	request collection.PublishCollectionRequest,
) (collection.PublishCollectionResult, error) {
	return a.components.ManagedArtifacts.PublishCollection(ctx, request)
}

func (a *API) RemoveManagedArtifact(
	ctx context.Context,
	request artifact.RemoveArtifactRequest,
) error {
	return a.components.ManagedArtifacts.Remove(ctx, request)
}

func cloneAttachmentDrafts(
	values []collection.AttachmentDraft,
) []collection.AttachmentDraft {
	if values == nil {
		return nil
	}

	output := make(
		[]collection.AttachmentDraft,
		len(values),
	)
	for index, value := range values {
		output[index] = value
		output[index].Data = append(
			json.RawMessage(nil),
			value.Data...,
		)
	}
	return output
}
