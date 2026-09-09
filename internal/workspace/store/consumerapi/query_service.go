package consumerapi

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/resource"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
	"github.com/flexigpt/flexigpt-app/internal/workspace/store/domain/artifactadapter"
)

type QueryService struct {
	workspaces *Service
	resources  compositionapi.ResourceAPI
}

func NewQueryService(
	workspaces *Service,
	resources compositionapi.ResourceAPI,
) (*QueryService, error) {
	if workspaces == nil ||
		resources == nil {
		return nil, fmt.Errorf(
			"%w: Workspace query dependencies are incomplete",
			workspaceDomain.ErrInvalidWorkspace,
		)
	}
	return &QueryService{
		workspaces: workspaces,
		resources:  resources,
	}, nil
}

func (q *QueryService) GetWorkspace(
	ctx context.Context,
	workspace collection.CollectionRef,
) (workspaceDomain.Workspace, error) {
	return q.workspaces.Get(ctx, workspace)
}

func (q *QueryService) ResolveArtifact(
	ctx context.Context,
	ref artifact.ArtifactRef,
) (workspaceDomain.Workspace, workspaceDomain.Resource, error) {
	if err := ref.Validate(); err != nil {
		return workspaceDomain.Workspace{}, workspaceDomain.Resource{}, err
	}
	resolved, err := q.resources.ResolveArtifact(
		ctx,
		ref,
		resource.ResolveOptions{},
	)
	if err != nil {
		return workspaceDomain.Workspace{}, workspaceDomain.Resource{}, err
	}
	workspace, err := q.workspaces.Get(ctx, resolved.Collection.Ref())
	if err != nil {
		return workspaceDomain.Workspace{}, workspaceDomain.Resource{}, err
	}
	return workspace, workspaceResourceOfResolved(resolved), nil
}

func (q *QueryService) ComposeLoadPlan(
	ctx context.Context,
	workspace collection.CollectionRef,
	artifactRefs []artifact.ArtifactRef,
) (workspaceDomain.LoadPlan, error) {
	view, err := q.Catalog(ctx, workspace)
	if err != nil {
		return workspaceDomain.LoadPlan{}, err
	}
	return q.ComposeLoadPlanFromCatalog(view, artifactRefs)
}

