package store

import (
	"context"
	"errors"
)

// StoreAPI is the transport-facing MCP Store facade.
//
// It exposes only Store queries and Store-only creation. Mutations that must
// invalidate a live MCP runtime remain MCP Aggregate operations.
type StoreAPI struct {
	service *API
}

func NewStoreAPI(service *API) (*StoreAPI, error) {
	if service == nil {
		return nil, wrapStoreError(
			"facade initialization",
			errors.New("MCP Store service is required"),
		)
	}
	return &StoreAPI{service: service}, nil
}

func (a *StoreAPI) CreateMCPBundle(
	ctx context.Context,
	request *CreateMCPBundleRequest,
) (*CreateMCPBundleResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"MCP Bundle creation",
	); err != nil {
		return nil, err
	}

	value, err := a.service.Create(ctx, *request.Body)
	if err != nil {
		return nil, wrapStoreError("create Bundle", err)
	}
	return &CreateMCPBundleResponse{Body: &value}, nil
}

func (a *StoreAPI) GetMCPBundle(
	ctx context.Context,
	request *GetMCPBundleRequest,
) (*GetMCPBundleResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"MCP Bundle get",
	); err != nil {
		return nil, err
	}

	value, err := a.service.Get(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("get Bundle", err)
	}
	return &GetMCPBundleResponse{Body: &value}, nil
}

func (a *StoreAPI) ListMCPBundles(
	ctx context.Context,
	request *ListMCPBundlesRequest,
) (*ListMCPBundlesResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"MCP Bundle list",
	); err != nil {
		return nil, err
	}

	values, err := a.service.List(ctx, request.RootID)
	if err != nil {
		return nil, wrapStoreError("list Bundles", err)
	}
	return &ListMCPBundlesResponse{
		Body: &ListMCPBundlesResponseBody{
			Bundles: values,
		},
	}, nil
}

func (a *StoreAPI) GetMCPBundleDocument(
	ctx context.Context,
	request *GetMCPBundleDocumentRequest,
) (*GetMCPBundleDocumentResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"MCP Bundle document get",
	); err != nil {
		return nil, err
	}

	value, err := a.service.GetDocument(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("get Bundle document", err)
	}
	return &GetMCPBundleDocumentResponse{Body: &value}, nil
}

func (a *StoreAPI) ListMCPBundleServers(
	ctx context.Context,
	request *ListMCPBundleServersRequest,
) (*ListMCPBundleServersResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"MCP Bundle Server list",
	); err != nil {
		return nil, err
	}

	values, err := a.service.ListServers(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("list Bundle Servers", err)
	}
	return &ListMCPBundleServersResponse{
		Body: &ListMCPBundleServersResponseBody{
			Servers: values,
		},
	}, nil
}

func (a *StoreAPI) ListMCPBundlePolicies(
	ctx context.Context,
	request *ListMCPBundlePoliciesRequest,
) (*ListMCPBundlePoliciesResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"MCP Bundle Policy list",
	); err != nil {
		return nil, err
	}

	values, err := a.service.ListPolicies(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("list Bundle Policies", err)
	}
	return &ListMCPBundlePoliciesResponse{
		Body: &ListMCPBundlePoliciesResponseBody{
			Policies: values,
		},
	}, nil
}

func (a *StoreAPI) GetMCPServerInstallation(
	ctx context.Context,
	request *GetMCPServerInstallationRequest,
) (*GetMCPServerInstallationResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"MCP Server installation get",
	); err != nil {
		return nil, err
	}

	value, err := a.service.GetServerInstallation(ctx, request.Server)
	if err != nil {
		return nil, wrapStoreError("get Server installation", err)
	}
	return &GetMCPServerInstallationResponse{Body: &value}, nil
}

func (a *StoreAPI) InspectMCPServer(
	ctx context.Context,
	request *InspectMCPServerRequest,
) (*InspectMCPServerResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"MCP Server inspection",
	); err != nil {
		return nil, err
	}

	value, err := a.service.InspectMCPServer(ctx, request.Server)
	if err != nil {
		return nil, wrapStoreError("inspect Server", err)
	}
	return &InspectMCPServerResponse{Body: &value}, nil
}

func (a *StoreAPI) InspectMCPPolicy(
	ctx context.Context,
	request *InspectMCPPolicyRequest,
) (*InspectMCPPolicyResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"MCP Policy inspection",
	); err != nil {
		return nil, err
	}

	value, err := a.service.InspectMCPPolicy(ctx, request.Policy)
	if err != nil {
		return nil, wrapStoreError("inspect Policy", err)
	}
	return &InspectMCPPolicyResponse{Body: &value}, nil
}

func (a *StoreAPI) GetMCPBundleInstallation(
	ctx context.Context,
	request *GetMCPBundleInstallationRequest,
) (*GetMCPBundleInstallationResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"MCP Bundle installation get",
	); err != nil {
		return nil, err
	}

	value, err := a.service.GetBundleInstallation(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("get Bundle installation", err)
	}
	return &GetMCPBundleInstallationResponse{Body: &value}, nil
}
