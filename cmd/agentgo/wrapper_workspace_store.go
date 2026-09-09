package main

import (
	"context"
	"errors"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
	workspaceAggregate "github.com/flexigpt/flexigpt-app/internal/workspace/aggregate"
	workspaceConsumerAPI "github.com/flexigpt/flexigpt-app/internal/workspace/store/consumerapi"
)

type WorkspaceStoreWrapper struct {
	api *workspaceConsumerAPI.StoreAPI
}

func InitWorkspaceWrappers(
	storeWrapper *WorkspaceStoreWrapper,
	runtimeWrapper *WorkspaceRuntimeWrapper,
	aggregateWrapper *WorkspaceAggregateWrapper,
	sources compositionapi.SourceAPI,
	collections compositionapi.CollectionAPI,
	artifacts compositionapi.ArtifactAPI,
	catalogs compositionapi.CatalogAPI,
	resources compositionapi.ResourceAPI,
) error {
	if storeWrapper == nil ||
		runtimeWrapper == nil ||
		aggregateWrapper == nil {
		return errors.New("workspace wrapper dependencies are incomplete")
	}

	storeAPI, err := workspaceConsumerAPI.NewStoreAPI(
		sources,
		collections,
		artifacts,
		catalogs,
		resources,
		workspaceConsumerAPI.DefaultConfig(),
	)
	if err != nil {
		return err
	}

	runtimeAPI, err := workspaceAggregate.NewRuntimeAPI(
		storeAPI.ContextService(),
		storeAPI.SkillAdapter(),
	)
	if err != nil {
		return err
	}

	aggregateAPI, err := workspaceAggregate.NewAggregateAPI(
		storeAPI,
		runtimeAPI,
		storeAPI,
	)
	if err != nil {
		return err
	}

	storeWrapper.api = storeAPI
	runtimeWrapper.api = runtimeAPI
	aggregateWrapper.api = aggregateAPI
	return nil
}

func withWorkspaceStore[T any](
	w *WorkspaceStoreWrapper,
	fn func(*workspaceConsumerAPI.StoreAPI) (T, error),
) (T, error) {
	return middleware.WithRecoveryResp(func() (T, error) {
		var zero T
		if w == nil || w.api == nil {
			return zero, errors.New("workspace Store API is unavailable")
		}
		return fn(w.api)
	})
}

func (w *WorkspaceStoreWrapper) CreateFilesystemWorkspace(
	input workspaceConsumerAPI.CreateFilesystemWorkspaceInput,
) (workspaceConsumerAPI.WorkspaceView, error) {
	return withWorkspaceStore(w, func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceView, error) {
		return api.CreateFilesystemWorkspace(context.Background(), input)
	})
}

func (w *WorkspaceStoreWrapper) CreateEmptyWorkspace(
	input workspaceConsumerAPI.CreateEmptyWorkspaceInput,
) (workspaceConsumerAPI.WorkspaceView, error) {
	return withWorkspaceStore(w, func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceView, error) {
		return api.CreateEmptyWorkspace(context.Background(), input)
	})
}

func (w *WorkspaceStoreWrapper) UpdateWorkspace(
	workspace workspaceConsumerAPI.WorkspaceRef,
	input workspaceConsumerAPI.UpdateWorkspaceInput,
) (workspaceConsumerAPI.WorkspaceView, error) {
	return withWorkspaceStore(w, func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceView, error) {
		return api.UpdateWorkspace(context.Background(), workspace, input)
	})
}

func (w *WorkspaceStoreWrapper) SetWorkspacePrimarySource(
	workspace workspaceConsumerAPI.WorkspaceRef,
	input workspaceConsumerAPI.SetWorkspacePrimarySourceInput,
) (workspaceConsumerAPI.WorkspaceView, error) {
	return withWorkspaceStore(w, func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceView, error) {
		return api.SetWorkspacePrimarySource(context.Background(), workspace, input)
	})
}

func (w *WorkspaceStoreWrapper) AttachWorkspaceSource(
	workspace workspaceConsumerAPI.WorkspaceRef,
	input workspaceConsumerAPI.AttachWorkspaceSourceInput,
) (workspaceConsumerAPI.WorkspaceView, error) {
	return withWorkspaceStore(w, func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceView, error) {
		return api.AttachWorkspaceSource(context.Background(), workspace, input)
	})
}

