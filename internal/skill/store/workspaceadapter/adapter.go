package workspaceadapter

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/flexigpt/agentskills-go/document"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/resource"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
	skillDomain "github.com/flexigpt/flexigpt-app/internal/skill/store/domain"
	workspaceRuntime "github.com/flexigpt/flexigpt-app/internal/workspace/runtime"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
	"github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/artifactadapter"
)

type SkillArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Default     string `json:"default,omitempty"`
}

type SkillSummary struct {
	SchemaVersion string              `json:"schemaVersion"`
	ID            artifact.ArtifactID `json:"id"`
	Slug          string              `json:"slug"`
	Name          string              `json:"name"`
	DisplayName   string              `json:"displayName"`
	Description   string              `json:"description"`
	Tags          []string            `json:"tags,omitempty"`
	Insert        string              `json:"insert"`
	Arguments     []SkillArgument     `json:"arguments,omitempty"`
	IsEnabled     bool                `json:"isEnabled"`
	CreatedAt     time.Time           `json:"createdAt"`
	ModifiedAt    time.Time           `json:"modifiedAt"`
}

type WorkspaceSkill struct {
	Workspace        collection.CollectionRef `json:"workspace"`
	Artifact         artifact.ArtifactRef     `json:"artifact"`
	DefinitionDigest cryptoutil.Digest        `json:"definitionDigest"`
	SourceID         source.SourceID          `json:"sourceID"`
	Locator          basespec.Locator         `json:"locator"`
	Skill            SkillSummary             `json:"skill"`
	MarkdownBody     string                   `json:"markdownBody,omitempty"`
	ArtifactRevision uint64                   `json:"artifactRevision"`
	State            artifact.State           `json:"state"`
	CatalogCurrent   bool                     `json:"catalogCurrent"`
	WorkspaceEnabled bool                     `json:"-"`
	RuntimeDisabled  bool                     `json:"runtimeDisabled"`
	Diagnostics      []diagnostic.Diagnostic  `json:"diagnostics,omitempty"`

	ProjectionValid     bool              `json:"-"`
	RuntimePathBacked   bool              `json:"-"`
	SourceContentDigest cryptoutil.Digest `json:"-"`
	SourceGeneration    string            `json:"-"`
	RuntimeLocation     string            `json:"-"`
}

type SkillLoadPlan struct {
	Workspace       collection.CollectionRef `json:"workspace"`
	CatalogRevision uint64                   `json:"catalogRevision"`
	Skills          []WorkspaceSkill         `json:"skills"`
	Diagnostics     []diagnostic.Diagnostic  `json:"diagnostics,omitempty"`
}

type WorkspaceDataSource interface {
	Catalog(
		ctx context.Context,
		workspace collection.CollectionRef,
	) (workspaceDomain.CatalogView, error)

	ResolveArtifact(
		ctx context.Context,
		ref artifact.ArtifactRef,
	) (workspaceDomain.Workspace, workspaceDomain.Resource, error)

	ComposeLoadPlanFromCatalog(
		view workspaceDomain.CatalogView,
		artifactRefs []artifact.ArtifactRef,
	) (workspaceDomain.LoadPlan, error)
}

type ArtifactResourceReader interface {
	ResolveVerifiedLocalPath(
		ctx context.Context,
		resolved resource.ResolvedArtifact,
		localLocator basespec.Locator,
	) (string, error)

	SupportsLocalPath(
		kind source.SourceKind,
	) bool
}

type Adapter struct {
	query         WorkspaceDataSource
	runtimePolicy artifactadapter.SourceUsePolicy
	resourceAPI   ArtifactResourceReader
}

