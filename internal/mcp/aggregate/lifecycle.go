package aggregate

import (
	"context"
	"errors"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	mcpServer "github.com/flexigpt/flexigpt-app/internal/mcp/runtime/server"
	mcpConsumerAPI "github.com/flexigpt/flexigpt-app/internal/mcp/store/consumerapi"
	mcpDomainServer "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain/server"
)

// RuntimeInvalidator is the only Runtime capability required by Store mutation
// coordination. Runtime does not know who persists a server.
type RuntimeInvalidator interface {
	Invalidate(
		ctx context.Context,
		server mcpServer.ServerID,
	) error

	InvalidateCollection(
		ctx context.Context,
		catalog mcpServer.CatalogID,
	) error
}

type Lifecycle struct {
	store   mcpConsumerAPI.BundleMutator
	runtime RuntimeInvalidator
}

func NewLifecycle(
	store mcpConsumerAPI.BundleMutator,
	runtime RuntimeInvalidator,
) (*Lifecycle, error) {
	if store == nil || runtime == nil {
		return nil, errors.New("MCP lifecycle dependencies are incomplete")
	}
	return &Lifecycle{
		store:   store,
		runtime: runtime,
	}, nil
}

func (l *Lifecycle) InvalidateServer(
	ctx context.Context,
	ref artifact.ArtifactRef,
) error {
	if l == nil || l.runtime == nil {
		return mcpServer.ErrClosed
	}
	serverID, err := RuntimeServerIDForArtifact(ref)
	if err != nil {
		return err
	}
	return l.runtime.Invalidate(ctx, serverID)
}

func (l *Lifecycle) InvalidateCollection(
	ctx context.Context,
	ref collection.CollectionRef,
) error {
	if l == nil || l.runtime == nil {
		return mcpServer.ErrClosed
	}
	catalogID, err := RuntimeCatalogIDForCollection(ref)
	if err != nil {
		return err
	}
	return l.runtime.InvalidateCollection(ctx, catalogID)
}

func (l *Lifecycle) ReplaceDocument(
	ctx context.Context,
	request mcpConsumerAPI.ReplaceDocumentRequest,
) (mcpConsumerAPI.Bundle, error) {
	commit, err := l.store.PrepareReplaceDocument(ctx, request)
	if err != nil {
		return mcpConsumerAPI.Bundle{}, err
	}
	if err := l.InvalidateCollection(ctx, request.Bundle); err != nil {
		return mcpConsumerAPI.Bundle{}, err
	}
	return commit(ctx)
}

func (l *Lifecycle) RefreshBundle(
	ctx context.Context,
	ref collection.CollectionRef,
	allowProtected bool,
) (mcpConsumerAPI.Bundle, error) {
	commit, err := l.store.PrepareRefresh(ctx, ref, allowProtected)
	if err != nil {
		return mcpConsumerAPI.Bundle{}, err
	}
	if err := l.InvalidateCollection(ctx, ref); err != nil {
		return mcpConsumerAPI.Bundle{}, err
	}
	return commit(ctx)
}

func (l *Lifecycle) UpdateBundleEnabled(
	ctx context.Context,
	ref collection.CollectionRef,
	expectedRevision uint64,
	enabled bool,
) (mcpConsumerAPI.Bundle, error) {
	commit, err := l.store.PrepareUpdateBundleEnabled(
		ctx,
		ref,
		expectedRevision,
		enabled,
	)
	if err != nil {
		return mcpConsumerAPI.Bundle{}, err
	}
	if err := l.InvalidateCollection(ctx, ref); err != nil {
		return mcpConsumerAPI.Bundle{}, err
	}
	return commit(ctx)
}

func (l *Lifecycle) RetireBundle(
	ctx context.Context,
	ref collection.CollectionRef,
	expectedRevision uint64,
) (collection.Collection, error) {
	commit, err := l.store.PrepareRetire(ctx, ref, expectedRevision)
	if err != nil {
		return collection.Collection{}, err
	}
	if err := l.InvalidateCollection(ctx, ref); err != nil {
		return collection.Collection{}, err
	}
	return commit(ctx)
}

func (l *Lifecycle) PurgeBundle(
	ctx context.Context,
	ref collection.CollectionRef,
	expectedRevision uint64,
) error {
	commit, err := l.store.PreparePurge(ctx, ref, expectedRevision)
	if err != nil {
		return err
	}
	if err := l.InvalidateCollection(ctx, ref); err != nil {
		return err
	}
	return commit(ctx)
}

func (l *Lifecycle) UpdateProtectedBundleInstallation(
	ctx context.Context,
	ref collection.CollectionRef,
	expectedOverlayRevision uint64,
	runtimeEnabled bool,
) error {
	commit, err := l.store.PrepareUpdateProtectedBundleInstallation(
		ctx,
		ref,
		expectedOverlayRevision,
		runtimeEnabled,
	)
	if err != nil {
		return err
	}
	if err := l.InvalidateCollection(ctx, ref); err != nil {
		return err
	}
	return commit(ctx)
}

func (l *Lifecycle) UpdateServerInstallation(
	ctx context.Context,
	ref artifact.ArtifactRef,
	expectedArtifactRevision uint64,
	data mcpDomainServer.ServerData,
) (artifact.Artifact, error) {
	commit, err := l.store.PrepareUpdateServerInstallation(
		ctx,
		ref,
		expectedArtifactRevision,
		data,
	)
	if err != nil {
		return artifact.Artifact{}, err
	}
	if err := l.InvalidateServer(ctx, ref); err != nil {
		return artifact.Artifact{}, err
	}
	return commit(ctx)
}

func (l *Lifecycle) UpdateProtectedServerInstallation(
	ctx context.Context,
	ref artifact.ArtifactRef,
	expectedOverlayRevision uint64,
	runtimeEnabled bool,
	data mcpDomainServer.ServerData,
) error {
	commit, err := l.store.PrepareUpdateProtectedServerInstallation(
		ctx,
		ref,
		expectedOverlayRevision,
		runtimeEnabled,
		data,
	)
	if err != nil {
		return err
	}
	if err := l.InvalidateServer(ctx, ref); err != nil {
		return err
	}
	return commit(ctx)
}
