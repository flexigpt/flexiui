package artifactstore

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/resource"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/internal/system"
)

// API implements the transport-independent Artifact Store consumer contract.
//
// Construction and shutdown are owned by compositionapi.Store.
type API struct {
	components *system.Components
	resources  *resource.Service

	closeOnce sync.Once
	closeErr  error
	closed    atomic.Bool
}

func (a *API) CreateRoot(
	ctx context.Context,
	draft root.RootDraft,
) (root.Root, error) {
	if err := a.check(ctx); err != nil {
		return root.Root{}, err
	}
	return a.components.Roots.Create(ctx, draft)
}

func (a *API) GetRoot(
	ctx context.Context,
	rootID root.RootID,
) (root.Root, error) {
	if err := a.check(ctx); err != nil {
		return root.Root{}, err
	}
	return a.components.Roots.Get(ctx, rootID)
}

func (a *API) ListRoots(
	ctx context.Context,
) ([]root.Root, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	return a.components.Roots.List(ctx)
}

func (a *API) UpdateRoot(
	ctx context.Context,
	rootID root.RootID,
	update root.RootUpdate,
) (root.Root, error) {
	if err := a.check(ctx); err != nil {
		return root.Root{}, err
	}
	return a.components.Roots.Update(ctx, rootID, update)
}

func (a *API) RetireRoot(
	ctx context.Context,
	rootID root.RootID,
	expectedRevision uint64,
) (root.Root, error) {
	if err := a.check(ctx); err != nil {
		return root.Root{}, err
	}
	return a.components.Roots.Retire(
		ctx,
		rootID,
		expectedRevision,
	)
}

func (a *API) PurgeRoot(
	ctx context.Context,
	rootID root.RootID,
	expectedRevision uint64,
) error {
	if err := a.check(ctx); err != nil {
		return err
	}
	return a.components.Roots.Purge(
		ctx,
		rootID,
		expectedRevision,
	)
}

func (a *API) CreateSource(
	ctx context.Context,
	rootID root.RootID,
	draft source.Draft,
) (source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, err
	}
	draft.Config = append(json.RawMessage(nil), draft.Config...)
	return a.components.Sources.Create(ctx, rootID, draft)
}

func (a *API) CreateSourceWithStatus(
	ctx context.Context,
	rootID root.RootID,
	draft source.Draft,
) (source.Summary, bool, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, false, err
	}
	draft.Config = append(json.RawMessage(nil), draft.Config...)
	return a.components.Sources.CreateWithStatus(ctx, rootID, draft)
}

func (a *API) DiscardSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
	expectedRevision uint64,
) error {
	if err := a.check(ctx); err != nil {
		return err
	}
	return a.components.Sources.Discard(
		ctx,
		rootID,
		sourceID,
		expectedRevision,
	)
}

func (a *API) GetSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
) (source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, err
	}
	return a.components.Sources.Get(ctx, rootID, sourceID)
}

func (a *API) ListSources(
	ctx context.Context,
	rootID root.RootID,
) ([]source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	return a.components.Sources.List(ctx, rootID)
}

func (a *API) UpdateSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
	update source.Update,
) (source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, err
	}
	update.Config = append(json.RawMessage(nil), update.Config...)
	return a.components.Sources.Update(
		ctx,
		rootID,
		sourceID,
		update,
	)
}

func (a *API) RetireSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
	expectedRevision uint64,
) (source.Summary, error) {
	if err := a.check(ctx); err != nil {
		return source.Summary{}, err
	}
	return a.components.Sources.Retire(
		ctx,
		rootID,
		sourceID,
		expectedRevision,
	)
}

func (a *API) PurgeSource(
	ctx context.Context,
	rootID root.RootID,
	sourceID source.SourceID,
	expectedRevision uint64,
) error {
	if err := a.check(ctx); err != nil {
		return err
	}
	return a.components.Sources.Purge(
		ctx,
		rootID,
		sourceID,
		expectedRevision,
	)
}

func (a *API) ListSourceKinds(
	ctx context.Context,
) ([]source.SourceKind, error) {
	if err := a.check(ctx); err != nil {
		return nil, err
	}
	return a.components.Sources.Kinds(), nil
}

func (a *API) Close() error {
	if a == nil {
		return nil
	}
	a.closeOnce.Do(func() {
		a.closed.Store(true)
		if a.components != nil {
			a.closeErr = a.components.Close()
		}
		a.resources = nil
		a.components = nil
	})
	return a.closeErr
}

func (a *API) check(ctx context.Context) error {
	if a == nil ||
		a.closed.Load() ||
		a.components == nil ||
		a.components.Roots == nil ||
		a.components.Sources == nil {
		return basespec.ErrClosed
	}
	if ctx == nil {
		return fmt.Errorf(
			"%w: artifact store API context is nil",
			basespec.ErrInvalid,
		)
	}
	return ctx.Err()
}
