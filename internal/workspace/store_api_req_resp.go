package workspace

import (
	"fmt"
	"time"

	"github.com/flexigpt/agentskills-go/document"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/diagnostic"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
	"github.com/flexigpt/flexigpt-app/internal/workspace/contextadapter"
	"github.com/flexigpt/flexigpt-app/internal/workspace/spec"
)

type WorkspaceRef = collection.CollectionRef

type WorkspaceDiscoveryRoot struct {
	Root            basespec.Locator `json:"root"`
	Recursive       bool             `json:"recursive"`
	IncludePatterns []string         `json:"includePatterns,omitempty"`
}
type WorkspaceDiscovery struct {
	AdditionalLocators []basespec.Locator       `json:"additionalLocators,omitempty"`
	AdditionalRoots    []WorkspaceDiscoveryRoot `json:"additionalRoots,omitempty"`
	IncludeReadme      bool                     `json:"includeReadme,omitempty"`
}

type WorkspaceAttachmentSettings struct {
	Recursive     *bool `json:"recursive,omitempty"`
	Authoritative *bool `json:"authoritative,omitempty"`
}

type WorkspaceArtifactSettings struct {
	RuntimeDisabled bool `json:"runtimeDisabled"`
}

type WorkspaceAttachmentView struct {
	SourceID          source.SourceID             `json:"sourceID"`
	Revision          uint64                      `json:"revision"`
	Role              collection.AttachmentRole   `json:"role"`
	Enabled           bool                        `json:"enabled"`
	SourceDisplayName string                      `json:"sourceDisplayName,omitempty"`
	SourceKind        string                      `json:"sourceKind,omitempty"`
	Path              string                      `json:"path,omitempty"`
	Settings          WorkspaceAttachmentSettings `json:"settings"`
	Diagnostics       []diagnostic.Diagnostic     `json:"diagnostics,omitempty"`
}

// WorkspaceView is the API-safe representation of a workspace.
//
// It deliberately excludes source configuration, root data, attachment raw
// data, and the trust-reference contents. Local filesystem paths are included
// because the local Workspace management UI intentionally displays them.
type WorkspaceView struct {
	Workspace       WorkspaceRef              `json:"workspace"`
	Revision        uint64                    `json:"revision"`
	DisplayName     string                    `json:"displayName"`
	Description     string                    `json:"description,omitempty"`
	Enabled         bool                      `json:"enabled"`
	Mode            spec.Mode                 `json:"mode"`
	PrimarySourceID source.SourceID           `json:"primarySourceID,omitempty"`
	PrimaryPath     string                    `json:"primaryPath,omitempty"`
	Discovery       WorkspaceDiscovery        `json:"discovery"`
	Attachments     []WorkspaceAttachmentView `json:"attachments"`
}

type WorkspaceArtifactView struct {
	Artifact           artifact.ArtifactRef        `json:"artifact"`
	Revision           uint64                      `json:"revision"`
	Name               string                      `json:"name"`
	Kind               artifact.ArtifactKind       `json:"kind"`
	Enabled            bool                        `json:"enabled"`
	State              artifact.State              `json:"state"`
	Adoption           artifact.AdoptionMode       `json:"adoption"`
	ResolvedDefinition *cryptoutil.Digest          `json:"resolvedDefinition,omitempty"`
	SourceID           source.SourceID             `json:"sourceID"`
	Locator            basespec.Locator            `json:"locator"`
	SubresourceLocator basespec.SubresourceLocator `json:"subresourceLocator,omitempty"`
	RuntimeDisabled    bool                        `json:"runtimeDisabled"`
	Diagnostics        []diagnostic.Diagnostic     `json:"diagnostics,omitempty"`
}

type WorkspaceResourceView struct {
	Artifact         WorkspaceArtifactView   `json:"artifact"`
	DefinitionDigest cryptoutil.Digest       `json:"definitionDigest"`
	SourceID         source.SourceID         `json:"sourceID"`
	Locator          basespec.Locator        `json:"locator"`
	CatalogCurrent   bool                    `json:"catalogCurrent"`
	ProjectionValid  bool                    `json:"projectionValid"`
	Diagnostics      []diagnostic.Diagnostic `json:"diagnostics,omitempty"`
}

type WorkspaceOccurrenceView struct {
	SourceID            source.SourceID             `json:"sourceID"`
	Locator             basespec.Locator            `json:"locator"`
	SubresourceLocator  basespec.SubresourceLocator `json:"subresourceLocator,omitempty"`
	Kind                artifact.ArtifactKind       `json:"kind,omitempty"`
	LogicalName         basespec.LogicalName        `json:"logicalName,omitempty"`
	LogicalVersion      basespec.LogicalVersion     `json:"logicalVersion,omitempty"`
	DefinitionDigest    *cryptoutil.Digest          `json:"definitionDigest,omitempty"`
	SourceContentDigest *cryptoutil.Digest          `json:"sourceContentDigest,omitempty"`
	State               string                      `json:"state"`
	Recorded            bool                        `json:"recorded"`
	Artifact            *artifact.ArtifactRef       `json:"artifact,omitempty"`
	Diagnostics         []diagnostic.Diagnostic     `json:"diagnostics,omitempty"`
}

