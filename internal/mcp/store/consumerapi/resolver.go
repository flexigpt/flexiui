package consumerapi

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/catalog"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/resource"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
	mcpDomain "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain"
	mcpDomainPolicy "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain/policy"
	mcpDomainServer "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain/server"
)

func (a *API) ResolveMCPServer(
	ctx context.Context,
	ref artifact.ArtifactRef,
) (mcpDomainServer.Resolved, error) {
	return a.resolveMCPServer(ctx, ref, true)
}

// InspectMCPServerForRuntime establishes Artifact, Collection, Catalog, Definition,
// installation, and policy validity without opening a Source snapshot or
// resolving a secret. It is deliberately for read-only status and setup
// projections. Connection and explicit runtime refresh must use
// ResolveMCPServer, which verifies exact current source bytes.
func (a *API) InspectMCPServerForRuntime(
	ctx context.Context,
	ref artifact.ArtifactRef,
) (mcpDomainServer.Resolved, error) {
	return a.resolveMCPServer(ctx, ref, false)
}

func (a *API) resolveMCPServer(
	ctx context.Context,
	ref artifact.ArtifactRef,
	verifySource bool,
) (mcpDomainServer.Resolved, error) {
	if a == nil {
		return mcpDomainServer.Resolved{}, basespec.ErrClosed
	}
	if err := ref.Validate(); err != nil {
		return mcpDomainServer.Resolved{}, err
	}

	record, err := a.artifacts.Get(ctx, ref)
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}
	if record.Kind != artifactbuiltin.ServerKind ||
		record.State != artifact.StateAvailable ||
		record.ResolvedDefinition == nil {
		return mcpDomainServer.Resolved{}, fmt.Errorf(
			"%w: MCP Server Artifact is not available",
			basespec.ErrReferenceUnresolved,
		)
	}

	bundle, err := a.Get(ctx, collection.CollectionRef{
		RootID:       record.RootID,
		CollectionID: record.CollectionID,
	})
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}

	snapshot, err := a.currentCatalog(ctx, bundle)
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}
	occurrence, err := currentServerOccurrence(snapshot, record)
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}

	definitionValue, err := mcpDomain.DefinitionForArtifact(snapshot, record)
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}
	document, err := mcpDomainServer.ServerDocumentFromDefinition(definitionValue)
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}

	sourceRevision := snapshot.SourceRevisions[record.Binding.SourceID]
	sourceGeneration := snapshot.SourceGenerations[record.Binding.SourceID]
	if sourceRevision == 0 || sourceGeneration == "" ||
		occurrence.SourceContentDigest == nil {
		return mcpDomainServer.Resolved{}, fmt.Errorf(
			"%w: MCP Server Source has no current Catalog state",
			basespec.ErrCatalogStale,
		)
	}

	resolvedResource, err := a.resources.ResolveArtifact(
		ctx,
		ref,
		resource.ResolveOptions{
			VerifySourceContent: verifySource,
		},
	)
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}
	if resolvedResource.Artifact.Revision != record.Revision ||
		resolvedResource.CatalogRevision != snapshot.Revision ||
		resolvedResource.Definition.Digest != definitionValue.Digest ||
		resolvedResource.Source.Revision != sourceRevision ||
		resolvedResource.SourceGeneration != sourceGeneration ||
		resolvedResource.Occurrence.SourceContentDigest == nil ||
		*resolvedResource.Occurrence.SourceContentDigest !=
			*occurrence.SourceContentDigest {
		return mcpDomainServer.Resolved{}, fmt.Errorf(
			"%w: MCP resource changed during resolution",
			basespec.ErrCatalogStale,
		)
	}

	installationData, installationRevision, _, runtimeEnabled, err := a.effectiveInstallation(
		ctx,
		bundle,
		record,
		document,
	)
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}

	policyValue, err := a.effectivePolicy(
		ctx,
		bundle,
		document,
		installationData.AdditionalPolicies,
	)
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}

	version, err := resolvedVersion(struct {
		ServerRef            artifact.ArtifactRef     `json:"server"`
		CollectionRef        collection.CollectionRef `json:"collection"`
		ArtifactRevision     uint64                   `json:"artifactRevision"`
		CatalogRevision      uint64                   `json:"catalogRevision"`
		DefinitionDigest     cryptoutil.Digest        `json:"definitionDigest"`
		SourceContentDigest  cryptoutil.Digest        `json:"sourceContentDigest"`
		SourceGeneration     string                   `json:"sourceGeneration"`
		InstallationRevision uint64                   `json:"installationRevision"`
		PolicyDigest         cryptoutil.Digest        `json:"policyDigest"`
	}{
		ServerRef:            ref,
		CollectionRef:        bundle.Collection.Ref(),
		ArtifactRevision:     record.Revision,
		CatalogRevision:      snapshot.Revision,
		DefinitionDigest:     *record.ResolvedDefinition,
		SourceContentDigest:  *occurrence.SourceContentDigest,
		SourceGeneration:     sourceGeneration,
		InstallationRevision: installationRevision,
		PolicyDigest:         policyValue.Digest,
	})
	if err != nil {
		return mcpDomainServer.Resolved{}, err
	}

	resolved := mcpDomainServer.Resolved{
		Server:               ref,
		Collection:           bundle.Collection.Ref(),
		ArtifactRevision:     record.Revision,
		CatalogRevision:      snapshot.Revision,
		DefinitionDigest:     *record.ResolvedDefinition,
		SourceContentDigest:  *occurrence.SourceContentDigest,
		SourceGeneration:     sourceGeneration,
		Document:             document,
		Installation:         installationData,
		Policy:               policyValue,
		InstallationRevision: installationRevision,
		RuntimeEnabled:       runtimeEnabled,
		BuiltIn: a.protection.IsProtectedRoot(
			ref.RootID,
		),
		Version: version,
	}
	if err := resolved.Validate(); err != nil {
		return mcpDomainServer.Resolved{}, err
	}
	return resolved, nil
}

