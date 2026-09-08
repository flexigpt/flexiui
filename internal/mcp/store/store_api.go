package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sort"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/definition"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/installerapi"
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
	mcpPolicy "github.com/flexigpt/flexigpt-app/internal/mcp/runtime/policy"
	mcpOverlay "github.com/flexigpt/flexigpt-app/internal/mcp/store/overlay"
	mcpStorePolicy "github.com/flexigpt/flexigpt-app/internal/mcp/store/policy"
	mcpStoreServer "github.com/flexigpt/flexigpt-app/internal/mcp/store/server"
)

type StoreAPI struct {
	sources          compositionapi.SourceAPI
	collections      compositionapi.CollectionAPI
	artifacts        compositionapi.ArtifactAPI
	catalogs         compositionapi.CatalogAPI
	resources        compositionapi.ResourceAPI
	schemas          compositionapi.SchemaAPI
	managedArtifacts compositionapi.ManagedArtifactAPI
	protection       compositionapi.ProtectionAPI

	userRootID     root.RootID
	overlays       mcpOverlay.OverlayRepository
	secretCleaner  mcpStoreServer.SecretCleaner
	baselinePolicy mcpPolicy.MCPPolicy
}

func NewStoreAPI(
	sources compositionapi.SourceAPI,
	collections compositionapi.CollectionAPI,
	artifacts compositionapi.ArtifactAPI,
	catalogs compositionapi.CatalogAPI,
	resources compositionapi.ResourceAPI,
	schemas compositionapi.SchemaAPI,
	managedArtifacts compositionapi.ManagedArtifactAPI,
	protection compositionapi.ProtectionAPI,
	userRootID root.RootID,
	overlays mcpOverlay.OverlayRepository,
	secretCleaner mcpStoreServer.SecretCleaner,
	baselinePolicy mcpPolicy.MCPPolicy,
) (*StoreAPI, error) {
	if sources == nil ||
		collections == nil ||
		artifacts == nil ||
		catalogs == nil ||
		resources == nil ||
		schemas == nil ||
		managedArtifacts == nil ||
		protection == nil ||
		secretCleaner == nil {
		return nil, fmt.Errorf(
			"%w: MCP Store dependencies are incomplete",
			basespec.ErrInvalid,
		)
	}
	if userRootID != "" {
		if err := userRootID.Validate(); err != nil {
			return nil, err
		}
	}

	return &StoreAPI{
		sources:          sources,
		collections:      collections,
		artifacts:        artifacts,
		catalogs:         catalogs,
		resources:        resources,
		schemas:          schemas,
		managedArtifacts: managedArtifacts,
		protection:       protection,
		userRootID:       userRootID,
		overlays:         overlays,
		secretCleaner:    secretCleaner,
		baselinePolicy:   baselinePolicy,
	}, nil
}

type Bundle struct {
	Collection      collection.Collection        `json:"collection"`
	Data            CollectionData               `json:"data"`
	Attachment      collection.Attachment        `json:"attachment"`
	Source          source.Summary               `json:"source"`
	PackageAddress  source.ManagedPackageAddress `json:"packageAddress"`
	DocumentLocator basespec.Locator             `json:"documentLocator"`
}

type Registration struct {
	ArtifactID  artifact.ArtifactID
	Subresource basespec.SubresourceLocator
	Kind        artifact.ArtifactKind
	Enabled     bool
	Data        json.RawMessage
}