func NewAdapter(
	query WorkspaceDataSource,
	runtimePolicy artifactadapter.SourceUsePolicy,
	resourceAPI ArtifactResourceReader,
) (*Adapter, error) {
	if query == nil || runtimePolicy == nil || resourceAPI == nil {
		return nil, fmt.Errorf(
			"%w: Workspace Skill adapter dependencies are incomplete",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}

	return &Adapter{
		query:         query,
		runtimePolicy: runtimePolicy,
		resourceAPI:   resourceAPI,
	}, nil
}

func (f *Adapter) List(
	ctx context.Context,
	workspace collection.CollectionRef,
) ([]WorkspaceSkill, error) {
	if err := workspace.Validate(); err != nil {
		return nil, err
	}
	view, err := f.query.Catalog(ctx, workspace)
	if err != nil {
		return nil, err
	}
	output := make([]WorkspaceSkill, 0)
	for _, resourceValue := range view.Resources {
		if resourceValue.Definition.Kind != artifactbuiltin.AgentSkillArtifactKind ||
			resourceValue.Definition.SchemaID != artifactbuiltin.AgentSkillSchemaID {
			continue
		}
		value, err := projectWorkspaceSkill(
			workspace,
			resourceValue,
			view.Workspace.Collection.Enabled,
			false,
			f.supportsRuntimePath(resourceValue.Source.Kind),
		)
		if err != nil {
			value.Diagnostics = diagnostic.Append(
				value.Diagnostics,
				skillProjectionDiagnostic(resourceValue.Artifact, err),
			)
		}
		output = append(output, value)
	}
	sortWorkspaceSkills(output)
	return output, nil
}

func (f *Adapter) LoadArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
) (WorkspaceSkill, error) {
	workspaceValue, resourceValue, err := f.query.ResolveArtifact(ctx, ref)
	if err != nil {
		return WorkspaceSkill{}, err
	}
	if resourceValue.Definition.Kind != artifactbuiltin.AgentSkillArtifactKind ||
		resourceValue.Definition.SchemaID != artifactbuiltin.AgentSkillSchemaID {
		return WorkspaceSkill{}, fmt.Errorf(
			"%w: Artifact %q is not an Agent Skill",
			workspaceDomain.ErrReferenceUnresolved,
			ref.ArtifactID,
		)
	}
	value, failure, err := f.prepareWorkspaceSkill(
		ctx,
		workspaceValue,
		resourceValue,
	)
	if err != nil {
		return WorkspaceSkill{}, err
	}
	if failure != nil {
		return WorkspaceSkill{}, fmt.Errorf(
			"%w: Artifact %q is unavailable for runtime loading: %s",
			workspaceDomain.ErrReferenceUnresolved,
			ref.ArtifactID,
			failure.Message,
		)
	}
	return value, nil
}

func (f *Adapter) Load(
	ctx context.Context,
	workspace collection.CollectionRef,
	artifactRefs []artifact.ArtifactRef,
) (SkillLoadPlan, error) {
	view, err := f.query.Catalog(ctx, workspace)
	if err != nil {
		return SkillLoadPlan{}, err
	}
	return f.loadFromCatalog(ctx, view, artifactRefs)
}

// LoadAll resolves every currently runtime-eligible Workspace Skill. It is
// used by the Skill catalog bridge, which must fail closed when one selected
// runtime registration cannot be materialized.
func (f *Adapter) LoadAll(
	ctx context.Context,
	workspace collection.CollectionRef,
) (SkillLoadPlan, error) {
	view, err := f.query.Catalog(ctx, workspace)
	if err != nil {
		return SkillLoadPlan{}, err
	}

	artifactRefs := f.runtimeEligibleSkillRefs(view)
	plan, err := f.loadFromCatalog(ctx, view, artifactRefs)
	if err != nil {
		return SkillLoadPlan{}, err
	}
	if len(plan.Skills) != len(artifactRefs) {
		return SkillLoadPlan{}, fmt.Errorf(
			"%w: one or more Workspace Skills could not be projected",
			basespec.ErrCatalogStale,
		)
	}
	return plan, nil
}

func (f *Adapter) loadFromCatalog(
	ctx context.Context,
	view workspaceDomain.CatalogView,
	artifactRefs []artifact.ArtifactRef,
) (SkillLoadPlan, error) {
	loadPlan, err := f.query.ComposeLoadPlanFromCatalog(
		view,
		artifactRefs,
	)
	if err != nil {
		return SkillLoadPlan{}, err
	}
	workspaceValue := loadPlan.WorkspaceState

	output := SkillLoadPlan{
		Workspace:       loadPlan.Workspace,
		CatalogRevision: loadPlan.CatalogRevision,
		Diagnostics:     diagnostic.Clone(loadPlan.Diagnostics),
	}

	for _, item := range loadPlan.Items {
		resourceValue := workspaceResourceFromLoadPlanItem(item)
		projected, failure, err := f.prepareWorkspaceSkill(
			ctx,
			workspaceValue,
			resourceValue,
		)
		if err != nil {
			return SkillLoadPlan{}, err
		}
		if failure != nil {
			output.Diagnostics = diagnostic.Append(
				output.Diagnostics,
				*failure,
			)
			continue
		}

		output.Skills = append(output.Skills, projected)
	}
	sortWorkspaceSkills(output.Skills)
	return output, nil
}