func (a *API) currentCatalog(
	ctx context.Context,
	bundle Bundle,
) (catalog.Snapshot, error) {
	return a.catalogs.CurrentCatalog(
		ctx,
		bundle.Collection.Ref(),
	)
}

func currentServerOccurrence(
	snapshot catalog.Snapshot,
	record artifact.Artifact,
) (catalog.Occurrence, error) {
	key := catalog.OccurrenceKey{
		CollectionID:       record.CollectionID,
		SourceID:           record.Binding.SourceID,
		Locator:            record.Binding.Locator,
		SubresourceLocator: record.Binding.SubresourceLocator,
	}
	for _, occurrence := range snapshot.Occurrences {
		if occurrence.Key != key {
			continue
		}
		if occurrence.State != catalog.OccurrenceValid ||
			occurrence.Kind != artifactbuiltin.ServerKind ||
			occurrence.DefinitionDigest == nil ||
			occurrence.SourceContentDigest == nil ||
			record.ResolvedDefinition == nil ||
			*occurrence.DefinitionDigest != *record.ResolvedDefinition {
			break
		}
		return occurrence.Clone(), nil
	}
	return catalog.Occurrence{}, fmt.Errorf(
		"%w: MCP Server does not match its current Catalog occurrence",
		basespec.ErrCatalogStale,
	)
}

func (a *API) effectiveInstallation(
	ctx context.Context,
	bundle Bundle,
	record artifact.Artifact,
	document mcpDomainServer.ServerDocument,
) (
	installationData mcpDomainServer.ServerData,
	installationRevision uint64,
	installationEnabled bool,
	runtimeEnabled bool,
	err error,
) {
	if !a.protection.IsProtectedRoot(record.RootID) {
		data, err := mcpDomainServer.DecodeServerData(record.Data)
		if err != nil {
			return mcpDomainServer.ServerData{}, 0, false, false, err
		}
		if err := mcpDomainServer.ValidateServerDataForDocument(record.Ref(), document, data); err != nil {
			return mcpDomainServer.ServerData{}, 0, false, false, err
		}
		return data,
			record.Revision,
			record.Enabled,
			bundle.Collection.Enabled && record.Enabled,
			nil
	}

	if a.overlays == nil {
		return mcpDomainServer.ServerData{},
			0,
			false,
			false,
			fmt.Errorf(
				"%w: protected MCP installation overlay store is unavailable",
				basespec.ErrReferenceUnresolved,
			)
	}

	serverOverlay, found, err := a.overlays.GetServerOverlay(
		ctx,
		record.Ref(),
	)
	if err != nil {
		return mcpDomainServer.ServerData{}, 0, false, false, err
	}
	if !found {
		return mcpDomainServer.DefaultServerData(), 0, false, false, nil
	}
	if err := mcpDomainServer.ValidateServerDataForDocument(
		record.Ref(),
		document,
		serverOverlay.ServerData,
	); err != nil {
		return mcpDomainServer.ServerData{}, 0, false, false, err
	}
	bundleOverlay, bundleFound, err := a.overlays.GetBundleOverlay(
		ctx,
		record.RootID,
		record.CollectionID,
	)
	if err != nil {
		return mcpDomainServer.ServerData{}, 0, false, false, err
	}
	if !bundleFound {
		return serverOverlay.ServerData,
			serverOverlay.Revision,
			serverOverlay.RuntimeEnabled,
			false,
			nil
	}
	return serverOverlay.ServerData,
		serverOverlay.Revision,
		serverOverlay.RuntimeEnabled,
		bundle.Collection.Enabled &&
			record.Enabled &&
			serverOverlay.RuntimeEnabled &&
			bundleOverlay.RuntimeEnabled,
		nil
}