func (a *StoreAPI) Create(
	ctx context.Context,
	request CreateMCPBundleBody,
) (Bundle, error) {
	if a == nil {
		return Bundle{}, basespec.ErrClosed
	}
	if err := request.RootID.Validate(); err != nil {
		return Bundle{}, err
	}
	if a.userRootID != "" &&
		request.RootID != a.userRootID {
		return Bundle{}, fmt.Errorf(
			"%w: user MCP Bundles must be created in Root %q",
			basespec.ErrInvalid,
			a.userRootID,
		)
	}
	if err := a.requireBundleMutation(
		ctx,
		request.RootID,
		false,
	); err != nil {
		return Bundle{}, err
	}
	if err := request.CollectionID.Validate(); err != nil {
		return Bundle{}, err
	}
	if err := request.SourceID.Validate(); err != nil {
		return Bundle{}, err
	}

	document, parsedDocument, err := a.canonicalizeBundleBytes(ctx, request.Document)
	if err != nil {
		return Bundle{}, err
	}
	if err := validateCreateRegistrations(
		request.RootID,
		document,
		request.Registrations,
	); err != nil {
		return Bundle{}, err
	}

	packageAddress, err := PackageAddressForBundle(
		document.LogicalName,
		document.LogicalVersion,
	)
	if err != nil {
		return Bundle{}, err
	}
	sourceValue, createdSource, err := a.sources.CreateWithStatus(
		ctx,
		request.RootID,
		source.Draft{
			ID:          request.SourceID,
			StorageKey:  request.SourceStorageKey,
			Kind:        source.SourceKindManagedDirectory,
			DisplayName: displayName(document),
			Enabled:     true,
			Config:      json.RawMessage(jsonutil.EmptyObject),
		},
	)
	if err != nil {
		return Bundle{}, err
	}
	cleanupSource := func(cause error) error {
		if !createdSource {
			return cause
		}
		return errors.Join(
			cause,
			a.sources.Discard(
				context.WithoutCancel(ctx),
				request.RootID,
				request.SourceID,
				sourceValue.Revision,
			),
		)
	}

	collectionData, err := EncodeCollectionData(CollectionData{
		SchemaVersion:           artifactbuiltin.MCPSchemaVersion,
		DiscoveryPolicyRevision: artifactbuiltin.DecoderRevision,
		LogicalName:             document.LogicalName,
		LogicalVersion:          document.LogicalVersion,
		Labels:                  maps.Clone(document.Labels),
		ManagedSourceID:         request.SourceID,
	})
	if err != nil {
		return Bundle{}, cleanupSource(err)
	}
	attachmentData, err := EncodeAttachmentData(AttachmentData{
		SchemaVersion:  artifactbuiltin.MCPSchemaVersion,
		PackageAddress: packageAddress,
	})
	if err != nil {
		return Bundle{}, cleanupSource(err)
	}

	created, _, err := a.collections.Create(
		ctx,
		request.RootID,
		collection.Draft{
			ID:          request.CollectionID,
			Kind:        artifactbuiltin.BundleKind,
			DisplayName: displayName(document),
			Description: document.Description,
			Enabled:     true,
			Data:        collectionData,
		},
		[]collection.AttachmentDraft{{
			SourceID: request.SourceID,
			Role:     artifactbuiltin.ManagedAttachmentRole,
			Enabled:  true,
			Data:     attachmentData,
		}},
	)
	if err != nil {
		return Bundle{}, cleanupSource(err)
	}

	bundle, err := a.Get(ctx, created.Ref())
	if err != nil {
		return Bundle{}, err
	}
	if err := validateCreateBundleIntent(
		bundle,
		request,
		document,
		packageAddress,
	); err != nil {
		return Bundle{}, cleanupSource(err)
	}
	if _, err := a.replaceCanonicalDocument(
		ctx,
		ReplaceDocumentRequest{
			Bundle:                     bundle.Collection.Ref(),
			ExpectedCollectionRevision: bundle.Collection.Revision,
			Document:                   parsedDocument.Raw,
			Registrations:              request.Registrations,
			AllowProtected:             false,
		},
		parsedDocument,
		nil,
	); err != nil {
		return Bundle{}, err
	}
	return a.Get(ctx, created.Ref())
}

func (a *StoreAPI) List(
	ctx context.Context,
	rootID root.RootID,
) ([]Bundle, error) {
	values, err := a.collections.ListByRoot(ctx, rootID)
	if err != nil {
		return nil, err
	}
	output := make([]Bundle, 0)
	for _, value := range values {
		if value.Kind != artifactbuiltin.BundleKind {
			continue
		}
		bundle, err := a.Get(ctx, value.Ref())
		if err != nil {
			return nil, err
		}
		output = append(output, bundle)
	}
	sort.Slice(output, func(left, right int) bool {
		return output[left].Collection.ID <
			output[right].Collection.ID
	})
	return output, nil
}

// Refresh performs an explicit MCP configuration refresh. The runtime session
// is invalidated before publication so a live client cannot continue using
// configuration that the caller explicitly asked to re-evaluate.
func (a *StoreAPI) Refresh(
	ctx context.Context,
	ref collection.CollectionRef,
	allowProtected bool,
) (Bundle, error) {
	if a == nil {
		return Bundle{}, basespec.ErrClosed
	}
	if err := ref.Validate(); err != nil {
		return Bundle{}, err
	}
	if err := a.requireBundleMutation(
		ctx,
		ref.RootID,
		allowProtected,
	); err != nil {
		return Bundle{}, err
	}

	bundle, err := a.Get(ctx, ref)
	if err != nil {
		return Bundle{}, err
	}
	if !bundle.Collection.Enabled {
		return Bundle{}, fmt.Errorf(
			"%w: MCP Bundle %q is disabled",
			basespec.ErrConflict,
			ref.CollectionID,
		)
	}

	if _, err := a.catalogs.RefreshCollection(
		ctx,
		ref,
	); err != nil {
		return Bundle{}, err
	}
	return a.Get(ctx, ref)
}

