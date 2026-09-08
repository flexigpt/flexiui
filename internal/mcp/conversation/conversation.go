package conversation

import (
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
	mcpServer "github.com/flexigpt/flexigpt-app/internal/mcp/runtime/server"
	mcpDomainPolicy "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain/policy"
)

type MCPAppModelContextUpdate struct {
	InstanceID string             `json:"instanceID,omitempty"`
	Server     mcpServer.ServerID `json:"server"`

	ResourceURI string `json:"resourceUri,omitempty"`

	Content           []mcpServer.MCPContent `json:"content,omitempty"`
	StructuredContent any                    `json:"structuredContent,omitempty"`
	UpdatedAt         string                 `json:"updatedAt,omitempty"`
	RawArguments      jsonutil.JSONRawString `json:"rawArguments,omitempty"`
}

type MCPToolSelection struct {
	Server           mcpServer.ServerID `json:"server"`
	ToolName         string             `json:"toolName"`
	ProviderToolName string             `json:"providerToolName,omitempty"`
	ChoiceID         string             `json:"choiceID,omitempty"`
	Digest           string             `json:"digest,omitempty"`

	ApprovalRule  *mcpDomainPolicy.MCPApprovalRule  `json:"approvalRule,omitempty"`
	ExecutionMode *mcpDomainPolicy.MCPExecutionMode `json:"executionMode,omitempty"`

	AppResourceURI string   `json:"appResourceUri,omitempty"`
	Visibility     []string `json:"visibility,omitempty"`
}

type MCPProviderToolMapping struct {
	Server mcpServer.ServerID `json:"server"`

	ProviderToolName string `json:"providerToolName"`
	ChoiceID         string `json:"choiceID"`

	ToolName   string `json:"toolName"`
	ToolDigest string `json:"toolDigest"`

	ApprovalRule   mcpDomainPolicy.MCPApprovalRule  `json:"approvalRule,omitempty"`
	ExecutionMode  mcpDomainPolicy.MCPExecutionMode `json:"executionMode,omitempty"`
	AppResourceURI string                           `json:"appResourceUri,omitempty"`
	Visibility     []string                         `json:"visibility,omitempty"`
}

type MCPToolExposure string

const (
	MCPToolExposureNone     MCPToolExposure = "none"
	MCPToolExposureAll      MCPToolExposure = "all"
	MCPToolExposureSelected MCPToolExposure = "selected"
)

type MCPServerSelection struct {
	Server mcpServer.ServerID `json:"server"`

	SnapshotDigest string `json:"snapshotDigest,omitempty"`

	ToolExposure  MCPToolExposure    `json:"toolExposure"` // none | all | selected
	SelectedTools []MCPToolSelection `json:"selectedTools,omitempty"`

	IncludeServerInstructions bool `json:"includeServerInstructions,omitempty"`
}

type MCPResourceTemplateSelection struct {
	mcpServer.MCPResourceTemplateRef

	ArgumentValues map[string]string `json:"argumentValues,omitempty"`
}

type MCPPromptSelection struct {
	mcpServer.MCPPromptRef

	ArgumentValues map[string]string `json:"argumentValues,omitempty"`
}

type MCPConversationContext struct {
	Servers           []MCPServerSelection           `json:"servers"`
	Resources         []mcpServer.MCPResourceRef     `json:"resources,omitempty"`
	ResourceTemplates []MCPResourceTemplateSelection `json:"resourceTemplates,omitempty"`
	Prompts           []MCPPromptSelection           `json:"prompts,omitempty"`
}