func (a *API) effectivePolicy(
	ctx context.Context,
	bundle Bundle,
	serverDocument mcpDomainServer.ServerDocument,
	additional []artifact.ArtifactRef,
) (mcpDomainPolicy.Effective, error) {
	values := make([]mcpDomainPolicy.MCPPolicy, 0, 1+len(additional))

	if reference := serverDocument.Extension.Policy; reference != nil {
		matches, err := a.policyBodiesByLogicalName(
			ctx,
			bundle.Collection.Ref(),
			reference.Ref,
		)
		if err != nil {
			return mcpDomainPolicy.Effective{}, err
		}
		switch len(matches) {
		case 0:
			if reference.Required {
				return mcpDomainPolicy.Effective{}, fmt.Errorf(
					"%w: required MCP policy %q is unavailable",
					basespec.ErrReferenceUnresolved,
					reference.Ref,
				)
			}
		case 1:
			values = append(values, matches[0])
		default:
			return mcpDomainPolicy.Effective{}, fmt.Errorf(
				"%w: MCP policy %q is ambiguous",
				basespec.ErrConflict,
				reference.Ref,
			)
		}
	}

	for _, ref := range sortedArtifactRefs(additional) {
		if ref.RootID != bundle.Collection.RootID {
			return mcpDomainPolicy.Effective{}, fmt.Errorf(
				"%w: additional MCP policy belongs to another Root",
				basespec.ErrInvalid,
			)
		}
		record, err := a.artifacts.Get(ctx, ref)
		if err != nil {
			return mcpDomainPolicy.Effective{}, err
		}
		if record.Kind != artifactbuiltin.PolicyKind ||
			record.CollectionID != bundle.Collection.ID ||
			!record.Enabled ||
			record.State != artifact.StateAvailable ||
			record.ResolvedDefinition == nil {
			return mcpDomainPolicy.Effective{}, fmt.Errorf(
				"%w: additional MCP policy %q is unavailable",
				basespec.ErrReferenceUnresolved,
				ref.ArtifactID,
			)
		}
		definitionValue, err := a.currentDefinitionForArtifact(ctx, record)
		if err != nil {
			return mcpDomainPolicy.Effective{}, err
		}
		body, err := mcpDomainPolicy.PolicyBodyFromDefinition(
			definitionValue,
		)
		if err != nil {
			return mcpDomainPolicy.Effective{}, err
		}
		values = append(values, body)
	}

	// BaselinePolicy is a fallback, not a mandatory policy floor. A primary
	// policy already replaces it. Apply the same rule when the installation
	// explicitly selects only additional policy Artifacts, otherwise allow,
	// trusted, auto, and Apps-enabled effects can never become effective.
	baseline := a.baselinePolicy
	if len(values) != 0 {
		baseline = values[0]
		values = values[1:]
	}
	return mcpDomainPolicy.Compose(baseline, values...)
}

func (a *API) policyBodiesByLogicalName(
	ctx context.Context,
	ref collection.CollectionRef,
	name basespec.LogicalName,
) ([]mcpDomainPolicy.MCPPolicy, error) {
	records, err := a.artifacts.ListByCollection(ctx, ref)
	if err != nil {
		return nil, err
	}
	output := make([]mcpDomainPolicy.MCPPolicy, 0)
	for _, record := range records {
		if record.Kind != artifactbuiltin.PolicyKind ||
			!record.Enabled ||
			record.State != artifact.StateAvailable ||
			record.ResolvedDefinition == nil {
			continue
		}
		definitionValue, err := a.currentDefinitionForArtifact(ctx, record)
		if err != nil {
			return nil, err
		}
		if definitionValue.LogicalName != name {
			continue
		}
		body, err := mcpDomainPolicy.PolicyBodyFromDefinition(
			definitionValue,
		)
		if err != nil {
			return nil, err
		}
		output = append(output, body)
	}
	return output, nil
}

func resolvedVersion(value any) (cryptoutil.Digest, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	canonical, err := jsonutil.Canonicalize(raw)
	if err != nil {
		return "", err
	}
	return cryptoutil.DigestBytes(canonical), nil
}

func sortedArtifactRefs(
	values []artifact.ArtifactRef,
) []artifact.ArtifactRef {
	output := append([]artifact.ArtifactRef(nil), values...)
	sort.Slice(output, func(left, right int) bool {
		if output[left].RootID != output[right].RootID {
			return output[left].RootID < output[right].RootID
		}
		return output[left].ArtifactID < output[right].ArtifactID
	})
	return output
}
