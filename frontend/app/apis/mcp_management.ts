import type { ArtifactCollection, ArtifactCollectionRef, ArtifactRecord, ArtifactRef } from '@/spec/artifact';
import type {
	MCPAuthHealth,
	MCPBundle,
	MCPBundleDocument,
	MCPBundleInstallation,
	MCPCreateBundleInput,
	MCPPolicyView,
	MCPReplaceBundleDocumentInput,
	MCPRuntimeServerID,
	MCPSecretKind,
	MCPSecretWriteResult,
	MCPServerData,
	MCPServerInstallation,
	MCPServerResolved,
	MCPServerSchemaIdentity,
} from '@/spec/mcp_artifact';

import type { IMCPAggregateAPI, IMCPStoreAPI } from '@/apis/interface';

export class MCPManagementAPI {
	constructor(
		// oxlint-disable-next-line typescript/parameter-properties
		private readonly store: IMCPStoreAPI,
		// oxlint-disable-next-line typescript/parameter-properties
		private readonly aggregate: IMCPAggregateAPI
	) {}

	listMCPBundlesForManagement(): Promise<MCPBundle[]> {
		return this.store.listMCPBundlesForManagement();
	}

	getMCPServerSchemaIdentity(): Promise<MCPServerSchemaIdentity> {
		return this.store.getMCPServerSchemaIdentity();
	}

	createMCPBundle(input: MCPCreateBundleInput): Promise<MCPBundle> {
		return this.store.createMCPBundle(input);
	}

	getMCPBundle(bundle: ArtifactCollectionRef): Promise<MCPBundle> {
		return this.store.getMCPBundle(bundle);
	}

	getMCPBundleDocument(bundle: ArtifactCollectionRef): Promise<MCPBundleDocument> {
		return this.store.getMCPBundleDocument(bundle);
	}

	getMCPBundleInstallation(bundle: ArtifactCollectionRef): Promise<MCPBundleInstallation> {
		return this.store.getMCPBundleInstallation(bundle);
	}

	listMCPBundleServers(bundle: ArtifactCollectionRef): Promise<ArtifactRecord[]> {
		return this.store.listMCPBundleServers(bundle);
	}

	listMCPBundlePolicies(bundle: ArtifactCollectionRef): Promise<ArtifactRecord[]> {
		return this.store.listMCPBundlePolicies(bundle);
	}

	updateMCPBundleEnabled(
		bundle: ArtifactCollectionRef,
		expectedRevision: number,
		enabled: boolean
	): Promise<MCPBundle> {
		return this.store.updateMCPBundleEnabled(bundle, expectedRevision, enabled);
	}

	retireMCPBundle(bundle: ArtifactCollectionRef, expectedRevision: number): Promise<ArtifactCollection> {
		return this.aggregate.retireMCPBundle(bundle, expectedRevision);
	}

	purgeMCPBundle(bundle: ArtifactCollectionRef, expectedRevision: number): Promise<void> {
		return this.aggregate.purgeMCPBundle(bundle, expectedRevision);
	}

	getMCPServerInstallation(server: ArtifactRef): Promise<MCPServerInstallation> {
		return this.store.getMCPServerInstallation(server);
	}

	inspectMCPServer(server: ArtifactRef): Promise<MCPServerResolved> {
		return this.store.inspectMCPServer(server);
	}

	inspectMCPPolicy(policy: ArtifactRef): Promise<MCPPolicyView> {
		return this.store.inspectMCPPolicy(policy);
	}

	runtimeServerIDForArtifact(artifact: ArtifactRef): Promise<MCPRuntimeServerID> {
		return this.aggregate.runtimeServerIDForArtifact(artifact);
	}

	artifactRefForRuntimeServerID(server: MCPRuntimeServerID): Promise<ArtifactRef> {
		return this.aggregate.artifactRefForRuntimeServerID(server);
	}

	refreshMCPBundle(bundle: ArtifactCollectionRef): Promise<MCPBundle> {
		return this.aggregate.refreshMCPBundle(bundle);
	}

	replaceMCPBundleDocument(input: MCPReplaceBundleDocumentInput): Promise<MCPBundle> {
		return this.aggregate.replaceMCPBundleDocument(input);
	}

	updateMCPServerInstallation(
		server: ArtifactRef,
		expectedArtifactRevision: number,
		data: MCPServerData
	): Promise<ArtifactRecord> {
		return this.aggregate.updateMCPServerInstallation(server, expectedArtifactRevision, data);
	}

	updateProtectedMCPBundleInstallation(
		bundle: ArtifactCollectionRef,
		expectedOverlayRevision: number,
		runtimeEnabled: boolean
	): Promise<void> {
		return this.aggregate.updateProtectedMCPBundleInstallation(bundle, expectedOverlayRevision, runtimeEnabled);
	}

	updateProtectedMCPServerInstallation(
		server: ArtifactRef,
		expectedOverlayRevision: number,
		runtimeEnabled: boolean,
		data: MCPServerData
	): Promise<void> {
		return this.aggregate.updateProtectedMCPServerInstallation(server, expectedOverlayRevision, runtimeEnabled, data);
	}

	putMCPServerSecret(
		server: ArtifactRef,
		kind: MCPSecretKind,
		slot: string,
		secret: string
	): Promise<MCPSecretWriteResult> {
		return this.aggregate.putMCPServerSecret(server, kind, slot, secret);
	}

	deleteMCPServerSecret(server: ArtifactRef, kind: MCPSecretKind, slot: string): Promise<void> {
		return this.aggregate.deleteMCPServerSecret(server, kind, slot);
	}

	getMCPServerAuthHealth(server: ArtifactRef): Promise<MCPAuthHealth> {
		return this.aggregate.getMCPServerAuthHealth(server);
	}
}