// ComposeLoadPlanFromCatalog performs Workspace selection validation and
// runtime-read eligibility projection from one already loaded Catalog view.
//
// Context and Skill adapters use this to avoid loading a second Catalog after
// they have already selected Artifact references from the first one.
func (q *QueryService) ComposeLoadPlanFromCatalog(
	view workspaceDomain.CatalogView,
	artifactRefs []artifact.ArtifactRef,
) (workspaceDomain.LoadPlan, error) {
	workspace := view.Workspace.Collection.Ref()
	requested, err := workspaceArtifactSelection(workspace, artifactRefs)
	if err != nil {
		return workspaceDomain.LoadPlan{}, err
	}

	plan := workspaceDomain.LoadPlan{
		Workspace:       workspace,
		WorkspaceState:  view.Workspace,
		CatalogRevision: view.Catalog.Revision,
		Diagnostics: diagnostic.Append(
			view.Catalog.Diagnostics,
			view.FreshnessDiagnostics...,
		),
	}
	resources := make(map[artifact.ArtifactID]workspaceDomain.Resource, len(view.Resources))
	for _, value := range view.Resources {
		resources[value.Artifact.ID] = value
	}
	unresolved := make(
		map[artifact.ArtifactID]artifact.Artifact,
		len(view.UnresolvedArtifacts),
	)
	for _, value := range view.UnresolvedArtifacts {
		unresolved[value.ID] = value
	}

	ordered := make([]artifact.ArtifactID, 0, len(requested))
	for artifactID := range requested {
		ordered = append(ordered, artifactID)
	}
	slices.Sort(ordered)

	for _, artifactID := range ordered {
		resourceValue, found := resources[artifactID]
		if !found {
			if unresolvedValue, exists := unresolved[artifactID]; exists {
				plan.Diagnostics = diagnostic.Append(
					plan.Diagnostics,
					unresolvedValue.Diagnostics...,
				)
				plan.Diagnostics = diagnostic.Append(
					plan.Diagnostics,
					recordAvailabilityDiagnostic(
						unresolvedValue,
						workspaceDomain.DiagnosticCodeArtifactUnresolved,
						"the Workspace Artifact is unavailable for loading",
					),
				)
			} else {
				plan.Diagnostics = diagnostic.Append(
					plan.Diagnostics,
					diagnostic.Diagnostic{
						Severity: diagnostic.SeverityError,
						Code:     workspaceDomain.DiagnosticCodeArtifactUnresolved,
						Message:  "the requested Workspace Artifact was not found",
					},
				)
			}
			continue
		}

		switch {
		case !view.CatalogCurrent:
			plan.Diagnostics = diagnostic.Append(
				plan.Diagnostics,
				recordAvailabilityDiagnostic(
					resourceValue.Artifact,
					workspaceDomain.DiagnosticCodeArtifactUnavailable,
					"the Workspace catalog is stale and must be refreshed",
				),
			)
			continue

		case !resourceValue.Artifact.Enabled:
			plan.Diagnostics = diagnostic.Append(
				plan.Diagnostics,
				recordAvailabilityDiagnostic(
					resourceValue.Artifact,
					workspaceDomain.DiagnosticCodeArtifactUnavailable,
					"the Workspace Artifact is disabled",
				),
			)
			continue

		case resourceValue.Artifact.State != artifact.StateAvailable:
			plan.Diagnostics = diagnostic.Append(
				plan.Diagnostics,
				recordAvailabilityDiagnostic(
					resourceValue.Artifact,
					workspaceDomain.DiagnosticCodeArtifactUnavailable,
					"the Workspace Artifact is not available",
				),
			)
			continue

		case !resourceValue.CatalogCurrent:
			plan.Diagnostics = diagnostic.Append(
				plan.Diagnostics,
				recordAvailabilityDiagnostic(
					resourceValue.Artifact,
					workspaceDomain.DiagnosticCodeArtifactUnavailable,
					"the linked Workspace Artifact is not catalog-current",
				),
			)
			continue

		case !resourceValue.ProjectionValid:
			plan.Diagnostics = diagnostic.Append(
				plan.Diagnostics,
				resourceValue.Diagnostics...,
			)
			continue

		case resourceValue.Resolved == nil:
			plan.Diagnostics = diagnostic.Append(
				plan.Diagnostics,
				recordAvailabilityDiagnostic(
					resourceValue.Artifact,
					workspaceDomain.DiagnosticCodeArtifactUnavailable,
					"the Workspace Artifact has no current resolved resource chain",
				),
			)
			continue
		}

		resolved := resourceValue.Resolved.Clone()
		plan.Items = append(plan.Items, workspaceDomain.LoadPlanItem{
			Resolved:          resolved,
			ArtifactData:      resourceValue.ArtifactData,
			ArtifactDataValid: resourceValue.ArtifactDataValid,
			ProjectionValid:   resourceValue.ProjectionValid,
		})
		plan.Diagnostics = diagnostic.Append(
			plan.Diagnostics,
			resourceValue.Artifact.Diagnostics...,
		)
	}
	sort.Slice(plan.Items, func(left, right int) bool {
		return plan.Items[left].Resolved.Artifact.ID < plan.Items[right].Resolved.Artifact.ID
	})
	return plan, nil
}

