package consumerapi

import (
	"encoding/json"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	mcpDomain "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain"
	mcpDomainServer "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain/server"
)

type CreateMCPBundleBody struct {
	RootID           root.RootID
	CollectionID     collection.CollectionID
	SourceID         source.SourceID
	SourceStorageKey basespec.StorageKey
	Document         json.RawMessage
	Registrations    []Registration
}

type CreateMCPBundleRequest struct {
	Body *CreateMCPBundleBody `json:"body"`
}

type CreateMCPBundleResponse struct {
	Body *Bundle `json:"body"`
}

type GetMCPBundleRequest struct {
	Bundle collection.CollectionRef `json:"bundle"`
}

type GetMCPBundleResponse struct {
	Body *Bundle `json:"body"`
}

type ListMCPBundlesRequest struct {
	RootID root.RootID `json:"rootID"`
}

type ListMCPBundlesResponseBody struct {
	Bundles []Bundle `json:"bundles"`
}

type ListMCPBundlesResponse struct {
	Body *ListMCPBundlesResponseBody `json:"body"`
}

type GetMCPBundleDocumentRequest struct {
	Bundle collection.CollectionRef `json:"bundle"`
}

type GetMCPBundleDocumentResponse struct {
	Body *mcpDomain.BundleDocument `json:"body"`
}

type ListMCPBundleServersRequest struct {
	Bundle collection.CollectionRef `json:"bundle"`
}

type ListMCPBundleServersResponseBody struct {
	Servers []artifact.Artifact `json:"servers"`
}

type ListMCPBundleServersResponse struct {
	Body *ListMCPBundleServersResponseBody `json:"body"`
}

type ListMCPBundlePoliciesRequest struct {
	Bundle collection.CollectionRef `json:"bundle"`
}

type ListMCPBundlePoliciesResponseBody struct {
	Policies []artifact.Artifact `json:"policies"`
}

type ListMCPBundlePoliciesResponse struct {
	Body *ListMCPBundlePoliciesResponseBody `json:"body"`
}

type GetMCPServerInstallationRequest struct {
	Server artifact.ArtifactRef `json:"server"`
}

type GetMCPServerInstallationResponse struct {
	Body *ServerInstallationView `json:"body"`
}

type InspectMCPServerRequest struct {
	Server artifact.ArtifactRef `json:"server"`
}

type InspectMCPServerResponse struct {
	Body *mcpDomainServer.Resolved `json:"body"`
}

type InspectMCPPolicyRequest struct {
	Policy artifact.ArtifactRef `json:"policy"`
}

type InspectMCPPolicyResponse struct {
	Body *PolicyView `json:"body"`
}

type GetMCPBundleInstallationRequest struct {
	Bundle collection.CollectionRef `json:"bundle"`
}

type GetMCPBundleInstallationResponse struct {
	Body *BundleInstallationView `json:"body"`
}

func requireStoreRequest[T any](
	request *T,
	requireBody bool,
	bodyPresent bool,
	subject string,
) error {
	if request == nil {
		return fmt.Errorf(
			"%w: %s request is required",
			basespec.ErrInvalid,
			subject,
		)
	}
	if requireBody && !bodyPresent {
		return fmt.Errorf(
			"%w: %s request body is required",
			basespec.ErrInvalid,
			subject,
		)
	}
	return nil
}

func wrapStoreError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("MCP Store %s: %w", operation, err)
}
