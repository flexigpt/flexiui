import type { ArtifactAddress, ArtifactDiagnostic, ArtifactRef, ArtifactRootID } from '@/spec/artifact';
import { ArtifactAdoptionMode, ArtifactState } from '@/spec/artifact';
import type {
	AdoptSkillBody,
	CreateManagedSkillBody,
	CreateManagedSkillResult,
	CreateSkillBundleBody,
	ManagedSkillDocumentView,
	PinSkillBody,
	RegisterSkillBundleDirectoryInput,
	RetireSkillBundleResult,
	SetSkillEnabledBody,
	SkillArtifactView,
	SkillBundleRef,
	SkillBundleView,
	SkillDocumentInput,
	UpdateSkillBundleBody,
} from '@/spec/skill';
import { SkillBundleAttachmentRole } from '@/spec/skill';

import type { ISkillStoreAPI } from '@/apis/interface';
import { byteArrayToWails, enumFromWails, requireWailsBody, wailsObjectArrayOrEmpty } from '@/apis/wailsapi/transport';
import {
	AdoptSkill,
	CreateManagedSkill,
	CreateSkillBundle,
	GetManagedSkillDocument,
	GetSkillBundle,
	ListBundleSkills,
	ListSkillBundles,
	ListSkillBundlesForManagement,
	PinSkill,
	PurgeSkill,
	PurgeSkillBundle,
	RefreshSkillBundle,
	RegisterSkillBundleDirectory,
	RetireSkillBundle,
	SetSkillEnabled,
	UnadoptSkill,
	UpdateSkillBundle,
} from '@/apis/wailsjs/go/main/SkillStoreWrapper';
import type {
	artifact as wailsArtifact,
	consumerapi as wailsConsumerAPI,
	domain as wailsDomain,
} from '@/apis/wailsjs/go/models';

function skillArtifactFromWails(artifactValue: wailsArtifact.Artifact, field = 'skillArtifact'): SkillArtifactView {
	const artifact = requireWailsBody(artifactValue, field);
	const ref: ArtifactRef = {
		rootID: artifact.rootID,
		artifactID: artifact.id,
	};

	const address: ArtifactAddress = {
		...ref,
		collectionID: artifact.collectionID,
		kind: artifact.kind,
	};

	return {
		artifact: ref,
		address,
		revision: artifact.revision,
		name: artifact.name,
		kind: artifact.kind,
		enabled: artifact.enabled,
		adoption: enumFromWails(artifact.adoption, ArtifactAdoptionMode, `${field}.adoption`),
		state: enumFromWails(artifact.state, ArtifactState, `${field}.state`),
		binding: requireWailsBody(artifact.binding, `${field}.binding`),
		definitionDigest: artifact.resolvedDefinition ?? undefined,
		diagnostics:
			artifact.diagnostics === null || artifact.diagnostics === undefined
				? undefined
				: wailsObjectArrayOrEmpty<ArtifactDiagnostic>(artifact.diagnostics, `${field}.diagnostics`),
		createdAt: artifact.createdAt,
		modifiedAt: artifact.modifiedAt,
	};
}

function skillBundleFromWails(bundleValue: wailsDomain.SkillBundle, field = 'skillBundle'): SkillBundleView {
	const bundle = requireWailsBody(bundleValue, field);
	const collection = requireWailsBody(bundle.collection, `${field}.collection`);
	const data = requireWailsBody(bundle.data, `${field}.data`);
	const sources = wailsObjectArrayOrEmpty<wailsDomain.SkillBundle['sources'][number]>(
		bundle.sources,
		`${field}.sources`
	);
	const attachments = wailsObjectArrayOrEmpty<wailsDomain.SkillBundle['attachments'][number]>(
		bundle.attachments,
		`${field}.attachments`
	);
	const sourcesByID = new Map(sources.map(source => [source.id, source] as const));

	return {
		bundle: {
			rootID: collection.rootID,
			collectionID: collection.id,
		},
		revision: collection.revision,
		displayName: collection.displayName,
		description: collection.description ?? undefined,
		enabled: collection.enabled,
		retiredAt: collection.retiredAt ?? undefined,
		logicalName: data.logicalName,
		logicalVersion: data.logicalVersion ?? undefined,
		labels: data.labels ?? undefined,
		managedSourceID: data.managedSourceID ?? undefined,
		attachments: attachments.map(attachment => {
			const source = sourcesByID.get(attachment.sourceID);

			return {
				sourceID: attachment.sourceID,
				revision: attachment.revision,
				role: enumFromWails(attachment.role, SkillBundleAttachmentRole, `${field}.attachment.role`),
				enabled: attachment.enabled,
				sourceDisplayName: source?.displayName,
				sourceKind: source?.kind,
			};
		}),
		createdAt: collection.createdAt,
		modifiedAt: collection.modifiedAt,
	};
}