func (q *QueryService) Catalog(
	ctx context.Context,
	workspace collection.CollectionRef,
) (workspaceDomain.CatalogView, error) {
	workspaceValue, err := q.workspaces.Get(ctx, workspace)
	if err != nil {
		return workspaceDomain.CatalogView{}, err
	}
	inspection, err := q.resources.InspectCollectionResources(
		ctx,
		workspace,
	)
	if err != nil {
		return workspaceDomain.CatalogView{}, err
	}
	snapshot := inspection.Catalog.Catalog
	catalogCurrent := inspection.Catalog.IsCurrent()

	freshnessDiagnostics := make([]diagnostic.Diagnostic, 0)
	if inspection.Catalog.MetadataChanged {
		freshnessDiagnostics = diagnostic.Append(
			freshnessDiagnostics,
			diagnostic.Diagnostic{
				Severity: diagnostic.SeverityWarning,
				Code:     workspaceDomain.DiagnosticCodeCatalogStale,
				Message:  "the Workspace catalog no longer matches current collection metadata",
			},
		)
	}
	if inspection.Catalog.DecoderChanged {
		freshnessDiagnostics = diagnostic.Append(
			freshnessDiagnostics,
			diagnostic.Diagnostic{
				Severity: diagnostic.SeverityWarning,
				Code:     workspaceDomain.DiagnosticCodeCatalogDecoderStale,
				Message:  "the Workspace decoder capability set changed after this catalog was published",
			},
		)
	}
	if inspection.Catalog.PlanChanged {
		freshnessDiagnostics = diagnostic.Append(
			freshnessDiagnostics,
			diagnostic.Diagnostic{
				Severity: diagnostic.SeverityWarning,
				Code:     workspaceDomain.DiagnosticCodeCatalogPlanStale,
				Message:  "the Workspace provider discovery behavior changed after this catalog was published",
			},
		)
	}
	view := workspaceDomain.CatalogView{
		Workspace:            workspaceValue,
		Catalog:              snapshot.Clone(),
		CatalogCurrent:       catalogCurrent,
		FreshnessDiagnostics: freshnessDiagnostics,
	}
	for _, value := range inspection.Resources {
		view.Resources = append(
			view.Resources,
			workspaceResourceOfCollectionArtifact(value),
		)
	}
	for _, value := range inspection.UnresolvedArtifacts {
		view.UnresolvedArtifacts = append(
			view.UnresolvedArtifacts,
			workspaceUnresolvedArtifactOf(value),
		)
	}
	for _, value := range inspection.UnrecordedOccurrences {
		view.Unrecorded = append(view.Unrecorded, value.Clone())
	}
	sort.Slice(view.Resources, func(left, right int) bool {
		if view.Resources[left].Artifact.Kind !=
			view.Resources[right].Artifact.Kind {
			return view.Resources[left].Artifact.Kind <
				view.Resources[right].Artifact.Kind
		}
		if view.Resources[left].Artifact.Name !=
			view.Resources[right].Artifact.Name {
			return view.Resources[left].Artifact.Name <
				view.Resources[right].Artifact.Name
		}
		return view.Resources[left].Artifact.ID <
			view.Resources[right].Artifact.ID
	})
	view.Groups = groupCatalogResources(view.Resources, view.Unrecorded)
	return view, nil
}

func workspaceArtifactSelection(
	workspace collection.CollectionRef,
	artifactRefs []artifact.ArtifactRef,
) (map[artifact.ArtifactID]struct{}, error) {
	requested := make(
		map[artifact.ArtifactID]struct{},
		len(artifactRefs),
	)
	for _, ref := range artifactRefs {
		if err := ref.Validate(); err != nil {
			return nil, err
		}
		if ref.RootID != workspace.RootID {
			return nil, fmt.Errorf(
				"%w: Workspace Artifact belongs to another Root",
				workspaceDomain.ErrReferenceUnresolved,
			)
		}
		if _, duplicate := requested[ref.ArtifactID]; duplicate {
			return nil, fmt.Errorf(
				"%w: duplicate Workspace Artifact %q",
				workspaceDomain.ErrInvalidWorkspace,
				ref.ArtifactID,
			)
		}
		requested[ref.ArtifactID] = struct{}{}
	}
	return requested, nil
}

func workspaceResourceOfResolved(
	value resource.ResolvedArtifact,
) workspaceDomain.Resource {
	resolved := value.Clone()
	return workspaceResourceOfCollectionArtifact(
		resource.CollectionArtifactResource{
			Artifact:       resolved.Artifact,
			Definition:     resolved.Definition,
			Occurrence:     resolved.Occurrence,
			Source:         resolved.Source,
			CatalogCurrent: true,
			Resolved:       &resolved,
		},
	)
}

func workspaceResourceOfCollectionArtifact(
	value resource.CollectionArtifactResource,
) workspaceDomain.Resource {
	artifactData, artifactDataValid, diagnostics := workspaceArtifactDataOf(value.Artifact)
	projectionValid := artifactDataValid

	if value.Definition.Kind != value.Artifact.Kind {
		projectionValid = false
		diagnostics = append(
			diagnostics,
			projectionDiagnostic(
				value.Artifact,
				fmt.Errorf(
					"artifact kind %q does not match resolved definition kind %q",
					value.Artifact.Kind,
					value.Definition.Kind,
				),
			),
		)
	}

	occurrence := value.Occurrence.Clone()
	output := workspaceDomain.Resource{
		Artifact:          value.Artifact.Clone(),
		ArtifactData:      artifactData,
		ArtifactDataValid: artifactDataValid,
		Definition:        value.Definition.Clone(),
		Occurrence:        &occurrence,
		Source:            value.Source.Clone(),
		CatalogCurrent:    value.CatalogCurrent,
		ProjectionValid:   projectionValid,
		Diagnostics:       diagnostics,
	}
	if value.Resolved != nil {
		resolved := value.Resolved.Clone()
		output.Resolved = &resolved
	}
	return output
}

