package aggregate

import workspaceConsumerAPI "github.com/flexigpt/flexigpt-app/internal/workspace/store/consumerapi"

// Aggregate transport types are aliases. Storage transport models remain
// defined exactly once in store/consumerapi.
type (
	WorkspaceRef = workspaceConsumerAPI.WorkspaceRef

	WorkspaceArtifactView          = workspaceConsumerAPI.WorkspaceArtifactView
	WorkspaceContextContribution   = workspaceConsumerAPI.WorkspaceContextContribution
	WorkspaceContextDecision       = workspaceConsumerAPI.WorkspaceContextDecision
	WorkspaceContextInspectionView = workspaceConsumerAPI.WorkspaceContextInspectionView
	WorkspaceContextLoadPlan       = workspaceConsumerAPI.WorkspaceContextLoadPlan
	WorkspaceContextView           = workspaceConsumerAPI.WorkspaceContextView
	WorkspaceSkillLoadView         = workspaceConsumerAPI.WorkspaceSkillLoadView
	WorkspaceSkillView             = workspaceConsumerAPI.WorkspaceSkillView
	WorkspaceSkillSummary          = workspaceConsumerAPI.WorkspaceSkillSummary
	WorkspaceSkillArgument         = workspaceConsumerAPI.WorkspaceSkillArgument

	ListWorkspaceContextsRequest      = workspaceConsumerAPI.ListWorkspaceContextsRequest
	ListWorkspaceContextsResponseBody = workspaceConsumerAPI.ListWorkspaceContextsResponseBody
	ListWorkspaceContextsResponse     = workspaceConsumerAPI.ListWorkspaceContextsResponse

	LoadWorkspaceContextsRequestBody = workspaceConsumerAPI.LoadWorkspaceContextsRequestBody
	LoadWorkspaceContextsRequest     = workspaceConsumerAPI.LoadWorkspaceContextsRequest
	LoadWorkspaceContextsResponse    = workspaceConsumerAPI.LoadWorkspaceContextsResponse

	ComposeWorkspaceContextRequestBody = workspaceConsumerAPI.ComposeWorkspaceContextRequestBody
	ComposeWorkspaceContextRequest     = workspaceConsumerAPI.ComposeWorkspaceContextRequest
	ComposeWorkspaceContextResponse    = workspaceConsumerAPI.ComposeWorkspaceContextResponse

	ListWorkspaceSkillsRequest      = workspaceConsumerAPI.ListWorkspaceSkillsRequest
	ListWorkspaceSkillsResponseBody = workspaceConsumerAPI.ListWorkspaceSkillsResponseBody
	ListWorkspaceSkillsResponse     = workspaceConsumerAPI.ListWorkspaceSkillsResponse

	LoadWorkspaceSkillsRequestBody = workspaceConsumerAPI.LoadWorkspaceSkillsRequestBody
	LoadWorkspaceSkillsRequest     = workspaceConsumerAPI.LoadWorkspaceSkillsRequest
	LoadWorkspaceSkillsResponse    = workspaceConsumerAPI.LoadWorkspaceSkillsResponse

	GetWorkspaceRequest  = workspaceConsumerAPI.GetWorkspaceRequest
	GetWorkspaceResponse = workspaceConsumerAPI.GetWorkspaceResponse

	SetWorkspaceArtifactRuntimeDisabledRequest     = workspaceConsumerAPI.SetWorkspaceArtifactRuntimeDisabledRequest
	SetWorkspaceArtifactRuntimeDisabledRequestBody = workspaceConsumerAPI.SetWorkspaceArtifactRuntimeDisabledRequestBody
	SetWorkspaceArtifactRuntimeDisabledResponse    = workspaceConsumerAPI.SetWorkspaceArtifactRuntimeDisabledResponse

	ContextDocument   = workspaceConsumerAPI.ContextDocument
	ContextInspection = workspaceConsumerAPI.ContextInspection
	ContextLoadPlan   = workspaceConsumerAPI.ContextLoadPlan
)

func requireRequestBody[T any](
	request *T,
	bodyPresent bool,
	requireBody bool,
	subject string,
) error {
	return workspaceConsumerAPI.RequireRequestBody(
		request,
		bodyPresent,
		requireBody,
		subject,
	)
}