type WorkspaceResourceGroupView struct {
	Kind       artifact.ArtifactKind     `json:"kind"`
	Resources  []WorkspaceResourceView   `json:"resources"`
	Unrecorded []WorkspaceOccurrenceView `json:"unrecorded"`
}

type WorkspaceCatalogView struct {
	Workspace               WorkspaceView                `json:"workspace"`
	CatalogRevision         uint64                       `json:"catalogRevision"`
	CatalogCurrent          bool                         `json:"catalogCurrent"`
	Diagnostics             []diagnostic.Diagnostic      `json:"diagnostics,omitempty"`
	Resources               []WorkspaceResourceView      `json:"resources"`
	Groups                  []WorkspaceResourceGroupView `json:"groups"`
	Occurrences             []WorkspaceOccurrenceView    `json:"occurrences"`
	ValidOccurrences        []WorkspaceOccurrenceView    `json:"validOccurrences"`
	InvalidOccurrences      []WorkspaceOccurrenceView    `json:"invalidOccurrences"`
	MissingOccurrences      []WorkspaceOccurrenceView    `json:"missingOccurrences"`
	UnrecordedOccurrences   []WorkspaceOccurrenceView    `json:"unrecordedOccurrences"`
	UnresolvedArtifacts     []WorkspaceArtifactView      `json:"unresolvedArtifacts"`
	UnrecordedCount         int                          `json:"unrecordedCount"`
	UnresolvedArtifactCount int                          `json:"unresolvedArtifactCount"`
}