func (w *WorkspaceStoreWrapper) UpdateWorkspaceAttachment(
	workspace workspaceConsumerAPI.WorkspaceRef,
	input workspaceConsumerAPI.UpdateWorkspaceAttachmentInput,
) (workspaceConsumerAPI.WorkspaceView, error) {
	return withWorkspaceStore(w, func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceView, error) {
		return api.UpdateWorkspaceAttachment(context.Background(), workspace, input)
	})
}

func (w *WorkspaceStoreWrapper) DetachWorkspaceSource(
	workspace workspaceConsumerAPI.WorkspaceRef,
	input workspaceConsumerAPI.DetachWorkspaceSourceInput,
) (workspaceConsumerAPI.WorkspaceView, error) {
	return withWorkspaceStore(w, func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceView, error) {
		return api.DetachWorkspaceSource(context.Background(), workspace, input)
	})
}

func (w *WorkspaceStoreWrapper) RefreshWorkspace(
	workspace workspaceConsumerAPI.WorkspaceRef,
) (workspaceConsumerAPI.WorkspaceRefreshResult, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceRefreshResult, error) {
			return api.RefreshWorkspace(context.Background(), workspace)
		},
	)
}

func (w *WorkspaceStoreWrapper) RetireWorkspace(
	workspace workspaceConsumerAPI.WorkspaceRef,
	expectedRevision uint64,
) (workspaceConsumerAPI.RetireWorkspaceResult, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.RetireWorkspaceResult, error) {
			return api.RetireWorkspace(context.Background(), workspace, expectedRevision)
		},
	)
}

func (w *WorkspaceStoreWrapper) PurgeWorkspace(
	workspace workspaceConsumerAPI.WorkspaceRef,
	expectedRevision uint64,
) (workspaceConsumerAPI.WorkspaceRef, error) {
	return withWorkspaceStore(w, func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceRef, error) {
		return api.PurgeWorkspace(context.Background(), workspace, expectedRevision)
	})
}

func (w *WorkspaceStoreWrapper) ListWorkspaceSourcesForManagement() (
	[]workspaceConsumerAPI.WorkspaceSourceSummary,
	error,
) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) ([]workspaceConsumerAPI.WorkspaceSourceSummary, error) {
			return api.ListWorkspaceSourcesForManagement(context.Background())
		},
	)
}

func (w *WorkspaceStoreWrapper) RegisterWorkspaceDirectory(
	workspace workspaceConsumerAPI.WorkspaceRef,
	input workspaceConsumerAPI.RegisterWorkspaceDirectoryInput,
) (workspaceConsumerAPI.WorkspaceDirectoryRegistrationResult, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceDirectoryRegistrationResult, error) {
			return api.RegisterWorkspaceDirectory(context.Background(), workspace, input)
		},
	)
}

func (w *WorkspaceStoreWrapper) SetWorkspaceSourceEnabled(
	sourceID string,
	expectedRevision uint64,
	enabled bool,
) (workspaceConsumerAPI.WorkspaceSourceSummary, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceSourceSummary, error) {
			return api.SetWorkspaceSourceEnabled(
				context.Background(),
				source.SourceID(sourceID),
				expectedRevision,
				enabled,
			)
		},
	)
}

func (w *WorkspaceStoreWrapper) AdoptWorkspaceOccurrence(
	workspace workspaceConsumerAPI.WorkspaceRef,
	input workspaceConsumerAPI.AdoptWorkspaceOccurrenceInput,
) (workspaceConsumerAPI.WorkspaceArtifactView, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceArtifactView, error) {
			return api.AdoptWorkspaceOccurrence(context.Background(), workspace, input)
		},
	)
}

func (w *WorkspaceStoreWrapper) PinWorkspaceArtifact(
	workspace workspaceConsumerAPI.WorkspaceRef,
	input workspaceConsumerAPI.PinWorkspaceArtifactInput,
) (workspaceConsumerAPI.WorkspaceArtifactView, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceArtifactView, error) {
			return api.PinWorkspaceArtifact(context.Background(), workspace, input)
		},
	)
}

func (w *WorkspaceStoreWrapper) ListWorkspaceSuppressions(
	workspace workspaceConsumerAPI.WorkspaceRef,
) ([]workspaceConsumerAPI.WorkspaceSuppressionView, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) ([]workspaceConsumerAPI.WorkspaceSuppressionView, error) {
			return api.ListWorkspaceSuppressions(context.Background(), workspace)
		},
	)
}