func (f *Adapter) runtimeEligibleSkillRefs(
	view workspaceDomain.CatalogView,
) []artifact.ArtifactRef {
	if !view.Workspace.Collection.Enabled {
		return nil
	}

	output := make([]artifact.ArtifactRef, 0)
	for _, value := range view.Resources {
		if value.Definition.Kind != artifactbuiltin.AgentSkillArtifactKind ||
			value.Definition.SchemaID != artifactbuiltin.AgentSkillSchemaID ||
			!value.ProjectionValid ||
			!value.CatalogCurrent ||
			value.Resolved == nil ||
			!value.Artifact.Enabled ||
			value.Artifact.State != artifact.StateAvailable ||
			value.ArtifactData.RuntimeDisabled ||
			!f.supportsRuntimePath(value.Source.Kind) {
			continue
		}
		output = append(output, value.Artifact.Ref())
	}
	return output
}

func workspaceResourceFromLoadPlanItem(
	value workspaceDomain.LoadPlanItem,
) workspaceDomain.Resource {
	resolved := value.Resolved.Clone()
	occurrence := resolved.Occurrence.Clone()
	return workspaceDomain.Resource{
		Artifact:          resolved.Artifact.Clone(),
		ArtifactData:      value.ArtifactData,
		ArtifactDataValid: value.ArtifactDataValid,
		Definition:        resolved.Definition.Clone(),
		Occurrence:        &occurrence,
		Source:            resolved.Source.Clone(),
		CatalogCurrent:    true,
		ProjectionValid:   value.ProjectionValid,
		Resolved:          &resolved,
	}
}

func (f *Adapter) prepareWorkspaceSkill(
	ctx context.Context,
	workspaceValue workspaceDomain.Workspace,
	resourceValue workspaceDomain.Resource,
) (WorkspaceSkill, *diagnostic.Diagnostic, error) {
	workspace := workspaceValue.Collection.Ref()
	projected, err := projectWorkspaceSkill(
		workspace,
		resourceValue,
		workspaceValue.Collection.Enabled,
		true,
		f.supportsRuntimePath(resourceValue.Source.Kind),
	)
	if err != nil {
		failure := skillProjectionDiagnostic(resourceValue.Artifact, err)
		return WorkspaceSkill{}, &failure, nil
	}
	if !projected.ProjectionValid {
		failure := skillProjectionDiagnostic(
			resourceValue.Artifact,
			fmt.Errorf(
				"%w: Workspace Skill projection is invalid",
				workspaceDomain.ErrInvalidWorkspace,
			),
		)
		return WorkspaceSkill{}, &failure, nil
	}

	decision := f.runtimePolicy.Decide(ctx, artifactadapter.RuntimePolicyRequest{
		Use:                    workspaceRuntime.RuntimeUseSkill,
		Workspace:              workspaceValue,
		Artifact:               resourceValue.Artifact,
		DefinitionDigest:       resourceValue.Definition.Digest,
		SourceID:               resourceValue.Source.ID,
		RuntimeDisabled:        resourceValue.ArtifactData.RuntimeDisabled,
		RuntimeSettingsInvalid: !resourceValue.ArtifactDataValid,
	})
	if err := decision.Validate(); err != nil {
		return WorkspaceSkill{}, nil, err
	}
	if decision.Disposition != workspaceRuntime.RuntimeAllowed {
		failure := artifactadapter.RuntimeDecisionDiagnostic(
			decision,
			resourceValue.Artifact,
		)
		return WorkspaceSkill{}, &failure, nil
	}
	if resourceValue.Resolved == nil {
		failure := runtimeLocationDiagnostic(
			resourceValue.Artifact,
			fmt.Errorf(
				"%w: Workspace Skill has no current resolved resource chain",
				basespec.ErrCatalogStale,
			),
		)
		return WorkspaceSkill{}, &failure, nil
	}

	resolved := resourceValue.Resolved.Clone()
	packageLocator, err := skillDomain.RuntimePackageLocator(
		resolved.Artifact.Binding.Locator,
		resolved.Artifact.Binding.SubresourceLocator,
	)
	if err != nil {
		failure := runtimeLocationDiagnostic(resolved.Artifact, err)
		return WorkspaceSkill{}, &failure, nil
	}
	runtimeLocation, err := f.resourceAPI.ResolveVerifiedLocalPath(
		ctx,
		resolved,
		packageLocator,
	)
	if err != nil {
		failure := runtimeLocationDiagnostic(resolved.Artifact, err)
		return WorkspaceSkill{}, &failure, nil
	}
	if resolved.Occurrence.SourceContentDigest == nil {
		failure := runtimeLocationDiagnostic(
			resolved.Artifact,
			fmt.Errorf(
				"%w: resolved Workspace Skill has no source digest",
				basespec.ErrDigestMismatch,
			),
		)
		return WorkspaceSkill{}, &failure, nil
	}

	projected.SourceContentDigest = *resolved.Occurrence.SourceContentDigest
	projected.SourceGeneration = resolved.SourceGeneration
	projected.RuntimeLocation = runtimeLocation
	return projected, nil, nil
}

