package consumerapi

import (
	"context"
	"fmt"
	"sort"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/definition"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
	workspaceRuntime "github.com/flexigpt/flexigpt-app/internal/workspace/runtime"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
	"github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/artifactadapter"
	workspaceDomainContext "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/context"
)

type ContextContribution struct {
	Artifact         artifact.ArtifactRef                      `json:"artifact"`
	ArtifactRevision uint64                                    `json:"artifactRevision"`
	DefinitionDigest cryptoutil.Digest                         `json:"definitionDigest"`
	SourceID         source.SourceID                           `json:"sourceID"`
	Locator          basespec.Locator                          `json:"locator"`
	Name             string                                    `json:"name"`
	Role             artifactbuiltin.WorkspaceContextRole      `json:"role"`
	MediaType        artifactbuiltin.WorkspaceContextMediaType `json:"mediaType"`
	Content          string                                    `json:"content"`
	ConventionOrder  int                                       `json:"conventionOrder"`
	OriginalBytes    int                                       `json:"originalBytes"`
	IncludedBytes    int                                       `json:"includedBytes"`
	Truncated        bool                                      `json:"truncated"`
}

type ContextLoadPlan struct {
	Workspace       collection.CollectionRef `json:"workspace"`
	CatalogRevision uint64                   `json:"catalogRevision"`
	Contributions   []ContextContribution    `json:"contributions"`
	Prompt          string                   `json:"prompt"`
	Diagnostics     []diagnostic.Diagnostic  `json:"diagnostics,omitempty"`
	Decisions       []CompositionDecision    `json:"decisions"`
	PromptBytes     int                      `json:"promptBytes"`
}

type ContextDocument struct {
	Artifact         artifact.ArtifactRef                      `json:"artifact"`
	ArtifactRevision uint64                                    `json:"artifactRevision"`
	DefinitionDigest cryptoutil.Digest                         `json:"definitionDigest"`
	SourceID         source.SourceID                           `json:"sourceID"`
	Locator          basespec.Locator                          `json:"locator"`
	Name             string                                    `json:"name"`
	Role             artifactbuiltin.WorkspaceContextRole      `json:"role"`
	MediaType        artifactbuiltin.WorkspaceContextMediaType `json:"mediaType"`
	Enabled          bool                                      `json:"enabled"`
	State            artifact.State                            `json:"state"`
	CatalogCurrent   bool                                      `json:"catalogCurrent"`
	ProjectionValid  bool                                      `json:"projectionValid"`
	RuntimeDisabled  bool                                      `json:"runtimeDisabled"`
	Diagnostics      []diagnostic.Diagnostic                   `json:"diagnostics,omitempty"`
}

type ContextInspection struct {
	Workspace       collection.CollectionRef `json:"workspace"`
	CatalogRevision uint64                   `json:"catalogRevision"`
	Contributions   []ContextContribution    `json:"contributions"`
	Diagnostics     []diagnostic.Diagnostic  `json:"diagnostics,omitempty"`
}

// ContextService is the aggregate-facing storage port for Workspace Context
// list, inspection, and composition operations.
type ContextService interface {
	List(
		ctx context.Context,
		workspace collection.CollectionRef,
	) ([]ContextDocument, error)

	Load(
		ctx context.Context,
		workspace collection.CollectionRef,
		artifactRefs []artifact.ArtifactRef,
	) (ContextInspection, error)

	Compose(
		ctx context.Context,
		workspace collection.CollectionRef,
		artifactRefs []artifact.ArtifactRef,
	) (ContextLoadPlan, error)
}

type workspaceDataSource interface {
	GetWorkspace(
		ctx context.Context,
		workspace collection.CollectionRef,
	) (workspaceDomain.Workspace, error)

	Catalog(
		ctx context.Context,
		workspace collection.CollectionRef,
	) (workspaceDomain.CatalogView, error)

	ComposeLoadPlan(
		ctx context.Context,
		workspace collection.CollectionRef,
		artifactRefs []artifact.ArtifactRef,
	) (workspaceDomain.LoadPlan, error)
}

type contextService struct {
	query             workspaceDataSource
	runtimePolicy     artifactadapter.SourceUsePolicy
	compositionPolicy workspaceRuntime.CompositionPolicy
	engine            *workspaceRuntime.Engine
}