func (w *WorkspaceStoreWrapper) SuppressWorkspaceBinding(
	workspace workspaceConsumerAPI.WorkspaceRef,
	input workspaceConsumerAPI.SuppressWorkspaceBindingInput,
) (workspaceConsumerAPI.WorkspaceSuppressionView, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceSuppressionView, error) {
			return api.SuppressWorkspaceBinding(context.Background(), workspace, input)
		},
	)
}

func (w *WorkspaceStoreWrapper) UnsuppressWorkspaceBinding(
	workspace workspaceConsumerAPI.WorkspaceRef,
	binding artifact.SourceBinding,
	expectedRevision uint64,
) (workspaceConsumerAPI.UnsuppressWorkspaceBindingResult, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.UnsuppressWorkspaceBindingResult, error) {
			return api.UnsuppressWorkspaceBinding(
				context.Background(),
				workspace,
				binding,
				expectedRevision,
			)
		},
	)
}

func (w *WorkspaceStoreWrapper) SetWorkspaceArtifactEnabled(
	workspace workspaceConsumerAPI.WorkspaceRef,
	art artifact.ArtifactRef,
	expectedRevision uint64,
	enabled bool,
) (workspaceConsumerAPI.WorkspaceArtifactView, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.WorkspaceArtifactView, error) {
			return api.SetWorkspaceArtifactEnabled(
				context.Background(),
				workspace,
				art,
				expectedRevision,
				enabled,
			)
		},
	)
}

func (w *WorkspaceStoreWrapper) UnadoptWorkspaceArtifact(
	workspace workspaceConsumerAPI.WorkspaceRef,
	art artifact.ArtifactRef,
	input workspaceConsumerAPI.UnadoptWorkspaceArtifactInput,
) (workspaceConsumerAPI.UnadoptWorkspaceArtifactResult, error) {
	return withWorkspaceStore(
		w,
		func(api *workspaceConsumerAPI.StoreAPI) (workspaceConsumerAPI.UnadoptWorkspaceArtifactResult, error) {
			return api.UnadoptWorkspaceArtifact(context.Background(), workspace, art, input)
		},
	)
}

func (w *WorkspaceStoreWrapper) PurgeWorkspaceArtifact(
	workspace workspaceConsumerAPI.WorkspaceRef,
	art artifact.ArtifactRef,
	expectedRevision uint64,
) (artifact.ArtifactRef, error) {
	return withWorkspaceStore(w, func(api *workspaceConsumerAPI.StoreAPI) (artifact.ArtifactRef, error) {
		return api.PurgeWorkspaceArtifact(
			context.Background(),
			workspace,
			art,
			expectedRevision,
		)
	})
}

func (w *WorkspaceStoreWrapper) GetWorkspace(
	request *workspaceConsumerAPI.GetWorkspaceRequest,
) (*workspaceConsumerAPI.GetWorkspaceResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.GetWorkspaceResponse, error) {
			return w.api.GetWorkspace(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) ListWorkspaces(
	request *workspaceConsumerAPI.ListWorkspacesRequest,
) (*workspaceConsumerAPI.ListWorkspacesResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.ListWorkspacesResponse, error) {
			return w.api.ListWorkspaces(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) GetWorkspaceCatalog(
	request *workspaceConsumerAPI.GetWorkspaceCatalogRequest,
) (*workspaceConsumerAPI.GetWorkspaceCatalogResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.GetWorkspaceCatalogResponse, error) {
			return w.api.GetWorkspaceCatalog(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) GetWorkspaceArtifact(
	request *workspaceConsumerAPI.GetWorkspaceArtifactRequest,
) (*workspaceConsumerAPI.GetWorkspaceArtifactResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.GetWorkspaceArtifactResponse, error) {
			return w.api.GetWorkspaceArtifact(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) ListWorkspaceArtifacts(
	request *workspaceConsumerAPI.ListWorkspaceArtifactsRequest,
) (*workspaceConsumerAPI.ListWorkspaceArtifactsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*workspaceConsumerAPI.ListWorkspaceArtifactsResponse, error) {
			return w.api.ListWorkspaceArtifacts(ctx, request)
		},
	)
}

func (w *WorkspaceStoreWrapper) close() {
	if w == nil {
		return
	}
	w.api = nil
}
