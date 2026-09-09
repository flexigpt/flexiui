import type { ArtifactCollectionRef, ArtifactRecord, ArtifactRef, ArtifactRootID } from '@/spec/artifact';
import type {
	MCPArtifactRegistration,
	MCPBundle,
	MCPBundleDocument,
	MCPBundleInstallation,
	MCPCreateBundleInput,
	MCPPolicyView,
	MCPServerInstallation,
	MCPServerResolved,
	MCPServerSchemaIdentity,
} from '@/spec/mcp_artifact';

import type { IMCPStoreAPI } from '@/apis/interface';
import {
	rawJSONFromWails,
	rawJSONObjectToWails,
	requiredObject,
	requireWailsBody,
	wailsObjectArrayOrEmpty,
} from '@/apis/wailsapi/transport';
import {
	CreateMCPBundle,
	GetMCPBundle,
	GetMCPBundleDocument,
	GetMCPBundleInstallation,
	GetMCPServerInstallation,
	GetMCPServerSchemaIdentity,
	InspectMCPPolicy,
	InspectMCPServer,
	ListMCPBundlePolicies,
	ListMCPBundles,
	ListMCPBundleServers,
	ListMCPBundlesForManagement,
	UpdateBundleEnabled,
} from '@/apis/wailsjs/go/main/MCPStoreWrapper';
import type { consumerapi as wailsConsumerAPI } from '@/apis/wailsjs/go/models';

function registrationToWails(value: MCPArtifactRegistration, field: string): unknown {
	return {
		ArtifactID: value.artifactID,
		Subresource: value.subresource,
		Kind: value.kind,
		Enabled: value.enabled,
		...(value.data === undefined ? {} : { Data: rawJSONObjectToWails(value.data, `${field}.data`) }),
	};
}

function documentToWails(value: MCPBundleDocument, field: string): unknown {
	return rawJSONObjectToWails(JSON.stringify(value), field);
}

export class WailsMCPStoreAPI implements IMCPStoreAPI {
	async listMCPBundlesForManagement(): Promise<MCPBundle[]> {
		return wailsObjectArrayOrEmpty<MCPBundle>(await ListMCPBundlesForManagement(), 'ListMCPBundlesForManagement');
	}

	async getMCPServerSchemaIdentity(): Promise<MCPServerSchemaIdentity> {
		return requiredObject<MCPServerSchemaIdentity>(await GetMCPServerSchemaIdentity(), 'GetMCPServerSchemaIdentity');
	}

	async createMCPBundle(input: MCPCreateBundleInput): Promise<MCPBundle> {
		const response = await CreateMCPBundle({
			body: {
				RootID: input.rootID,
				CollectionID: input.collectionID,
				SourceID: input.sourceID,
				SourceStorageKey: input.sourceStorageKey,
				Document: documentToWails(input.document, 'CreateMCPBundle.document'),
				Registrations: input.registrations.map((registration, index) =>
					registrationToWails(registration, `CreateMCPBundle.registrations[${index}]`)
				),
			},
		} as wailsConsumerAPI.CreateMCPBundleRequest);

		return requiredObject<MCPBundle>(requireWailsBody(response.body, 'CreateMCPBundle.body'), 'CreateMCPBundle.body');
	}

	async getMCPBundle(bundle: ArtifactCollectionRef): Promise<MCPBundle> {
		const response = await GetMCPBundle({
			bundle,
		} as wailsConsumerAPI.GetMCPBundleRequest);

		return requiredObject<MCPBundle>(requireWailsBody(response.body, 'GetMCPBundle.body'), 'GetMCPBundle.body');
	}

	async listMCPBundles(rootID: ArtifactRootID): Promise<MCPBundle[]> {
		const response = await ListMCPBundles({
			rootID,
		} as wailsConsumerAPI.ListMCPBundlesRequest);

		const body = requiredObject<{ bundles?: MCPBundle[] }>(
			requireWailsBody(response.body, 'ListMCPBundles.body'),
			'ListMCPBundles.body'
		);

		return wailsObjectArrayOrEmpty<MCPBundle>(body.bundles, 'ListMCPBundles.body.bundles');
	}

