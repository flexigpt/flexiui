package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	mcpConsumerAPI "github.com/flexigpt/flexigpt-app/internal/mcp/store/consumerapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
)

type MCPStoreWrapper struct {
	api   *mcpConsumerAPI.API
	roots compositionapi.RootAPI
}

func withMCPStore[T any](
	w *MCPStoreWrapper,
	fn func(*mcpConsumerAPI.API) (T, error),
) (T, error) {
	return middleware.WithRecoveryResp(func() (T, error) {
		var zero T
		if w == nil || w.api == nil {
			return zero, basespec.ErrClosed
		}
		return fn(w.api)
	})
}

func (w *MCPStoreWrapper) CreateMCPBundle(
	request *mcpConsumerAPI.CreateMCPBundleRequest,
) (*mcpConsumerAPI.CreateMCPBundleResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.CreateMCPBundleResponse, error) {
		return api.CreateMCPBundle(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) GetMCPBundle(
	request *mcpConsumerAPI.GetMCPBundleRequest,
) (*mcpConsumerAPI.GetMCPBundleResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.GetMCPBundleResponse, error) {
		return api.GetMCPBundle(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) ListMCPBundles(
	request *mcpConsumerAPI.ListMCPBundlesRequest,
) (*mcpConsumerAPI.ListMCPBundlesResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.ListMCPBundlesResponse, error) {
		return api.ListMCPBundles(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) ListMCPBundlesForManagement() ([]mcpConsumerAPI.Bundle, error) {
	return middleware.WithRecoveryResp(func() ([]mcpConsumerAPI.Bundle, error) {
		if w == nil || w.api == nil || w.roots == nil {
			return nil, basespec.ErrClosed
		}

		roots, err := w.roots.List(context.Background())
		if err != nil {
			return nil, err
		}

		bundles := make([]mcpConsumerAPI.Bundle, 0)
		for _, rootValue := range roots {
			response, err := w.api.ListMCPBundles(
				context.Background(),
				&mcpConsumerAPI.ListMCPBundlesRequest{
					RootID: rootValue.ID,
				},
			)
			if err != nil {
				return nil, err
			}
			if response == nil || response.Body == nil {
				return nil, fmt.Errorf(
					"%w: MCP Bundle list returned no response body for Root %q",
					basespec.ErrInvalid,
					rootValue.ID,
				)
			}
			bundles = append(bundles, response.Body.Bundles...)
		}

		sort.Slice(bundles, func(left, right int) bool {
			if bundles[left].Collection.RootID != bundles[right].Collection.RootID {
				return bundles[left].Collection.RootID < bundles[right].Collection.RootID
			}
			return bundles[left].Collection.ID < bundles[right].Collection.ID
		})
		return bundles, nil
	})
}

func (w *MCPStoreWrapper) GetMCPServerSchemaIdentity() (
	mcpConsumerAPI.MCPServerSchemaIdentity,
	error,
) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (mcpConsumerAPI.MCPServerSchemaIdentity, error) {
		return api.GetMCPServerSchemaIdentity(context.Background())
	})
}

func (w *MCPStoreWrapper) GetMCPBundleDocument(
	request *mcpConsumerAPI.GetMCPBundleDocumentRequest,
) (*mcpConsumerAPI.GetMCPBundleDocumentResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.GetMCPBundleDocumentResponse, error) {
		return api.GetMCPBundleDocument(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) ListMCPBundleServers(
	request *mcpConsumerAPI.ListMCPBundleServersRequest,
) (*mcpConsumerAPI.ListMCPBundleServersResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.ListMCPBundleServersResponse, error) {
		return api.ListMCPBundleServers(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) ListMCPBundlePolicies(
	request *mcpConsumerAPI.ListMCPBundlePoliciesRequest,
) (*mcpConsumerAPI.ListMCPBundlePoliciesResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.ListMCPBundlePoliciesResponse, error) {
		return api.ListMCPBundlePolicies(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) GetMCPServerInstallation(
	request *mcpConsumerAPI.GetMCPServerInstallationRequest,
) (*mcpConsumerAPI.GetMCPServerInstallationResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.GetMCPServerInstallationResponse, error) {
		return api.GetMCPServerInstallation(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) InspectMCPServer(
	request *mcpConsumerAPI.InspectMCPServerRequest,
) (*mcpConsumerAPI.InspectMCPServerResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.InspectMCPServerResponse, error) {
		return api.InspectMCPServer(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) InspectMCPPolicy(
	request *mcpConsumerAPI.InspectMCPPolicyRequest,
) (*mcpConsumerAPI.InspectMCPPolicyResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.InspectMCPPolicyResponse, error) {
		return api.InspectMCPPolicy(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) GetMCPBundleInstallation(
	request *mcpConsumerAPI.GetMCPBundleInstallationRequest,
) (*mcpConsumerAPI.GetMCPBundleInstallationResponse, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (*mcpConsumerAPI.GetMCPBundleInstallationResponse, error) {
		return api.GetMCPBundleInstallation(context.Background(), request)
	})
}

func (w *MCPStoreWrapper) UpdateBundleEnabled(
	ref collection.CollectionRef,
	expectedRevision uint64,
	enabled bool,
) (mcpConsumerAPI.Bundle, error) {
	return withMCPStore(w, func(api *mcpConsumerAPI.API) (mcpConsumerAPI.Bundle, error) {
		return api.UpdateBundleEnabled(
			context.Background(),
			ref,
			expectedRevision,
			enabled,
		)
	})
}

func (w *MCPStoreWrapper) close() {
	if w == nil {
		return
	}
	w.api = nil
	w.roots = nil
}
