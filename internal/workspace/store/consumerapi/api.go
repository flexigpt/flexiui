package consumerapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/flexigpt/agentskills-go/document"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
	"github.com/flexigpt/flexigpt-app/internal/skill/store/workspaceadapter"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
	"github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/artifactadapter"
	"github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/attachmentdata"
)

type StoreAPI struct {
	sources     compositionapi.SourceAPI
	collections compositionapi.CollectionAPI
	artifacts   compositionapi.ArtifactAPI
	catalogs    compositionapi.CatalogAPI
	resources   compositionapi.ResourceAPI
	workspace   *components
}

func NewStoreAPI(
	sources compositionapi.SourceAPI,
	collections compositionapi.CollectionAPI,
	artifacts compositionapi.ArtifactAPI,
	catalogs compositionapi.CatalogAPI,
	resources compositionapi.ResourceAPI,
	config Config,
) (*StoreAPI, error) {
	if sources == nil ||
		collections == nil ||
		artifacts == nil ||
		catalogs == nil ||
		resources == nil {
		return nil, fmt.Errorf(
			"%w: Workspace Store dependencies are incomplete",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}
	config = config.normalized()

	workspaceComponents, err := newComponents(
		sources,
		collections,
		artifacts,
		resources,
		config,
	)
	if err != nil {
		return nil, err
	}

	return &StoreAPI{
		sources:     sources,
		collections: collections,
		artifacts:   artifacts,
		catalogs:    catalogs,
		resources:   resources,
		workspace:   workspaceComponents,
	}, nil
}

func (a *StoreAPI) GetWorkspace(
	ctx context.Context,
	request *GetWorkspaceRequest,
) (*GetWorkspaceResponse, error) {
	if err := RequireRequestBody(
		request,
		false,
		false,
		"Workspace get",
	); err != nil {
		return nil, err
	}
	value, err := a.workspace.service.Get(ctx, request.Workspace)
	if err != nil {
		return nil, err
	}
	view, err := a.workspaceViewForAPI(ctx, value)
	if err != nil {
		return nil, err
	}
	return &GetWorkspaceResponse{Body: &view}, nil
}

func (a *StoreAPI) ListWorkspaces(
	ctx context.Context,
	request *ListWorkspacesRequest,
) (*ListWorkspacesResponse, error) {
	if err := RequireRequestBody(
		request,
		false,
		false,
		"Workspace list",
	); err != nil {
		return nil, err
	}

	// Listing follows the same single-Root contract as creation. In
	// particular, a stale UI cannot enumerate the protected built-in Root.
	values, err := a.workspace.service.List(ctx, a.workspace.workspaceRootID)
	if err != nil {
		return nil, err
	}
	output := make([]WorkspaceView, 0, len(values))
	for _, value := range values {
		view, err := a.workspaceViewForAPI(ctx, value)
		if err != nil {
			return nil, err
		}
		output = append(output, view)
	}
	return &ListWorkspacesResponse{
		Body: &ListWorkspacesResponseBody{Workspaces: output},
	}, nil
}

func (a *StoreAPI) GetWorkspaceCatalog(
	ctx context.Context,
	request *GetWorkspaceCatalogRequest,
) (*GetWorkspaceCatalogResponse, error) {
	if err := RequireRequestBody(
		request,
		false,
		false,
		"Workspace catalog get",
	); err != nil {
		return nil, err
	}
	value, err := a.workspace.query.Catalog(ctx, request.Workspace)
	if err != nil {
		return nil, err
	}
	output, err := a.workspaceCatalogViewForAPI(ctx, value)
	if err != nil {
		return nil, err
	}
	return &GetWorkspaceCatalogResponse{Body: &output}, nil
}

func (a *StoreAPI) GetWorkspaceArtifact(
	ctx context.Context,
	request *GetWorkspaceArtifactRequest,
) (*GetWorkspaceArtifactResponse, error) {
	if err := RequireRequestBody(
		request,
		false,
		false,
		"Workspace Artifact get",
	); err != nil {
		return nil, err
	}
	value, err := a.workspaceArtifact(
		ctx,
		request.Workspace,
		request.Artifact,
	)
	if err != nil {
		return nil, err
	}
	output := workspaceArtifactViewOf(value)
	return &GetWorkspaceArtifactResponse{Body: &output}, nil
}

func (a *StoreAPI) ListWorkspaceArtifacts(
	ctx context.Context,
	request *ListWorkspaceArtifactsRequest,
) (*ListWorkspaceArtifactsResponse, error) {
	if err := RequireRequestBody(
		request,
		false,
		false,
		"Workspace Artifact list",
	); err != nil {
		return nil, err
	}
	if _, err := a.workspace.service.Get(ctx, request.Workspace); err != nil {
		return nil, err
	}
	values, err := a.artifacts.ListByCollection(
		ctx,
		request.Workspace,
	)
	if err != nil {
		return nil, err
	}
	output := make([]WorkspaceArtifactView, 0, len(values))
	for _, value := range values {
		output = append(output, workspaceArtifactViewOf(value))
	}
	sort.Slice(output, func(left, right int) bool {
		if output[left].Name != output[right].Name {
			return output[left].Name < output[right].Name
		}
		return output[left].Artifact.ArtifactID < output[right].Artifact.ArtifactID
	})
	return &ListWorkspaceArtifactsResponse{
		Body: &ListWorkspaceArtifactsResponseBody{Artifacts: output},
	}, nil
}

func (a *StoreAPI) ContextService() ContextService {
	if a == nil || a.workspace == nil {
		return nil
	}
	return a.workspace.contextService
}

func (a *StoreAPI) SkillAdapter() *workspaceadapter.Adapter {
	return a.workspace.skillAdapter
}

func (a *StoreAPI) SetArtifactRuntimeDisabled(
	ctx context.Context,
	workspace workspaceDomain.WorkspaceRef,
	ref artifact.ArtifactRef,
	expectedRevision uint64,
	runtimeDisabled bool,
) (WorkspaceArtifactView, error) {
	current, err := a.workspaceArtifact(ctx, workspace, ref)
	if err != nil {
		return WorkspaceArtifactView{}, err
	}

	artifactData, err := artifactadapter.DecodeArtifactData(current.Data)
	if err != nil {
		return WorkspaceArtifactView{}, err
	}
	artifactData.RuntimeDisabled = runtimeDisabled

	data, err := artifactadapter.EncodeArtifactData(artifactData)
	if err != nil {
		return WorkspaceArtifactView{}, err
	}

	value, err := a.artifacts.UpdateData(
		ctx,
		ref,
		expectedRevision,
		data,
	)
	if err != nil {
		return WorkspaceArtifactView{}, err
	}

	return workspaceArtifactViewOf(value), nil
}

func (a *StoreAPI) workspaceArtifact(
	ctx context.Context,
	workspace collection.CollectionRef,
	ref artifact.ArtifactRef,
) (artifact.Artifact, error) {
	if _, err := a.workspace.service.Get(ctx, workspace); err != nil {
		return artifact.Artifact{}, err
	}
	value, err := a.artifacts.Get(ctx, ref)
	if err != nil {
		return artifact.Artifact{}, err
	}
	if value.RootID != workspace.RootID ||
		value.CollectionID != workspace.CollectionID {
		return artifact.Artifact{}, fmt.Errorf(
			"%w: Artifact %q does not belong to Workspace %q",
			workspaceDomain.ErrReferenceUnresolved,
			ref.ArtifactID,
			workspace.CollectionID,
		)
	}
	return value, nil
}

func (a *StoreAPI) workspaceViewForAPI(
	ctx context.Context,
	value workspaceDomain.Workspace,
) (WorkspaceView, error) {
	output, err := workspaceViewOf(value)
	if err != nil {
		return WorkspaceView{}, err
	}
	if err := a.enrichWorkspaceSourcePresentation(ctx, &output, value); err != nil {
		return WorkspaceView{}, err
	}
	return output, nil
}

func (a *StoreAPI) workspaceCatalogViewForAPI(
	ctx context.Context,
	value workspaceDomain.CatalogView,
) (WorkspaceCatalogView, error) {
	output, err := workspaceCatalogViewOf(value)
	if err != nil {
		return WorkspaceCatalogView{}, err
	}
	if err := a.enrichWorkspaceSourcePresentation(
		ctx,
		&output.Workspace,
		value.Workspace,
	); err != nil {
		return WorkspaceCatalogView{}, err
	}
	return output, nil
}

func (a *StoreAPI) enrichWorkspaceSourcePresentation(
	ctx context.Context,
	output *WorkspaceView,
	value workspaceDomain.Workspace,
) error {
	for index := range output.Attachments {
		attachment := &output.Attachments[index]
		var summaryFound bool

		for _, summary := range value.Sources {
			if summary.ID != attachment.SourceID {
				continue
			}
			attachment.SourceDisplayName = summary.DisplayName
			attachment.SourceKind = string(summary.Kind)
			summaryFound = true
			break
		}
		if !summaryFound {
			return fmt.Errorf(
				"%w: Workspace attachment source %q is unavailable",
				workspaceDomain.ErrInvalidWorkspace,
				attachment.SourceID,
			)
		}
		sourceKind := source.SourceKind(attachment.SourceKind)
		if !a.resources.SupportsLocalPath(sourceKind) {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}

		pathValue, err := a.resources.ResolveSourceLocalPath(
			ctx,
			value.Collection.RootID,
			attachment.SourceID,
			".",
		)
		if err != nil {
			attachment.Diagnostics = diagnostic.Append(
				attachment.Diagnostics,
				workspaceSourcePresentationDiagnostic(
					"workspace.source.path-unavailable",
					"the filesystem Source path is currently unavailable",
				),
			)
			continue
		}
		attachment.Path = pathValue
		if attachment.SourceID == output.PrimarySourceID {
			output.PrimaryPath = attachment.Path
		}
	}
	return nil
}

func workspaceSourcePresentationDiagnostic(
	code string,
	message string,
) diagnostic.Diagnostic {
	return diagnostic.Diagnostic{
		Severity: diagnostic.SeverityWarning,
		Code:     code,
		Message:  message,
	}
}

func workspaceCatalogViewOf(
	value workspaceDomain.CatalogView,
) (WorkspaceCatalogView, error) {
	workspaceValue, err := workspaceViewOf(value.Workspace)
	if err != nil {
		return WorkspaceCatalogView{}, err
	}
	output := WorkspaceCatalogView{
		Workspace:       workspaceValue,
		CatalogRevision: value.Catalog.Revision,
		CatalogCurrent:  value.CatalogCurrent,
		Diagnostics: diagnostic.Append(
			value.Catalog.Diagnostics,
			value.FreshnessDiagnostics...,
		),
		Resources:             make([]WorkspaceResourceView, 0, len(value.Resources)),
		Groups:                make([]WorkspaceResourceGroupView, 0, len(value.Groups)),
		Occurrences:           make([]WorkspaceOccurrenceView, 0, len(value.Catalog.Occurrences)),
		ValidOccurrences:      make([]WorkspaceOccurrenceView, 0),
		InvalidOccurrences:    make([]WorkspaceOccurrenceView, 0),
		MissingOccurrences:    make([]WorkspaceOccurrenceView, 0),
		UnrecordedOccurrences: make([]WorkspaceOccurrenceView, 0),
		UnresolvedArtifacts:   make([]WorkspaceArtifactView, 0, len(value.UnresolvedArtifacts)),

		UnrecordedCount:         len(value.Unrecorded),
		UnresolvedArtifactCount: len(value.UnresolvedArtifacts),
	}
	artifactsByOccurrence := make(map[string]artifact.Artifact, len(value.Resources))
	for _, resourceValue := range value.Resources {
		artifactView := workspaceArtifactViewOfResource(resourceValue)
		projected := WorkspaceResourceView{
			Artifact:         artifactView,
			DefinitionDigest: resourceValue.Definition.Digest,
			SourceID:         resourceValue.Source.ID,
			Locator:          resourceValue.Artifact.Binding.Locator,
			CatalogCurrent:   resourceValue.CatalogCurrent,
			ProjectionValid:  resourceValue.ProjectionValid,
			Diagnostics: diagnostic.Append(
				artifactView.Diagnostics,
				resourceValue.Diagnostics...,
			),
		}
		output.Resources = append(output.Resources, projected)
		artifactsByOccurrence[occurrenceViewKey(
			resourceValue.Artifact.Binding.SourceID,
			resourceValue.Artifact.Binding.Locator,
			resourceValue.Artifact.Binding.SubresourceLocator,
			resourceValue.Artifact.Kind,
		)] = resourceValue.Artifact
	}
	for _, localArtifact := range value.UnresolvedArtifacts {
		output.UnresolvedArtifacts = append(
			output.UnresolvedArtifacts,
			workspaceArtifactViewOf(localArtifact),
		)
		artifactsByOccurrence[occurrenceViewKey(
			localArtifact.Binding.SourceID,
			localArtifact.Binding.Locator,
			localArtifact.Binding.SubresourceLocator,
			localArtifact.Kind,
		)] = localArtifact
	}
	for _, occurrence := range value.Catalog.Occurrences {
		projected := workspaceOccurrenceViewOf(
			occurrence,
			artifactsByOccurrence,
		)
		output.Occurrences = append(output.Occurrences, projected)
		switch occurrence.State {
		case catalog.OccurrenceValid:
			output.ValidOccurrences = append(output.ValidOccurrences, projected)
		case catalog.OccurrenceInvalid:
			output.InvalidOccurrences = append(output.InvalidOccurrences, projected)
		case catalog.OccurrenceMissing:
			output.MissingOccurrences = append(output.MissingOccurrences, projected)
		default:
		}
		if !projected.Recorded {
			output.UnrecordedOccurrences = append(
				output.UnrecordedOccurrences,
				projected,
			)
		}
	}
	for _, group := range value.Groups {
		projected := WorkspaceResourceGroupView{
			Kind:       group.Kind,
			Resources:  make([]WorkspaceResourceView, 0, len(group.Resources)),
			Unrecorded: make([]WorkspaceOccurrenceView, 0, len(group.Unrecorded)),
		}
		for _, resourceValue := range group.Resources {
			artifactView := workspaceArtifactViewOfResource(resourceValue)
			projected.Resources = append(
				projected.Resources,
				WorkspaceResourceView{
					Artifact:         artifactView,
					DefinitionDigest: resourceValue.Definition.Digest,
					SourceID:         resourceValue.Source.ID,
					Locator:          resourceValue.Artifact.Binding.Locator,
					CatalogCurrent:   resourceValue.CatalogCurrent,
					ProjectionValid:  resourceValue.ProjectionValid,
					Diagnostics: diagnostic.Append(
						artifactView.Diagnostics,
						resourceValue.Diagnostics...,
					),
				},
			)
		}
		for _, occurrence := range group.Unrecorded {
			projected.Unrecorded = append(
				projected.Unrecorded,
				workspaceOccurrenceViewOf(occurrence, artifactsByOccurrence),
			)
		}
		output.Groups = append(output.Groups, projected)
	}
	return output, nil
}

func workspaceViewOf(value workspaceDomain.Workspace) (WorkspaceView, error) {
	output := WorkspaceView{
		Workspace:       value.Collection.Ref(),
		Revision:        value.Collection.Revision,
		DisplayName:     value.Collection.DisplayName,
		Description:     value.Collection.Description,
		Enabled:         value.Collection.Enabled,
		Mode:            value.Mode,
		PrimarySourceID: value.PrimarySourceID,
		Discovery:       workspaceDiscoveryOf(value.Data.Discovery),
		Attachments:     make([]WorkspaceAttachmentView, 0, len(value.Attachments)),
	}

	for _, attachment := range value.Attachments {
		settings, err := workspaceAttachmentSettingsOf(attachment.Data)
		if err != nil {
			return WorkspaceView{}, err
		}
		output.Attachments = append(output.Attachments, WorkspaceAttachmentView{
			SourceID: attachment.SourceID,
			Revision: attachment.Revision,
			Role:     attachment.Role,
			Enabled:  attachment.Enabled,
			Settings: settings,
		})
	}
	return output, nil
}

func workspaceAttachmentSettingsOf(
	raw json.RawMessage,
) (WorkspaceAttachmentSettings, error) {
	value, err := attachmentdata.DecodeAttachmentData(raw)
	if err != nil {
		return WorkspaceAttachmentSettings{}, fmt.Errorf(
			"%w: decode workspace attachment settings: %w",
			workspaceDomain.ErrInvalidWorkspace,
			err,
		)
	}

	return WorkspaceAttachmentSettings{
		Recursive:     cloneBool(value.Recursive),
		Authoritative: cloneBool(value.Authoritative),
	}, nil
}

func workspaceDiscoveryOf(value workspaceDomain.DiscoveryPreferences) WorkspaceDiscovery {
	output := WorkspaceDiscovery{
		AdditionalLocators: append(
			[]basespec.Locator(nil),
			value.AdditionalLocators...,
		),
		IncludeReadme: value.IncludeReadme,
	}
	for _, root := range value.AdditionalRoots {
		output.AdditionalRoots = append(output.AdditionalRoots, WorkspaceDiscoveryRoot{
			Root:            root.Root,
			Recursive:       root.Recursive,
			IncludePatterns: append([]string(nil), root.IncludePatterns...),
		})
	}
	return output
}

func workspaceOccurrenceViewOf(
	value catalog.Occurrence,
	artifacts map[string]artifact.Artifact,
) WorkspaceOccurrenceView {
	output := WorkspaceOccurrenceView{
		SourceID:            value.Key.SourceID,
		Locator:             value.Key.Locator,
		SubresourceLocator:  value.Key.SubresourceLocator,
		Kind:                value.Kind,
		LogicalName:         value.LogicalName,
		LogicalVersion:      value.LogicalVersion,
		DefinitionDigest:    cryptoutil.CloneDigest(value.DefinitionDigest),
		SourceContentDigest: cryptoutil.CloneDigest(value.SourceContentDigest),
		State:               string(value.State),
		Diagnostics:         diagnostic.Clone(value.Diagnostics),
	}
	if localArtifact, found := artifacts[occurrenceViewKey(
		value.Key.SourceID,
		value.Key.Locator,
		value.Key.SubresourceLocator,
		value.Kind,
	)]; found {
		artifactRef := localArtifact.Ref()
		output.Recorded = true
		output.Artifact = &artifactRef
	}
	return output
}

func ContextLoadPlanViewOf(
	value ContextLoadPlan,
) WorkspaceContextLoadPlan {
	output := WorkspaceContextLoadPlan{
		Workspace:       value.Workspace,
		CatalogRevision: value.CatalogRevision,
		Prompt:          value.Prompt,
		Diagnostics:     diagnostic.Clone(value.Diagnostics),
		Contributions:   make([]WorkspaceContextContribution, 0, len(value.Contributions)),
		Decisions:       make([]WorkspaceContextDecision, 0, len(value.Decisions)),
		PromptBytes:     value.PromptBytes,
	}
	for _, contribution := range value.Contributions {
		output.Contributions = append(
			output.Contributions,
			ContextContributionViewOf(contribution),
		)
	}
	for _, decision := range value.Decisions {
		output.Decisions = append(output.Decisions, WorkspaceContextDecision(decision))
	}
	return output
}

func ContextContributionViewOf(
	value ContextContribution,
) WorkspaceContextContribution {
	return WorkspaceContextContribution{
		Artifact:         value.Artifact,
		RecordRevision:   value.ArtifactRevision,
		DefinitionDigest: value.DefinitionDigest,
		SourceID:         value.SourceID,
		Locator:          value.Locator,
		Name:             value.Name,
		Role:             value.Role,
		MediaType:        value.MediaType,
		Content:          value.Content,
		ConventionOrder:  value.ConventionOrder,
		OriginalBytes:    value.OriginalBytes,
		IncludedBytes:    value.IncludedBytes,
		Truncated:        value.Truncated,
	}
}

func ContextViewOf(value ContextDocument) WorkspaceContextView {
	return WorkspaceContextView{
		Artifact:         value.Artifact,
		RecordRevision:   value.ArtifactRevision,
		DefinitionDigest: value.DefinitionDigest,
		SourceID:         value.SourceID,
		Locator:          value.Locator,
		Name:             value.Name,
		Role:             value.Role,
		MediaType:        value.MediaType,
		Enabled:          value.Enabled,
		State:            value.State,
		CatalogCurrent:   value.CatalogCurrent,
		ProjectionValid:  value.ProjectionValid,
		RuntimeDisabled:  value.RuntimeDisabled,
		Diagnostics:      diagnostic.Clone(value.Diagnostics),
	}
}

func WorkspaceSkillLoadViewOf(
	value workspaceadapter.SkillLoadPlan,
) WorkspaceSkillLoadView {
	output := WorkspaceSkillLoadView{
		Workspace:       value.Workspace,
		CatalogRevision: value.CatalogRevision,
		Diagnostics:     diagnostic.Clone(value.Diagnostics),
		Skills:          make([]WorkspaceSkillView, 0, len(value.Skills)),
	}
	for _, skill := range value.Skills {
		output.Skills = append(output.Skills, WorkspaceSkillViewOf(skill))
	}
	return output
}

func WorkspaceSkillViewOf(value workspaceadapter.WorkspaceSkill) WorkspaceSkillView {
	summary := WorkspaceSkillSummary{
		SchemaVersion: value.Skill.SchemaVersion,
		ID:            value.Skill.ID,
		Slug:          value.Skill.Slug,
		Name:          value.Skill.Name,
		DisplayName:   value.Skill.DisplayName,
		Description:   value.Skill.Description,
		Tags:          append([]string(nil), value.Skill.Tags...),
		Insert:        document.SkillInsert(value.Skill.Insert),
		IsEnabled:     value.Skill.IsEnabled,
		CreatedAt:     value.Skill.CreatedAt,
		ModifiedAt:    value.Skill.ModifiedAt,
		Arguments:     make([]WorkspaceSkillArgument, 0, len(value.Skill.Arguments)),
	}
	for _, argument := range value.Skill.Arguments {
		summary.Arguments = append(summary.Arguments, WorkspaceSkillArgument{
			Name:        argument.Name,
			Description: argument.Description,
			Default:     argument.Default,
		})
	}
	return WorkspaceSkillView{
		Workspace:        value.Workspace,
		Artifact:         value.Artifact,
		DefinitionDigest: value.DefinitionDigest,
		SourceID:         value.SourceID,
		Locator:          value.Locator,
		Skill:            summary,
		MarkdownBody:     value.MarkdownBody,
		RecordRevision:   value.ArtifactRevision,
		State:            value.State,
		ProjectionValid:  value.ProjectionValid,
		CatalogCurrent:   value.CatalogCurrent,
		RuntimeDisabled:  value.RuntimeDisabled,
		Diagnostics:      diagnostic.Clone(value.Diagnostics),
	}
}

func workspaceArtifactViewOf(
	value artifact.Artifact,
) WorkspaceArtifactView {
	runtimeDisabled, dataErr := artifactadapter.ArtifactRuntimeDisabled(value)
	output := workspaceArtifactView(value, runtimeDisabled)
	if dataErr != nil {
		output.Diagnostics = diagnostic.Append(
			output.Diagnostics,
			workspaceArtifactDataDiagnostic(value),
		)
	}
	return output
}

func workspaceArtifactViewOfResource(
	value workspaceDomain.Resource,
) WorkspaceArtifactView {
	return workspaceArtifactView(
		value.Artifact,
		value.ArtifactData.RuntimeDisabled,
	)
}

func workspaceArtifactView(
	value artifact.Artifact,
	runtimeDisabled bool,
) WorkspaceArtifactView {
	var digest *cryptoutil.Digest
	if value.ResolvedDefinition != nil {
		copyValue := *value.ResolvedDefinition
		digest = &copyValue
	}
	return WorkspaceArtifactView{
		Artifact:           value.Ref(),
		Revision:           value.Revision,
		Name:               value.Name,
		Kind:               value.Kind,
		Enabled:            value.Enabled,
		State:              value.State,
		Adoption:           value.Adoption,
		ResolvedDefinition: digest,
		SourceID:           value.Binding.SourceID,
		Locator:            value.Binding.Locator,
		SubresourceLocator: value.Binding.SubresourceLocator,
		RuntimeDisabled:    runtimeDisabled,
		Diagnostics:        diagnostic.Clone(value.Diagnostics),
	}
}

func workspaceArtifactDataDiagnostic(
	value artifact.Artifact,
) diagnostic.Diagnostic {
	return diagnostic.Diagnostic{
		Severity: diagnostic.SeverityError,
		Code:     workspaceDomain.DiagnosticCodeProjectionInvalid,
		Message:  "the Workspace Artifact has invalid local runtime settings",
		Location: &diagnostic.Location{
			Locator:            value.Binding.Locator,
			SubresourceLocator: value.Binding.SubresourceLocator,
		},
	}
}

func occurrenceViewKey(
	sourceID source.SourceID,
	locator basespec.Locator,
	subresource basespec.SubresourceLocator,
	kind artifact.ArtifactKind,
) string {
	return string(sourceID) + "\x00" +
		string(locator) + "\x00" +
		string(subresource) + "\x00" +
		string(kind)
}

func cloneBool(value *bool) *bool {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}
