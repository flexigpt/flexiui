import type { ArtifactCollection, ArtifactCollectionRef, ArtifactRecord, ArtifactRef } from '@/spec/artifact';
import type {
	MCPAuthHealth,
	MCPBundle,
	MCPReplaceBundleDocumentInput,
	MCPRuntimeServerID,
	MCPSecretKind,
	MCPSecretWriteResult,
	MCPServerData,
} from '@/spec/mcp_artifact';

import type { IMCPAggregateAPI } from '@/apis/interface';
import { rawJSONObjectToWails, requiredObject, requireNonBlankString } from '@/apis/wailsapi/transport';
import {
	ArtifactRefForRuntimeServerID,
	DeleteMCPServerSecret,
	GetMCPServerAuthHealth,
	PurgeMCPBundle,
	PutMCPServerSecret,
	RefreshMCPBundle,
	ReplaceMCPBundleDocument,
	RetireMCPBundle,
	RuntimeServerIDForArtifact,
	UpdateMCPServerInstallation,
	UpdateProtectedMCPBundleInstallation,
	UpdateProtectedMCPServerInstallation,
} from '@/apis/wailsjs/go/main/MCPAggregateWrapper';

function documentToWails(value: MCPReplaceBundleDocumentInput['document'], field: string): unknown {
	return rawJSONObjectToWails(JSON.stringify(value), field);
}

function registrationToWails(value: MCPReplaceBundleDocumentInput['registrations'][number], field: string): unknown {
	return {
		ArtifactID: value.artifactID,
		Subresource: value.subresource,
		Kind: value.kind,
		Enabled: value.enabled,
		...(value.data === undefined ? {} : { Data: rawJSONObjectToWails(value.data, `${field}.data`) }),
	};
}

export class WailsMCPAggregateAPI implements IMCPAggregateAPI {
	async runtimeServerIDForArtifact(artifact: ArtifactRef): Promise<MCPRuntimeServerID> {
		return requireNonBlankString(
			await RuntimeServerIDForArtifact(artifact as Parameters<typeof RuntimeServerIDForArtifact>[0]),
			'RuntimeServerIDForArtifact'
		);
	}

	async artifactRefForRuntimeServerID(server: MCPRuntimeServerID): Promise<ArtifactRef> {
		return requiredObject<ArtifactRef>(
			await ArtifactRefForRuntimeServerID(server as Parameters<typeof ArtifactRefForRuntimeServerID>[0]),
			'ArtifactRefForRuntimeServerID'
		);
	}

	async replaceMCPBundleDocument(input: MCPReplaceBundleDocumentInput): Promise<MCPBundle> {
		const response = await ReplaceMCPBundleDocument({
			Bundle: input.bundle,
			ExpectedCollectionRevision: input.expectedCollectionRevision,
			Document: documentToWails(input.document, 'ReplaceMCPBundleDocument.document'),
			Registrations: input.registrations.map((registration, index) =>
				registrationToWails(registration, `ReplaceMCPBundleDocument.registrations[${index}]`)
			),
			AllowProtected: false,
		} as Parameters<typeof ReplaceMCPBundleDocument>[0]);

		return requiredObject<MCPBundle>(response, 'ReplaceMCPBundleDocument');
	}

	async refreshMCPBundle(bundle: ArtifactCollectionRef): Promise<MCPBundle> {
		return requiredObject<MCPBundle>(
			await RefreshMCPBundle(bundle as Parameters<typeof RefreshMCPBundle>[0]),
			'RefreshMCPBundle'
		);
	}

	async retireMCPBundle(bundle: ArtifactCollectionRef, expectedRevision: number): Promise<ArtifactCollection> {
		return requiredObject<ArtifactCollection>(
			await RetireMCPBundle(bundle as Parameters<typeof RetireMCPBundle>[0], expectedRevision),
			'RetireMCPBundle'
		);
	}

	async purgeMCPBundle(bundle: ArtifactCollectionRef, expectedRevision: number): Promise<void> {
		await PurgeMCPBundle(bundle as Parameters<typeof PurgeMCPBundle>[0], expectedRevision);
	}

	async updateMCPServerInstallation(
		server: ArtifactRef,
		expectedArtifactRevision: number,
		data: MCPServerData
	): Promise<ArtifactRecord> {
		return requiredObject<ArtifactRecord>(
			await UpdateMCPServerInstallation(
				server as Parameters<typeof UpdateMCPServerInstallation>[0],
				expectedArtifactRevision,
				data as Parameters<typeof UpdateMCPServerInstallation>[2]
			),
			'UpdateMCPServerInstallation'
		);
	}

	async updateProtectedMCPBundleInstallation(
		bundle: ArtifactCollectionRef,
		expectedOverlayRevision: number,
		runtimeEnabled: boolean
	): Promise<void> {
		await UpdateProtectedMCPBundleInstallation(
			bundle as Parameters<typeof UpdateProtectedMCPBundleInstallation>[0],
			expectedOverlayRevision,
			runtimeEnabled
		);
	}

	async updateProtectedMCPServerInstallation(
		server: ArtifactRef,
		expectedOverlayRevision: number,
		runtimeEnabled: boolean,
		data: MCPServerData
	): Promise<void> {
		await UpdateProtectedMCPServerInstallation(
			server as Parameters<typeof UpdateProtectedMCPServerInstallation>[0],
			expectedOverlayRevision,
			runtimeEnabled,
			data as Parameters<typeof UpdateProtectedMCPServerInstallation>[3]
		);
	}

	async putMCPServerSecret(
		server: ArtifactRef,
		kind: MCPSecretKind,
		slot: string,
		secret: string
	): Promise<MCPSecretWriteResult> {
		return requiredObject<MCPSecretWriteResult>(
			await PutMCPServerSecret(server as Parameters<typeof PutMCPServerSecret>[0], kind, slot, secret),
			'PutMCPServerSecret'
		);
	}

	async deleteMCPServerSecret(server: ArtifactRef, kind: MCPSecretKind, slot: string): Promise<void> {
		await DeleteMCPServerSecret(server as Parameters<typeof DeleteMCPServerSecret>[0], kind, slot);
	}

	async getMCPServerAuthHealth(server: ArtifactRef): Promise<MCPAuthHealth> {
		return requiredObject<MCPAuthHealth>(
			await GetMCPServerAuthHealth(server as Parameters<typeof GetMCPServerAuthHealth>[0]),
			'GetMCPServerAuthHealth'
		);
	}
}
