package main

import (
	"context"

	mcpConsumerAPI "github.com/flexigpt/flexigpt-app/internal/mcp/store/consumerapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
)

// MCPStoreWrapper exposes only pure MCP Store operations.
//
// Runtime-affecting MCP changes belong to MCPAggregateWrapper. Runtime
// connection operations belong to MCPRuntimeWrapper.
type MCPStoreWrapper struct {
	api *mcpConsumerAPI.API
}

func (w *MCPStoreWrapper) CreateMCPBundle(
	request *mcpConsumerAPI.CreateMCPBundleRequest,
) (*mcpConsumerAPI.CreateMCPBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.CreateMCPBundleResponse, error) {
			return w.api.CreateMCPBundle(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) GetMCPBundle(
	request *mcpConsumerAPI.GetMCPBundleRequest,
) (*mcpConsumerAPI.GetMCPBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.GetMCPBundleResponse, error) {
			return w.api.GetMCPBundle(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) ListMCPBundles(
	request *mcpConsumerAPI.ListMCPBundlesRequest,
) (*mcpConsumerAPI.ListMCPBundlesResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.ListMCPBundlesResponse, error) {
			return w.api.ListMCPBundles(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) GetMCPBundleDocument(
	request *mcpConsumerAPI.GetMCPBundleDocumentRequest,
) (*mcpConsumerAPI.GetMCPBundleDocumentResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.GetMCPBundleDocumentResponse, error) {
			return w.api.GetMCPBundleDocument(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) ListMCPBundleServers(
	request *mcpConsumerAPI.ListMCPBundleServersRequest,
) (*mcpConsumerAPI.ListMCPBundleServersResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.ListMCPBundleServersResponse, error) {
			return w.api.ListMCPBundleServers(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) ListMCPBundlePolicies(
	request *mcpConsumerAPI.ListMCPBundlePoliciesRequest,
) (*mcpConsumerAPI.ListMCPBundlePoliciesResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.ListMCPBundlePoliciesResponse, error) {
			return w.api.ListMCPBundlePolicies(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) GetMCPServerInstallation(
	request *mcpConsumerAPI.GetMCPServerInstallationRequest,
) (*mcpConsumerAPI.GetMCPServerInstallationResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.GetMCPServerInstallationResponse, error) {
			return w.api.GetMCPServerInstallation(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) InspectMCPServer(
	request *mcpConsumerAPI.InspectMCPServerRequest,
) (*mcpConsumerAPI.InspectMCPServerResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.InspectMCPServerResponse, error) {
			return w.api.InspectMCPServer(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) InspectMCPPolicy(
	request *mcpConsumerAPI.InspectMCPPolicyRequest,
) (*mcpConsumerAPI.InspectMCPPolicyResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.InspectMCPPolicyResponse, error) {
			return w.api.InspectMCPPolicy(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) GetMCPBundleInstallation(
	request *mcpConsumerAPI.GetMCPBundleInstallationRequest,
) (*mcpConsumerAPI.GetMCPBundleInstallationResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*mcpConsumerAPI.GetMCPBundleInstallationResponse, error) {
			return w.api.GetMCPBundleInstallation(ctx, request)
		},
	)
}

func (w *MCPStoreWrapper) close() {
	if w == nil {
		return
	}
	w.api = nil
}
