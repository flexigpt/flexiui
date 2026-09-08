package bundle

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/schema"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
	mcpDomainPolicy "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain/policy"
	mcpDomainServer "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain/server"
)

type BundleExtension struct {
	Servers  map[string]mcpDomainServer.ServerExtension `json:"servers,omitempty"`
	Policies map[string]mcpDomainPolicy.PolicyDocument  `json:"policies,omitempty"`
}

type BundleDocument struct {
	Kind          collection.CollectionKind `json:"kind"`
	SchemaID      schema.SchemaID           `json:"schemaID"`
	SchemaVersion string                    `json:"schemaVersion"`
	Digest        cryptoutil.Digest         `json:"digest,omitempty"`

	LogicalName    basespec.LogicalName    `json:"logicalName"`
	LogicalVersion basespec.LogicalVersion `json:"logicalVersion,omitempty"`
	DisplayName    string                  `json:"displayName,omitempty"`
	Description    string                  `json:"description,omitempty"`
	Labels         map[string]string       `json:"labels,omitempty"`

	MCPServers      map[string]mcpDomainServer.CoreServer `json:"mcpServers"`
	BundleExtension BundleExtension                       `json:"bundleExtension"`
}

func (value BundleDocument) Validate() error {
	return validateDocument(value)
}