func workspaceArtifactDataOf(
	value artifact.Artifact,
) (
	workspaceDomain.ArtifactData,
	bool,
	[]diagnostic.Diagnostic,
) {
	data, err := artifactadapter.DecodeArtifactData(value.Data)
	if err != nil {
		return workspaceDomain.ArtifactData{},
			false,
			[]diagnostic.Diagnostic{
				projectionDiagnostic(value, err),
			}
	}
	return data, true, nil
}

func workspaceUnresolvedArtifactOf(
	value resource.ArtifactResolutionIssue,
) artifact.Artifact {
	switch value.Status {
	case resource.ArtifactResolutionSourceUnavailable:
		return recordWithDiagnostic(
			value.Artifact,
			recordSourceUnavailableDiagnostic(value.Artifact),
		)

	case resource.ArtifactResolutionOccurrenceUnavailable:
		return recordWithDiagnostic(
			value.Artifact,
			recordDefinitionUnavailableDiagnostic(
				value.Artifact,
				errors.New("current catalog has no matching occurrence"),
			),
		)

	case resource.ArtifactResolutionDefinitionUnavailable:
		return recordWithDiagnostic(
			value.Artifact,
			recordDefinitionUnavailableDiagnostic(
				value.Artifact,
				errors.New("current catalog definition is unavailable"),
			),
		)

	case resource.ArtifactResolutionDefinitionMismatch:
		return recordWithDiagnostic(
			value.Artifact,
			recordDefinitionUnavailableDiagnostic(
				value.Artifact,
				errors.New("catalog definition fingerprint differs from artifact state"),
			),
		)

	default:
		return value.Artifact.Clone()
	}
}

func recordAvailabilityDiagnostic(
	value artifact.Artifact,
	code string,
	message string,
) diagnostic.Diagnostic {
	return diagnostic.Diagnostic{
		Severity: diagnostic.SeverityError,
		Code:     code,
		Message:  message,
		Location: &diagnostic.Location{
			Locator:            value.Binding.Locator,
			SubresourceLocator: value.Binding.SubresourceLocator,
		},
	}
}

func projectionDiagnostic(
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

func recordWithDiagnostic(
	value artifact.Artifact,
	d diagnostic.Diagnostic,
) artifact.Artifact {
	output := value
	output.Diagnostics = diagnostic.Append(
		[]diagnostic.Diagnostic{d},
		value.Diagnostics...,
	)
	return output
}

func recordSourceUnavailableDiagnostic(
	value artifact.Artifact,
) diagnostic.Diagnostic {
	return diagnostic.Diagnostic{
		Severity: diagnostic.SeverityError,
		Code:     workspaceDomain.DiagnosticCodeArtifactUnavailable,
		Message:  "the Artifact Source is no longer attached to this Workspace",
		Location: &diagnostic.Location{
			Locator:            value.Binding.Locator,
			SubresourceLocator: value.Binding.SubresourceLocator,
		},
	}
}

func recordDefinitionUnavailableDiagnostic(
	value artifact.Artifact,
	cause error,
) diagnostic.Diagnostic {
	return diagnostic.Diagnostic{
		Severity: diagnostic.SeverityError,
		Code:     workspaceDomain.DiagnosticCodeArtifactUnavailable,
		Message: diagnostic.BoundedMessage(
			fmt.Sprintf(
				"the resolved Workspace Artifact definition could not be read: %v",
				cause,
			),
		),
		Location: &diagnostic.Location{
			Locator:            value.Binding.Locator,
			SubresourceLocator: value.Binding.SubresourceLocator,
		},
	}
}

func groupCatalogResources(
	resources []workspaceDomain.Resource,
	unrecorded []catalog.Occurrence,
) []workspaceDomain.ResourceGroup {
	values := make(map[artifact.ArtifactKind]*workspaceDomain.ResourceGroup)
	for _, resourceValue := range resources {
		kind := resourceValue.Artifact.Kind
		group := values[kind]
		if group == nil {
			group = &workspaceDomain.ResourceGroup{Kind: kind}
			values[kind] = group
		}
		group.Resources = append(group.Resources, resourceValue)
	}
	for _, occurrence := range unrecorded {
		if occurrence.Kind == "" {
			continue
		}
		group := values[occurrence.Kind]
		if group == nil {
			group = &workspaceDomain.ResourceGroup{Kind: occurrence.Kind}
			values[occurrence.Kind] = group
		}
		group.Unrecorded = append(group.Unrecorded, occurrence)
	}

	output := make([]workspaceDomain.ResourceGroup, 0, len(values))
	for _, group := range values {
		output = append(output, *group)
	}
	sort.Slice(output, func(left, right int) bool {
		return output[left].Kind < output[right].Kind
	})
	return output
}