// EnsureBuiltInCurrent avoids managed package republishing for a current
// protected Bundle, but repairs a missing or stale Catalog after startup,
// decoder changes, or interrupted prior work.
func (a *StoreAPI) EnsureBuiltInCurrent(
	ctx context.Context,
	ref collection.CollectionRef,
) error {
	if a == nil {
		return basespec.ErrClosed
	}
	if err := installerapi.RequirePrivileged(ctx); err != nil {
		return err
	}
	if err := ref.Validate(); err != nil {
		return err
	}
	if !a.protection.IsProtectedRoot(ref.RootID) {
		return fmt.Errorf(
			"%w: MCP Bundle %q is not protected",
			basespec.ErrProtected,
			ref.CollectionID,
		)
	}

	bundle, err := a.Get(ctx, ref)
	if err != nil {
		return err
	}
	if _, err := a.currentCatalog(ctx, bundle); err == nil {
		return nil
	} else if !errors.Is(err, basespec.ErrCatalogUnavailable) &&
		!errors.Is(err, basespec.ErrCatalogStale) {
		return err
	}

	_, err = a.Refresh(ctx, ref, true)
	return err
}

func (a *StoreAPI) Get(
	ctx context.Context,
	ref collection.CollectionRef,
) (Bundle, error) {
	if a == nil {
		return Bundle{}, basespec.ErrClosed
	}
	value, err := a.collections.Get(ctx, ref)
	if err != nil {
		return Bundle{}, err
	}
	if value.Kind != artifactbuiltin.BundleKind {
		return Bundle{}, fmt.Errorf(
			"%w: Collection %q is not an MCP Bundle",
			basespec.ErrCollectionNotFound,
			ref.CollectionID,
		)
	}

	data, err := DecodeCollectionData(value.Data)
	if err != nil {
		return Bundle{}, err
	}
	attachments, err := a.collections.ListAttachments(
		ctx,
		ref,
	)
	if err != nil {
		return Bundle{}, err
	}
	if len(attachments) != 1 {
		return Bundle{}, fmt.Errorf(
			"%w: MCP Bundle must have exactly one Source Attachment",
			basespec.ErrInvalid,
		)
	}
	attachment := attachments[0]
	if attachment.Role != artifactbuiltin.ManagedAttachmentRole &&
		attachment.Role != artifactbuiltin.BuiltInAttachmentRole {
		return Bundle{}, fmt.Errorf(
			"%w: unsupported MCP Attachment role %q",
			basespec.ErrInvalid,
			attachment.Role,
		)
	}
	attachmentData, err := DecodeAttachmentData(attachment.Data)
	if err != nil {
		return Bundle{}, err
	}
	documentLocator, err := attachmentData.DocumentLocator()
	if err != nil {
		return Bundle{}, err
	}
	if err := validateBundlePackageAddress(attachmentData.PackageAddress); err != nil {
		return Bundle{}, err
	}

	sourceValue, err := a.sources.Get(
		ctx,
		ref.RootID,
		attachment.SourceID,
	)
	if err != nil {
		return Bundle{}, err
	}
	if sourceValue.Kind != source.SourceKindManagedDirectory {
		return Bundle{}, fmt.Errorf(
			"%w: MCP Bundle requires a managed Source",
			basespec.ErrInvalid,
		)
	}
	if data.ManagedSourceID != "" &&
		data.ManagedSourceID != sourceValue.ID {
		return Bundle{}, fmt.Errorf(
			"%w: MCP Bundle managed Source ownership mismatch",
			basespec.ErrInvalid,
		)
	}

	return Bundle{
		Collection: value,
		Data:       data,
		Attachment: attachment,
		Source:     sourceValue,

		PackageAddress:  attachmentData.PackageAddress,
		DocumentLocator: documentLocator,
	}, nil
}

func (a *StoreAPI) CreateMCPBundle(
	ctx context.Context,
	request *CreateMCPBundleRequest,
) (*CreateMCPBundleResponse, error) {
	if err := requireStoreRequest(
		request,
		true,
		request != nil && request.Body != nil,
		"MCP Bundle creation",
	); err != nil {
		return nil, err
	}

	value, err := a.Create(ctx, *request.Body)
	if err != nil {
		return nil, wrapStoreError("create Bundle", err)
	}
	return &CreateMCPBundleResponse{Body: &value}, nil
}