func newContextService(
	query workspaceDataSource,
	runtimePolicy artifactadapter.SourceUsePolicy,
	compositionPolicy workspaceRuntime.CompositionPolicy,
) (ContextService, error) {
	if query == nil || runtimePolicy == nil {
		return nil, fmt.Errorf(
			"%w: Workspace context adapter query is nil",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}

	compositionPolicy = compositionPolicy.Normalized()
	if err := compositionPolicy.Validate(); err != nil {
		return nil, fmt.Errorf(
			"%w: Workspace context composition policy: %w",
			workspaceDomain.ErrInvalidWorkspace,
			err,
		)
	}
	return &contextService{
		query:             query,
		runtimePolicy:     runtimePolicy,
		compositionPolicy: compositionPolicy,
		engine:            workspaceRuntime.NewEngine(),
	}, nil
}

func (p *contextService) Compose(
	ctx context.Context,
	workspace collection.CollectionRef,
	artifactRefs []artifact.ArtifactRef,
) (ContextLoadPlan, error) {
	if err := workspace.Validate(); err != nil {
		return ContextLoadPlan{}, err
	}
	if len(artifactRefs) == 0 {
		values, err := p.List(ctx, workspace)
		if err != nil {
			return ContextLoadPlan{}, err
		}
		for _, value := range values {
			if value.Enabled && value.State == artifact.StateAvailable {
				artifactRefs = append(artifactRefs, value.Artifact)
			}
		}
	}

	loadPlan, err := p.query.ComposeLoadPlan(ctx, workspace, artifactRefs)
	if err != nil {
		return ContextLoadPlan{}, err
	}

	workspaceValue, err := p.query.GetWorkspace(ctx, workspace)
	if err != nil {
		return ContextLoadPlan{}, err
	}
	output := ContextLoadPlan{
		Workspace:       workspace,
		CatalogRevision: loadPlan.CatalogRevision,
		Diagnostics:     diagnostic.Clone(loadPlan.Diagnostics),
	}
	handled := make(map[artifact.ArtifactID]struct{}, len(loadPlan.Items))
	for _, item := range loadPlan.Items {
		handled[item.Artifact.ID] = struct{}{}
		if err := workspaceDomainContext.ValidateContextDefinition(item.Definition); err != nil {
			output.Diagnostics = diagnostic.Append(
				output.Diagnostics,
				contextProjectionDiagnostic(item.Artifact, err),
			)
			output.Decisions = append(output.Decisions, CompositionDecision{
				Artifact: item.Artifact.Ref(),
				Status:   workspaceRuntime.CompositionUnavailable,
				Code:     workspaceDomain.DiagnosticCodeProjectionInvalid,
			})
			continue
		}
		decision := p.runtimePolicy.Decide(ctx, artifactadapter.RuntimePolicyRequest{
			Use:              workspaceRuntime.RuntimeUseContextPrompt,
			Workspace:        workspaceValue,
			Artifact:         item.Artifact,
			DefinitionDigest: item.Definition.Digest,
			SourceID:         item.Source.ID,
		})
		if err := decision.Validate(); err != nil {
			return ContextLoadPlan{}, err
		}
		if decision.Disposition != workspaceRuntime.RuntimeAllowed {
			output.Diagnostics = diagnostic.Append(
				output.Diagnostics,
				artifactadapter.RuntimeDecisionDiagnostic(decision, item.Artifact),
			)
			status := workspaceRuntime.CompositionDenied
			if decision.Disposition == workspaceRuntime.RuntimeUnavailable {
				status = workspaceRuntime.CompositionUnavailable
			}
			output.Decisions = append(output.Decisions, CompositionDecision{
				Artifact: item.Artifact.Ref(),
				Status:   status,
				Code:     decision.Code,
			})
			continue
		}
		body, err := definition.DecodeBody[workspaceDomainContext.Definition](
			item.Definition.Body,
		)
		if err != nil {
			handled[item.Artifact.ID] = struct{}{}
			output.Diagnostics = diagnostic.Append(
				output.Diagnostics,
				contextProjectionDiagnostic(item.Artifact, err),
			)
			output.Decisions = append(output.Decisions, CompositionDecision{
				Artifact: item.Artifact.Ref(),
				Status:   workspaceRuntime.CompositionUnavailable,
				Code:     workspaceDomain.DiagnosticCodeProjectionInvalid,
			})
			continue
		}
		output.Contributions = append(
			output.Contributions,
			ContextContribution{
				ConventionOrder: contextRuntimeOrder(
					item.Artifact.Binding.Locator,
				),
				ArtifactRevision: item.Artifact.Revision,
				Artifact:         item.Artifact.Ref(),
				DefinitionDigest: item.Definition.Digest,
				SourceID:         item.Source.ID,
				Locator:          item.Artifact.Binding.Locator,
				Name:             body.Name,
				Role:             body.Role,
				MediaType:        body.MediaType,
				Content:          body.Content,
			},
		)
	}
	for _, ref := range artifactRefs {
		if _, found := handled[ref.ArtifactID]; found {
			continue
		}
		output.Decisions = append(output.Decisions, CompositionDecision{
			Artifact: ref,
			Status:   workspaceRuntime.CompositionUnavailable,
			Code:     workspaceDomain.DiagnosticCodeArtifactUnresolved,
		})
	}

	sortContextContributions(output.Contributions)
	contributions,
		prompt,
		diagnostics,
		decisions,
		compositionErr := applyCompositionPolicy(
		p.engine,
		p.compositionPolicy,
		output.Contributions,
		output.Diagnostics,
		output.Decisions,
	)
	if compositionErr != nil {
		return ContextLoadPlan{}, fmt.Errorf(
			"%w: compose Workspace Context: %w",
			workspaceDomain.ErrInvalidWorkspace,
			compositionErr,
		)
	}
	output.Contributions = contributions
	output.Prompt = prompt
	output.Diagnostics = diagnostics
	output.Decisions = decisions
	output.PromptBytes = len(output.Prompt)
	return output, nil
}

func (p *contextService) List(
	ctx context.Context,
	workspace collection.CollectionRef,
) ([]ContextDocument, error) {
	view, err := p.query.Catalog(ctx, workspace)
	if err != nil {
		return nil, err
	}
	output := make([]ContextDocument, 0)
	for _, resourceValue := range view.Resources {
		if resourceValue.Definition.Kind != artifactbuiltin.WorkspaceContextArtifactKind ||
			resourceValue.Definition.SchemaID != artifactbuiltin.WorkspaceContextSchemaID {
			continue
		}
		value, err := projectContextDocument(resourceValue)
		if err != nil {
			value.Diagnostics = diagnostic.Append(
				value.Diagnostics,
				contextProjectionDiagnostic(resourceValue.Artifact, err),
			)
		}
		output = append(output, value)
	}
	sort.Slice(output, func(left, right int) bool {
		leftOrder := contextRuntimeOrder(output[left].Locator)
		rightOrder := contextRuntimeOrder(output[right].Locator)
		if leftOrder != rightOrder {
			return leftOrder < rightOrder
		}
		return output[left].Artifact.ArtifactID < output[right].Artifact.ArtifactID
	})
	return output, nil
}

func (p *contextService) Load(
	ctx context.Context,
	workspace collection.CollectionRef,
	artifactRefs []artifact.ArtifactRef,
) (ContextInspection, error) {
	view, err := p.query.Catalog(ctx, workspace)
	if err != nil {
		return ContextInspection{}, err
	}
	requested := make(
		map[artifact.ArtifactID]struct{},
		len(artifactRefs),
	)
	for _, ref := range artifactRefs {
		if err := ref.Validate(); err != nil {
			return ContextInspection{}, err
		}
		if ref.RootID != workspace.RootID {
			return ContextInspection{}, fmt.Errorf(
				"%w: Context Artifact belongs to another Root",
				workspaceDomain.ErrInvalidWorkspace,
			)
		}
		if _, duplicate := requested[ref.ArtifactID]; duplicate {
			return ContextInspection{}, fmt.Errorf(
				"%w: duplicate Context Artifact %q",
				workspaceDomain.ErrInvalidWorkspace,
				ref.ArtifactID,
			)
		}
		requested[ref.ArtifactID] = struct{}{}
	}
	output := ContextInspection{
		Workspace:       workspace,
		CatalogRevision: view.Catalog.Revision,
	}
	for _, resourceValue := range view.Resources {
		if resourceValue.Definition.Kind != artifactbuiltin.WorkspaceContextArtifactKind ||
			resourceValue.Definition.SchemaID != artifactbuiltin.WorkspaceContextSchemaID {
			continue
		}
		if len(requested) != 0 {
			if _, selected := requested[resourceValue.Artifact.ID]; !selected {
				continue
			}
		}
		contribution, err := projectContext(resourceValue)
		if err != nil {
			output.Diagnostics = diagnostic.Append(
				output.Diagnostics,
				contextProjectionDiagnostic(resourceValue.Artifact, err),
			)
			continue
		}
		output.Contributions = append(output.Contributions, contribution)
		output.Diagnostics = diagnostic.Append(
			output.Diagnostics,
			resourceValue.Artifact.Diagnostics...,
		)
	}
	sortContextContributions(output.Contributions)
	if len(requested) != 0 &&
		len(output.Contributions) != len(requested) {
		output.Diagnostics = diagnostic.Append(
			output.Diagnostics,
			diagnostic.Diagnostic{
				Severity: diagnostic.SeverityError,
				Code:     workspaceDomain.DiagnosticCodeArtifactUnresolved,
				Message:  "one or more requested Context Artifacts were not available for inspection",
			},
		)
	}
	return output, nil
}

func projectContextDocument(
	value workspaceDomain.Resource,
) (ContextDocument, error) {
	runtimeDisabled, dataErr := artifactadapter.ArtifactRuntimeDisabled(value.Artifact)
	output := ContextDocument{
		Artifact:         value.Artifact.Ref(),
		ArtifactRevision: value.Artifact.Revision,
		DefinitionDigest: value.Definition.Digest,
		SourceID:         value.Source.ID,
		Locator:          value.Artifact.Binding.Locator,
		Name:             value.Artifact.Name,
		Enabled:          value.Artifact.Enabled,
		State:            value.Artifact.State,
		CatalogCurrent:   value.CatalogCurrent,
		RuntimeDisabled:  runtimeDisabled,
		Diagnostics: diagnostic.Append(
			value.Artifact.Diagnostics,
			value.Diagnostics...,
		),
	}
	if dataErr != nil {
		return output, dataErr
	}
	if err := workspaceDomainContext.ValidateContextDefinition(value.Definition); err != nil {
		return output, err
	}
	body, err := definition.DecodeBody[workspaceDomainContext.Definition](
		value.Definition.Body,
	)
	if err != nil {
		return output, err
	}
	output.Name = body.Name
	output.Role = body.Role
	output.MediaType = body.MediaType
	output.ProjectionValid = true
	return output, nil
}

func projectContext(
	value workspaceDomain.Resource,
) (ContextContribution, error) {
	if err := workspaceDomainContext.ValidateContextDefinition(value.Definition); err != nil {
		return ContextContribution{}, err
	}
	body, err := definition.DecodeBody[workspaceDomainContext.Definition](value.Definition.Body)
	if err != nil {
		return ContextContribution{}, err
	}
	return ContextContribution{
		Artifact:         value.Artifact.Ref(),
		ArtifactRevision: value.Artifact.Revision,
		DefinitionDigest: value.Definition.Digest,
		SourceID:         value.Source.ID,
		Locator:          value.Artifact.Binding.Locator,
		ConventionOrder:  contextRuntimeOrder(value.Artifact.Binding.Locator),
		Name:             body.Name,
		Role:             body.Role,
		MediaType:        body.MediaType,
		Content:          body.Content,
	}, nil
}

func sortContextContributions(values []ContextContribution) {
	sort.Slice(values, func(left, right int) bool {
		if values[left].ConventionOrder != values[right].ConventionOrder {
			return values[left].ConventionOrder <
				values[right].ConventionOrder
		}
		if values[left].SourceID != values[right].SourceID {
			return values[left].SourceID < values[right].SourceID
		}
		return values[left].Locator < values[right].Locator
	})
}

func contextRuntimeOrder(locator basespec.Locator) int {
	return workspaceDomainContext.RuntimeOrder(locator)
}

func contextProjectionDiagnostic(
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
