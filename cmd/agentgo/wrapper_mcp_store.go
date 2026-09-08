package main

import (
	"context"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	mcpStore "github.com/flexigpt/flexigpt-app/internal/mcp/store"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
)

// MCPStoreWrapper exposes only pure MCP Store operations.
//
// Runtime-affecting MCP changes belong to MCPAggregateWrapper. Runtime
// connection operations belong to MCPRuntimeWrapper.
type MCPStoreWrapper struct {
	api *mcpStore.StoreAPI

	builtInInstaller artifactbuiltin.HydrationInstaller
}

func (w *MCPStoreWrapper) CreateMCPBundle(
	request *mcpStore.CreateMCPBundleRequest,
) (*mcpStore.CreateMCPBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.CreateMCPBundleResponse, error) {
			return w.api.CreateMCPBundle(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) GetMCPBundle(
	request *mcpStore.GetMCPBundleRequest,
) (*mcpStore.GetMCPBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.GetMCPBundleResponse, error) {
			return w.api.GetMCPBundle(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) ListMCPBundles(
	request *mcpStore.ListMCPBundlesRequest,
) (*mcpStore.ListMCPBundlesResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.ListMCPBundlesResponse, error) {
			return w.api.ListMCPBundles(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) GetMCPBundleDocument(
	request *mcpStore.GetMCPBundleDocumentRequest,
) (*mcpStore.GetMCPBundleDocumentResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.GetMCPBundleDocumentResponse, error) {
			return w.api.GetMCPBundleDocument(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) ListMCPBundleServers(
	request *mcpStore.ListMCPBundleServersRequest,
) (*mcpStore.ListMCPBundleServersResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.ListMCPBundleServersResponse, error) {
			return w.api.ListMCPBundleServers(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) ListMCPBundlePolicies(
	request *mcpStore.ListMCPBundlePoliciesRequest,
) (*mcpStore.ListMCPBundlePoliciesResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.ListMCPBundlePoliciesResponse, error) {
			return w.api.ListMCPBundlePolicies(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) GetMCPServerInstallation(
	request *mcpStore.GetMCPServerInstallationRequest,
) (*mcpStore.GetMCPServerInstallationResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.GetMCPServerInstallationResponse, error) {
			return w.api.GetMCPServerInstallation(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) InspectMCPServer(
	request *mcpStore.InspectMCPServerRequest,
) (*mcpStore.InspectMCPServerResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.InspectMCPServerResponse, error) {
			return w.api.InspectMCPServer(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) InspectMCPPolicy(
	request *mcpStore.InspectMCPPolicyRequest,
) (*mcpStore.InspectMCPPolicyResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.InspectMCPPolicyResponse, error) {
			return w.api.InspectMCPPolicy(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) GetMCPBundleInstallation(
	request *mcpStore.GetMCPBundleInstallationRequest,
) (*mcpStore.GetMCPBundleInstallationResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpStore.GetMCPBundleInstallationResponse, error) {
			return w.api.GetMCPBundleInstallation(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) close() {
	if w == nil {
		return
	}
	w.api = nil
	w.builtInInstaller = nil
}