func (a *StoreAPI) GetMCPBundle(
	ctx context.Context,
	request *GetMCPBundleRequest,
) (*GetMCPBundleResponse, error) {
	if err := requireStoreRequest(
		request,
		false,
		false,
		"MCP Bundle get",
	); err != nil {
		return nil, err
	}

	value, err := a.Get(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("get Bundle", err)
	}
	return &GetMCPBundleResponse{Body: &value}, nil
}

func (a *StoreAPI) ListMCPBundles(
	ctx context.Context,
	request *ListMCPBundlesRequest,
) (*ListMCPBundlesResponse, error) {
	if err := requireStoreRequest(
		request,
		false,
		false,
		"MCP Bundle list",
	); err != nil {
		return nil, err
	}

	values, err := a.List(ctx, request.RootID)
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
	if err := requireStoreRequest(
		request,
		false,
		false,
		"MCP Bundle document get",
	); err != nil {
		return nil, err
	}

	value, err := a.GetDocument(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("get Bundle document", err)
	}
	return &GetMCPBundleDocumentResponse{Body: &value}, nil
}

func (a *StoreAPI) ListMCPBundleServers(
	ctx context.Context,
	request *ListMCPBundleServersRequest,
) (*ListMCPBundleServersResponse, error) {
	if err := requireStoreRequest(
		request,
		false,
		false,
		"MCP Bundle Server list",
	); err != nil {
		return nil, err
	}

	values, err := a.ListServers(ctx, request.Bundle)
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
	if err := requireStoreRequest(
		request,
		false,
		false,
		"MCP Bundle Policy list",
	); err != nil {
		return nil, err
	}

	values, err := a.ListPolicies(ctx, request.Bundle)
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
	if err := requireStoreRequest(
		request,
		false,
		false,
		"MCP Server installation get",
	); err != nil {
		return nil, err
	}

	value, err := a.GetServerInstallation(ctx, request.Server)
	if err != nil {
		return nil, wrapStoreError("get Server installation", err)
	}
	return &GetMCPServerInstallationResponse{Body: &value}, nil
}

func (a *StoreAPI) InspectMCPServer(
	ctx context.Context,
	request *InspectMCPServerRequest,
) (*InspectMCPServerResponse, error) {
	if err := requireStoreRequest(
		request,
		false,
		false,
		"MCP Server inspection",
	); err != nil {
		return nil, err
	}

	value, err := a.InspectMCPServerForRuntime(ctx, request.Server)
	if err != nil {
		return nil, wrapStoreError("inspect Server", err)
	}
	return &InspectMCPServerResponse{Body: &value}, nil
}

func (a *StoreAPI) InspectMCPPolicy(
	ctx context.Context,
	request *InspectMCPPolicyRequest,
) (*InspectMCPPolicyResponse, error) {
	if err := requireStoreRequest(
		request,
		false,
		false,
		"MCP Policy inspection",
	); err != nil {
		return nil, err
	}

	value, err := a.InspectMCPPolicyForRuntime(ctx, request.Policy)
	if err != nil {
		return nil, wrapStoreError("inspect Policy", err)
	}
	return &InspectMCPPolicyResponse{Body: &value}, nil
}

func (a *StoreAPI) GetMCPBundleInstallation(
	ctx context.Context,
	request *GetMCPBundleInstallationRequest,
) (*GetMCPBundleInstallationResponse, error) {
	if err := requireStoreRequest(
		request,
		false,
		false,
		"MCP Bundle installation get",
	); err != nil {
		return nil, err
	}

	value, err := a.GetBundleInstallation(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("get Bundle installation", err)
	}
	return &GetMCPBundleInstallationResponse{Body: &value}, nil
}

func (a *StoreAPI) cleanupChangedServerInstallation(
	ctx context.Context,
	record artifact.Artifact,
	document mcpStoreServer.ServerDocument,
	after mcpStoreServer.ServerData,
) error {
	if record.Kind != artifactbuiltin.ServerKind {
		return nil
	}
	if err := mcpStoreServer.CleanupUnboundServerSecrets(
		ctx,
		record.Ref(),
		document,
		after,
		a.secretCleaner,
	); err != nil {
		return fmt.Errorf(
			"MCP server installation secret cleanup remains pending: %w",
			err,
		)
	}
	return nil
}