export class WailsSkillStoreAPI implements ISkillStoreAPI {
	async createSkillBundle(rootID: ArtifactRootID, body: CreateSkillBundleBody): Promise<SkillBundleView> {
		const managedSourceID = body.managedSourceID?.trim() || undefined;
		const managedSourceStorageKey = body.managedSourceStorageKey?.trim() || undefined;

		if (Boolean(managedSourceID) !== Boolean(managedSourceStorageKey)) {
			throw new Error('Managed Source ID and managed Source storage key must be supplied together.');
		}

		const response = await CreateSkillBundle({
			body: {
				RootID: rootID,
				CollectionID: body.collectionID,
				DisplayName: body.displayName,
				Description: body.description ?? '',
				Enabled: body.enabled,
				LogicalName: body.logicalName,
				LogicalVersion: body.logicalVersion ?? '',
				Labels: body.labels ?? {},
				ManagedSourceID: managedSourceID ?? '',
				ManagedSourceStorageKey: managedSourceStorageKey ?? '',
				Attachments: (body.attachments ?? []).map(attachment => ({
					sourceId: attachment.sourceID,
					role: attachment.role,
					enabled: attachment.enabled,
					discoveryRoot: attachment.discoveryRoot,
					expectedMemberDigests: attachment.expectedMemberDigests ?? {},
				})),
			},
		} as wailsConsumerAPI.CreateSkillBundleRequest);

		return skillBundleFromWails(requireWailsBody(response.body, 'CreateSkillBundle.body'), 'CreateSkillBundle.body');
	}

	async getSkillBundle(bundle: SkillBundleRef): Promise<SkillBundleView> {
		const response = await GetSkillBundle({
			bundle,
		} as wailsConsumerAPI.GetSkillBundleRequest);

		return skillBundleFromWails(requireWailsBody(response.body, 'GetSkillBundle.body'), 'GetSkillBundle.body');
	}

	async listSkillBundles(rootID: ArtifactRootID): Promise<SkillBundleView[]> {
		const response = await ListSkillBundles({
			rootID,
		} as wailsConsumerAPI.ListSkillBundlesRequest);
		const body = requireWailsBody(response.body, 'ListSkillBundles.body');

		return wailsObjectArrayOrEmpty<wailsDomain.SkillBundle>(body.bundles, 'ListSkillBundles.body.bundles').map(
			(bundle, index) => skillBundleFromWails(bundle, `ListSkillBundles.body.bundles[${index}]`)
		);
	}

	async listSkillBundlesForManagement(): Promise<SkillBundleView[]> {
		const bundles = await ListSkillBundlesForManagement();

		return wailsObjectArrayOrEmpty<wailsDomain.SkillBundle>(bundles, 'ListSkillBundlesForManagement').map(
			(bundle, index) => skillBundleFromWails(bundle, `ListSkillBundlesForManagement[${index}]`)
		);
	}

	async registerSkillBundleDirectory(
		bundle: SkillBundleRef,
		input: RegisterSkillBundleDirectoryInput
	): Promise<SkillBundleView> {
		const response = await RegisterSkillBundleDirectory({
			bundle,
			expectedCollectionRevision: input.expectedCollectionRevision,
			rootPath: input.rootPath,
			sourceDisplayName: input.sourceDisplayName,
		} as wailsConsumerAPI.RegisterSkillBundleDirectoryRequest);

		return skillBundleFromWails(
			requireWailsBody(response.body, 'RegisterSkillBundleDirectory.body'),
			'RegisterSkillBundleDirectory.body'
		);
	}

	async updateSkillBundle(bundle: SkillBundleRef, body: UpdateSkillBundleBody): Promise<SkillBundleView> {
		const response = await UpdateSkillBundle({
			body: {
				Bundle: bundle,
				ExpectedRevision: body.expectedRevision,
				DisplayName: body.displayName,
				Description: body.description ?? '',
				Enabled: body.enabled,
			},
		} as wailsConsumerAPI.UpdateSkillBundleRequest);

		return skillBundleFromWails(requireWailsBody(response.body, 'UpdateSkillBundle.body'), 'UpdateSkillBundle.body');
	}

	async retireSkillBundle(bundle: SkillBundleRef, expectedRevision: number): Promise<RetireSkillBundleResult> {
		const response = await RetireSkillBundle({
			bundle,
			expectedRevision,
		} as wailsConsumerAPI.RetireSkillBundleRequest);
		const retired = requireWailsBody(response.body, 'RetireSkillBundle.body');

		return {
			bundle,
			revision: retired.revision,
		};
	}

	async purgeSkillBundle(bundle: SkillBundleRef, expectedRevision: number): Promise<SkillBundleRef> {
		await PurgeSkillBundle({
			bundle,
			expectedRevision,
		} as wailsConsumerAPI.PurgeSkillBundleRequest);

		return bundle;
	}