func projectWorkspaceSkill(
	workspace collection.CollectionRef,
	resourceValue workspaceDomain.Resource,
	workspaceEnabled bool,
	includeMarkdown bool,
	runtimePathBacked bool,
) (WorkspaceSkill, error) {
	output := WorkspaceSkill{
		Workspace:        workspace,
		Artifact:         resourceValue.Artifact.Ref(),
		ArtifactRevision: resourceValue.Artifact.Revision,
		DefinitionDigest: resourceValue.Definition.Digest,
		SourceID:         resourceValue.Source.ID,
		Locator:          resourceValue.Artifact.Binding.Locator,
		State:            resourceValue.Artifact.State,
		CatalogCurrent:   resourceValue.CatalogCurrent,
		RuntimeDisabled:  resourceValue.ArtifactData.RuntimeDisabled,
		WorkspaceEnabled: workspaceEnabled,
		Diagnostics: diagnostic.Append(
			resourceValue.Artifact.Diagnostics,
			resourceValue.Diagnostics...,
		),
		RuntimePathBacked: runtimePathBacked,
	}
	if !resourceValue.ProjectionValid {
		return output, nil
	}
	doc, err := skillDomain.DocumentFromDefinition(
		resourceValue.Definition,
	)
	if err != nil {
		return output, err
	}
	markdownBody := ""
	if includeMarkdown {
		markdownBody = doc.MarkdownBody
	}
	output.Skill = skillSummary(resourceValue.Artifact, doc)
	output.MarkdownBody = markdownBody
	if resourceValue.Occurrence != nil &&
		resourceValue.Occurrence.SourceContentDigest != nil {
		output.SourceContentDigest = *resourceValue.Occurrence.SourceContentDigest
	}
	output.ProjectionValid = true
	return output, nil
}

func skillSummary(
	artifactValue artifact.Artifact,
	value document.SkillDocument,
) SkillSummary {
	arguments := make([]SkillArgument, 0, len(value.Arguments))
	for _, argument := range value.Arguments {
		arguments = append(arguments, SkillArgument{
			Name:        argument.Name,
			Description: argument.Description,
			Default:     argument.Default,
		})
	}
	return SkillSummary{
		SchemaVersion: artifactbuiltin.AgentSkillSchemaVersion,
		ID:            artifactValue.ID,
		Slug:          value.Name,
		Name:          value.Name,
		DisplayName:   value.DisplayName,
		Description:   value.Description,
		Tags:          append([]string(nil), value.Tags...),
		Insert:        string(value.Insert),
		Arguments:     arguments,
		IsEnabled:     artifactValue.Enabled,
		CreatedAt:     artifactValue.CreatedAt,
		ModifiedAt:    artifactValue.ModifiedAt,
	}
}

func sortWorkspaceSkills(values []WorkspaceSkill) {
	sort.Slice(values, func(left, right int) bool {
		if values[left].Skill.Name != values[right].Skill.Name {
			return values[left].Skill.Name < values[right].Skill.Name
		}
		return values[left].Artifact.ArtifactID < values[right].Artifact.ArtifactID
	})
}

func (f *Adapter) supportsRuntimePath(
	kind source.SourceKind,
) bool {
	return f.resourceAPI.SupportsLocalPath(kind)
}

func runtimeLocationDiagnostic(
	value artifact.Artifact,
	err error,
) diagnostic.Diagnostic {
	return diagnostic.Diagnostic{
		Severity: diagnostic.SeverityError,
		Code:     workspaceDomain.DiagnosticCodeRuntimeUnavailable,
		Message:  diagnostic.BoundedMessage(err.Error()),
		Location: &diagnostic.Location{
			Locator:            value.Binding.Locator,
			SubresourceLocator: value.Binding.SubresourceLocator,
		},
	}
}

func skillProjectionDiagnostic(
	value artifact.Artifact,
	err error,
) diagnostic.Diagnostic {
	return diagnostic.Diagnostic{
		Severity: diagnostic.SeverityError,
		Code:     workspaceDomain.DiagnosticCodeProjectionInvalid,
		Message:  diagnostic.BoundedMessage(err.Error()),
		Location: &diagnostic.Location{
			Locator:            value.Binding.Locator,
			SubresourceLocator: value.Binding.SubresourceLocator,
		},
	}
}