func (a *StoreAPI) cleanupRemovedServerInstallation(
	ctx context.Context,
	record artifact.Artifact,
) error {
	if record.Kind != artifactbuiltin.ServerKind {
		return nil
	}

	data, err := a.serverInstallationDataForCleanup(ctx, record)
	if err != nil {
		return err
	}
	if err := mcpStoreServer.CleanupRemovedServerSecrets(
		ctx,
		record.Ref(),
		data,
		a.secretCleaner,
	); err != nil {
		return fmt.Errorf(
			"MCP server removal secret cleanup remains pending: %w",
			err,
		)
	}
	return nil
}

func (a *StoreAPI) serverInstallationDataForCleanup(
	ctx context.Context,
	record artifact.Artifact,
) (mcpStoreServer.ServerData, error) {
	if record.Kind != artifactbuiltin.ServerKind {
		return mcpStoreServer.DefaultServerData(), nil
	}

	if !a.protection.IsProtectedRoot(record.RootID) {
		return mcpStoreServer.DecodeServerData(record.Data)
	}
	if a.overlays == nil {
		return mcpStoreServer.ServerData{}, fmt.Errorf(
			"%w: protected MCP overlay store is unavailable",
			basespec.ErrReferenceUnresolved,
		)
	}

	ovr, found, err := a.overlays.GetServerOverlay(
		ctx,
		record.Ref(),
	)
	if err != nil {
		return mcpStoreServer.ServerData{}, err
	}
	if !found {
		return mcpStoreServer.DefaultServerData(), nil
	}
	return ovr.ServerData, nil
}

// validateCreateRegistrations establishes all request-derived valid state
// before source or Collection mutation begins.
func validateCreateRegistrations(
	rootID root.RootID,
	document BundleDocument,
	registrations []Registration,
) error {
	definitions, err := definitionsForDocument(document)
	if err != nil {
		return err
	}
	values, err := registrationMap(registrations, definitions)
	if err != nil {
		return err
	}

	subresources := make([]basespec.SubresourceLocator, 0, len(values))
	for subresource := range values {
		subresources = append(subresources, subresource)
	}
	slices.Sort(subresources)

	for _, subresource := range subresources {
		registration := values[subresource]
		if registration.Kind != artifactbuiltin.ServerKind {
			continue
		}

		data, err := registrationData(registration)
		if err != nil {
			return err
		}
		serverDefinition, err := serverDocumentFromDefinition(
			definitions[subresource],
		)
		if err != nil {
			return err
		}
		serverData, err := mcpStoreServer.DecodeServerData(data)
		if err != nil {
			return err
		}
		if err := mcpStoreServer.ValidateServerDataForDocument(
			artifact.ArtifactRef{
				RootID:     rootID,
				ArtifactID: registration.ArtifactID,
			},
			serverDefinition,
			serverData,
		); err != nil {
			return err
		}
	}
	return nil
}

func definitionsForDocument(
	document BundleDocument,
) (
	map[basespec.SubresourceLocator]definition.Definition,
	error,
) {
	output := make(
		map[basespec.SubresourceLocator]definition.Definition,
		len(document.MCPServers)+len(document.BundleExtension.Policies),
	)
	for name := range document.MCPServers {
		serverDocument, err := ServerFromCanonicalBundle(document, name)
		if err != nil {
			return nil, err
		}
		value, err := mcpStoreServer.DefinitionForCanonicalServer(serverDocument)
		if err != nil {
			return nil, err
		}
		output[mcpStoreServer.ServerSubresource(
			basespec.LogicalName(name),
		)] = value
	}
	for name, policyDocument := range document.BundleExtension.Policies {
		value, err := mcpStorePolicy.DefinitionForCanonicalPolicy(policyDocument)
		if err != nil {
			return nil, err
		}
		output[mcpStorePolicy.PolicySubresource(
			basespec.LogicalName(name),
		)] = value
	}
	return output, nil
}

func validateCreateBundleIntent(
	value Bundle,
	request CreateMCPBundleBody,
	document BundleDocument,
	packageAddress source.ManagedPackageAddress,
) error {
	if value.Collection.RootID != request.RootID ||
		value.Collection.ID != request.CollectionID ||
		value.Collection.Kind != artifactbuiltin.BundleKind ||
		value.Source.ID != request.SourceID ||
		value.PackageAddress != packageAddress ||
		value.Data.ManagedSourceID != request.SourceID ||
		value.Data.LogicalName != document.LogicalName ||
		value.Data.LogicalVersion != document.LogicalVersion ||
		!maps.Equal(value.Data.Labels, document.Labels) {
		return fmt.Errorf(
			"%w: MCP Bundle creation intent differs from existing state",
			basespec.ErrConflict,
		)
	}
	return nil
}

func displayName(document BundleDocument) string {
	if document.DisplayName != "" {
		return document.DisplayName
	}
	return string(document.LogicalName)
}