	async getMCPBundleDocument(bundle: ArtifactCollectionRef): Promise<MCPBundleDocument> {
		const response = await GetMCPBundleDocument({
			bundle,
		} as wailsConsumerAPI.GetMCPBundleDocumentRequest);

		return requiredObject<MCPBundleDocument>(
			requireWailsBody(response.body, 'GetMCPBundleDocument.body'),
			'GetMCPBundleDocument.body'
		);
	}

	async listMCPBundleServers(bundle: ArtifactCollectionRef): Promise<ArtifactRecord[]> {
		const response = await ListMCPBundleServers({
			bundle,
		} as wailsConsumerAPI.ListMCPBundleServersRequest);

		const body = requiredObject<{ servers?: ArtifactRecord[] }>(
			requireWailsBody(response.body, 'ListMCPBundleServers.body'),
			'ListMCPBundleServers.body'
		);

		return wailsObjectArrayOrEmpty<ArtifactRecord>(body.servers, 'ListMCPBundleServers.body.servers');
	}

	async listMCPBundlePolicies(bundle: ArtifactCollectionRef): Promise<ArtifactRecord[]> {
		const response = await ListMCPBundlePolicies({
			bundle,
		} as wailsConsumerAPI.ListMCPBundlePoliciesRequest);

		const body = requiredObject<{ policies?: ArtifactRecord[] }>(
			requireWailsBody(response.body, 'ListMCPBundlePolicies.body'),
			'ListMCPBundlePolicies.body'
		);

		return wailsObjectArrayOrEmpty<ArtifactRecord>(body.policies, 'ListMCPBundlePolicies.body.policies');
	}

	async getMCPBundleInstallation(bundle: ArtifactCollectionRef): Promise<MCPBundleInstallation> {
		const response = await GetMCPBundleInstallation({
			bundle,
		} as wailsConsumerAPI.GetMCPBundleInstallationRequest);

		return requiredObject<MCPBundleInstallation>(
			requireWailsBody(response.body, 'GetMCPBundleInstallation.body'),
			'GetMCPBundleInstallation.body'
		);
	}

	async updateMCPBundleEnabled(
		bundle: ArtifactCollectionRef,
		expectedRevision: number,
		enabled: boolean
	): Promise<MCPBundle> {
		return requiredObject<MCPBundle>(
			await UpdateBundleEnabled(bundle as Parameters<typeof UpdateBundleEnabled>[0], expectedRevision, enabled),
			'UpdateBundleEnabled'
		);
	}

	async getMCPServerInstallation(server: ArtifactRef): Promise<MCPServerInstallation> {
		const response = await GetMCPServerInstallation({
			server,
		} as wailsConsumerAPI.GetMCPServerInstallationRequest);

		return requiredObject<MCPServerInstallation>(
			requireWailsBody(response.body, 'GetMCPServerInstallation.body'),
			'GetMCPServerInstallation.body'
		);
	}

	async inspectMCPServer(server: ArtifactRef): Promise<MCPServerResolved> {
		const response = await InspectMCPServer({
			server,
		} as wailsConsumerAPI.InspectMCPServerRequest);

		return requiredObject<MCPServerResolved>(
			requireWailsBody(response.body, 'InspectMCPServer.body'),
			'InspectMCPServer.body'
		);
	}

	async inspectMCPPolicy(policy: ArtifactRef): Promise<MCPPolicyView> {
		const response = await InspectMCPPolicy({
			policy,
		} as wailsConsumerAPI.InspectMCPPolicyRequest);

		const value = requiredObject<MCPPolicyView>(
			requireWailsBody(response.body, 'InspectMCPPolicy.body'),
			'InspectMCPPolicy.body'
		);
		const definition = requireWailsBody(value.definition, 'InspectMCPPolicy.body.definition');

		return {
			...value,
			definition: {
				...definition,
				body: rawJSONFromWails(definition.body, 'InspectMCPPolicy.body.definition.body'),
			},
		};
	}
}