	async refreshSkillBundle(bundle: SkillBundleRef): Promise<void> {
		await RefreshSkillBundle({
			bundle,
		} as wailsConsumerAPI.RefreshSkillBundleRequest);
	}

	async listSkillBundleArtifacts(bundle: SkillBundleRef): Promise<SkillArtifactView[]> {
		const response = await ListBundleSkills({
			bundle,
		} as wailsConsumerAPI.ListBundleSkillsRequest);
		const body = requireWailsBody(response.body, 'ListBundleSkills.body');

		return wailsObjectArrayOrEmpty<wailsArtifact.Artifact>(body.skills, 'ListBundleSkills.body.skills').map(
			(artifact, index) => skillArtifactFromWails(artifact, `ListBundleSkills.body.skills[${index}]`)
		);
	}

	async createManagedSkill(bundle: SkillBundleRef, body: CreateManagedSkillBody): Promise<CreateManagedSkillResult> {
		const response = await CreateManagedSkill({
			body: {
				Bundle: bundle,
				ExpectedCollectionRevision: body.expectedCollectionRevision,
				ExpectedArtifactRevision: body.expectedArtifactRevision ?? 0,
				ArtifactID: body.artifactID,
				SkillName: body.skillName,
				SKILLMD: body.skillMD === undefined ? [] : byteArrayToWails(body.skillMD),
				...(body.document === undefined ? {} : { Document: body.document }),
				Files: (body.files ?? []).map(file => ({
					locator: file.locator,
					content: byteArrayToWails(file.content),
				})),
				Enabled: body.enabled,
			},
		} as wailsConsumerAPI.CreateManagedSkillRequest);

		const result = requireWailsBody(response.body, 'CreateManagedSkill.body');

		return {
			artifact: skillArtifactFromWails(result.Artifact, 'CreateManagedSkill.body.Artifact'),
			address: requireWailsBody(result.Address, 'CreateManagedSkill.body.Address'),
		};
	}

	async getManagedSkillDocument(artifact: ArtifactRef): Promise<ManagedSkillDocumentView> {
		const response = await GetManagedSkillDocument({
			artifact,
		} as wailsConsumerAPI.GetManagedSkillDocumentRequest);
		const body = requireWailsBody(response.body, 'GetManagedSkillDocument.body');

		return {
			artifact: skillArtifactFromWails(body.artifact, 'GetManagedSkillDocument.body.artifact'),
			document: requireWailsBody(body.document, 'GetManagedSkillDocument.body.document') as SkillDocumentInput,
		};
	}

	async adoptSkill(bundle: SkillBundleRef, body: AdoptSkillBody): Promise<SkillArtifactView> {
		const response = await AdoptSkill({
			body: {
				Bundle: bundle,
				Occurrence: body.occurrence,
				ArtifactID: body.artifactID,
				ExpectedCatalogRevision: body.expectedCatalogRevision,
				Name: body.name,
				Enabled: body.enabled,
			},
		} as wailsConsumerAPI.AdoptSkillRequest);

		return skillArtifactFromWails(requireWailsBody(response.body, 'AdoptSkill.body'), 'AdoptSkill.body');
	}

	async pinSkill(bundle: SkillBundleRef, body: PinSkillBody): Promise<SkillArtifactView> {
		const response = await PinSkill({
			body: {
				Bundle: bundle,
				ExpectedCollectionRevision: body.expectedCollectionRevision,
				ArtifactID: body.artifactID,
				Binding: body.binding,
				Name: body.name,
				Enabled: body.enabled,
			},
		} as wailsConsumerAPI.PinSkillRequest);

		return skillArtifactFromWails(requireWailsBody(response.body, 'PinSkill.body'), 'PinSkill.body');
	}

	async setSkillEnabled(artifact: ArtifactRef, body: SetSkillEnabledBody): Promise<SkillArtifactView> {
		const response = await SetSkillEnabled({
			body: {
				Artifact: artifact,
				ExpectedRevision: body.expectedRevision,
				Enabled: body.enabled,
			},
		} as wailsConsumerAPI.SetSkillEnabledRequest);

		return skillArtifactFromWails(requireWailsBody(response.body, 'SetSkillEnabled.body'), 'SetSkillEnabled.body');
	}

	async unadoptSkill(artifact: ArtifactRef, expectedRevision: number, suppress: boolean): Promise<ArtifactRef> {
		await UnadoptSkill({
			artifact,
			expectedRevision,
			suppress,
		} as wailsConsumerAPI.UnadoptSkillRequest);

		return artifact;
	}

	async purgeSkill(artifact: ArtifactRef, expectedRevision: number): Promise<ArtifactRef> {
		await PurgeSkill({
			artifact,
			expectedRevision,
		} as wailsConsumerAPI.PurgeSkillRequest);

		return artifact;
	}
}