type WorkspaceContextContribution struct {
	Artifact         artifact.ArtifactRef                      `json:"artifact"`
	RecordRevision   uint64                                    `json:"recordRevision"`
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

type WorkspaceContextDecision struct {
	Artifact      artifact.ArtifactRef             `json:"artifact"`
	Status        contextadapter.CompositionStatus `json:"status"`
	Code          string                           `json:"code,omitempty"`
	OriginalBytes int                              `json:"originalBytes"`
	IncludedBytes int                              `json:"includedBytes"`
}

type WorkspaceContextLoadPlan struct {
	Workspace       WorkspaceRef                   `json:"workspace"`
	CatalogRevision uint64                         `json:"catalogRevision"`
	Contributions   []WorkspaceContextContribution `json:"contributions"`
	Prompt          string                         `json:"prompt"`
	Diagnostics     []diagnostic.Diagnostic        `json:"diagnostics,omitempty"`
	Decisions       []WorkspaceContextDecision     `json:"decisions"`
	PromptBytes     int                            `json:"promptBytes"`
}

type WorkspaceContextView struct {
	Artifact         artifact.ArtifactRef                      `json:"artifact"`
	RecordRevision   uint64                                    `json:"recordRevision"`
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

type WorkspaceContextInspectionView struct {
	Workspace       WorkspaceRef                   `json:"workspace"`
	CatalogRevision uint64                         `json:"catalogRevision"`
	Contributions   []WorkspaceContextContribution `json:"contributions"`
	Diagnostics     []diagnostic.Diagnostic        `json:"diagnostics,omitempty"`
}

type WorkspaceSkillArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Default     string `json:"default,omitempty"`
}

type WorkspaceSkillSummary struct {
	SchemaVersion string                   `json:"schemaVersion"`
	ID            artifact.ArtifactID      `json:"id"`
	Slug          string                   `json:"slug"`
	Name          string                   `json:"name"`
	DisplayName   string                   `json:"displayName"`
	Description   string                   `json:"description"`
	Tags          []string                 `json:"tags,omitempty"`
	Insert        document.SkillInsert     `json:"insert"`
	Arguments     []WorkspaceSkillArgument `json:"arguments,omitempty"`
	IsEnabled     bool                     `json:"isEnabled"`
	CreatedAt     time.Time                `json:"createdAt"`
	ModifiedAt    time.Time                `json:"modifiedAt"`
}

type WorkspaceSkillView struct {
	Workspace        WorkspaceRef            `json:"workspace"`
	Artifact         artifact.ArtifactRef    `json:"artifact"`
	DefinitionDigest cryptoutil.Digest       `json:"definitionDigest"`
	SourceID         source.SourceID         `json:"sourceID"`
	Locator          basespec.Locator        `json:"locator"`
	Skill            WorkspaceSkillSummary   `json:"skill"`
	MarkdownBody     string                  `json:"markdownBody,omitempty"`
	RecordRevision   uint64                  `json:"recordRevision"`
	State            artifact.State          `json:"state"`
	ProjectionValid  bool                    `json:"projectionValid"`
	CatalogCurrent   bool                    `json:"catalogCurrent"`
	RuntimeDisabled  bool                    `json:"runtimeDisabled"`
	Diagnostics      []diagnostic.Diagnostic `json:"diagnostics,omitempty"`
}

type WorkspaceSkillLoadView struct {
	Workspace       WorkspaceRef            `json:"workspace"`
	CatalogRevision uint64                  `json:"catalogRevision"`
	Skills          []WorkspaceSkillView    `json:"skills"`
	Diagnostics     []diagnostic.Diagnostic `json:"diagnostics,omitempty"`
}

type GetWorkspaceRequest struct {
	Workspace WorkspaceRef `json:"workspace" required:"true"`
}

type GetWorkspaceResponse struct {
	Body *WorkspaceView
}

type ListWorkspacesRequest struct{}

type ListWorkspacesResponseBody struct {
	Workspaces []WorkspaceView `json:"workspaces"`
}

type ListWorkspacesResponse struct {
	Body *ListWorkspacesResponseBody
}

type GetWorkspaceCatalogRequest struct {
	Workspace WorkspaceRef `json:"workspace" required:"true"`
}

type GetWorkspaceCatalogResponse struct {
	Body *WorkspaceCatalogView
}

type GetWorkspaceArtifactRequest struct {
	Workspace WorkspaceRef         `json:"workspace" required:"true"`
	Artifact  artifact.ArtifactRef `json:"artifact"  required:"true"`
}

type GetWorkspaceArtifactResponse struct {
	Body *WorkspaceArtifactView
}

type ListWorkspaceArtifactsRequest struct {
	Workspace WorkspaceRef `json:"workspace" required:"true"`
}

type ListWorkspaceArtifactsResponseBody struct {
	Artifacts []WorkspaceArtifactView `json:"artifacts"`
}

type ListWorkspaceArtifactsResponse struct {
	Body *ListWorkspaceArtifactsResponseBody
}

type ListWorkspaceContextsRequest struct {
	Workspace WorkspaceRef `json:"workspace" required:"true"`
}

type ListWorkspaceContextsResponseBody struct {
	Contexts []WorkspaceContextView `json:"contexts"`
}

type ListWorkspaceContextsResponse struct {
	Body *ListWorkspaceContextsResponseBody
}

type LoadWorkspaceContextsRequestBody struct {
	Artifacts []artifact.ArtifactRef `json:"artifacts,omitempty"`
}

type LoadWorkspaceContextsRequest struct {
	Workspace WorkspaceRef `json:"workspace" required:"true"`
	Body      *LoadWorkspaceContextsRequestBody
}

type LoadWorkspaceContextsResponse struct {
	Body *WorkspaceContextInspectionView
}

type ComposeWorkspaceContextRequestBody struct {
	Artifacts []artifact.ArtifactRef `json:"artifacts,omitempty"`
}

type ComposeWorkspaceContextRequest struct {
	Workspace WorkspaceRef `json:"workspace" required:"true"`
	Body      *ComposeWorkspaceContextRequestBody
}

type ComposeWorkspaceContextResponse struct {
	Body *WorkspaceContextLoadPlan
}

type ListWorkspaceSkillsRequest struct {
	Workspace WorkspaceRef `json:"workspace" required:"true"`
}

type ListWorkspaceSkillsResponseBody struct {
	Skills []WorkspaceSkillView `json:"skills"`
}

type ListWorkspaceSkillsResponse struct {
	Body *ListWorkspaceSkillsResponseBody
}

type LoadWorkspaceSkillsRequestBody struct {
	Artifacts []artifact.ArtifactRef `json:"artifacts"`
}

type LoadWorkspaceSkillsRequest struct {
	Workspace WorkspaceRef `json:"workspace" required:"true"`
	Body      *LoadWorkspaceSkillsRequestBody
}

type LoadWorkspaceSkillsResponse struct {
	Body *WorkspaceSkillLoadView
}

type SetWorkspaceArtifactRuntimeDisabledRequestBody struct {
	ExpectedRevision uint64 `json:"expectedRevision" required:"true"`
	RuntimeDisabled  bool   `json:"runtimeDisabled"  required:"true"`
}

type SetWorkspaceArtifactRuntimeDisabledRequest struct {
	Workspace WorkspaceRef         `json:"workspace" required:"true"`
	Artifact  artifact.ArtifactRef `json:"artifact"  required:"true"`
	Body      *SetWorkspaceArtifactRuntimeDisabledRequestBody
}

type SetWorkspaceArtifactRuntimeDisabledResponse struct {
	Body *WorkspaceArtifactView
}

// requireRequestBody performs transport-shape checks only.
//
// Domain validation belongs to Workspace services. Store validation and
// lifecycle enforcement belong to Artifact Store.
func requireRequestBody[T any](
	request *T,
	bodyPresent bool,
	requireBody bool,
	subject string,
) error {
	if request == nil {
		return fmt.Errorf(
			"%w: %s request is required",
			spec.ErrInvalidWorkspace,
			subject,
		)
	}
	if requireBody && !bodyPresent {
		return fmt.Errorf(
			"%w: %s request body is required",
			spec.ErrInvalidWorkspace,
			subject,
		)
	}
	return nil
}
