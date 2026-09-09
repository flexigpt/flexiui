import type { ArtifactRef, ArtifactSourceBinding, ArtifactSourceID } from '@/spec/artifact';
import type {
	AdoptWorkspaceOccurrenceBody,
	AttachWorkspaceSourceBody,
	CreateEmptyWorkspaceInput,
	CreateFilesystemWorkspaceInput,
	DetachWorkspaceSourceBody,
	PinWorkspaceArtifactBody,
	RegisterWorkspaceDirectoryInput,
	RetireWorkspaceResult,
	SetWorkspaceArtifactEnabledBody,
	SetWorkspacePrimarySourceBody,
	SuppressWorkspaceBindingBody,
	UnadoptWorkspaceArtifactBody,
	UnadoptWorkspaceArtifactResult,
	UnsuppressWorkspaceBindingResult,
	UpdateWorkspaceAttachmentBody,
	UpdateWorkspaceBody,
	WorkspaceArtifactView,
	WorkspaceCatalogView,
	WorkspaceDirectoryRegistrationResult,
	WorkspaceRef,
	WorkspaceRefreshResult,
	WorkspaceSourceSummary,
	WorkspaceSuppressionView,
	WorkspaceView,
} from '@/spec/workspace';

import type { IWorkspaceStoreAPI } from '@/apis/interface';
import { requiredObject, requireWailsBody, wailsObjectArrayOrEmpty } from '@/apis/wailsapi/transport';
import {
	AdoptWorkspaceOccurrence,
	AttachWorkspaceSource,
	CreateEmptyWorkspace,
	CreateFilesystemWorkspace,
	DetachWorkspaceSource,
	GetWorkspace,
	GetWorkspaceArtifact,
	GetWorkspaceCatalog,
	ListWorkspaceArtifacts,
	ListWorkspaces,
	ListWorkspaceSourcesForManagement,
	ListWorkspaceSuppressions,
	PinWorkspaceArtifact,
	PurgeWorkspace,
	PurgeWorkspaceArtifact,
	RefreshWorkspace,
	RegisterWorkspaceDirectory,
	RetireWorkspace,
	SetWorkspaceArtifactEnabled,
	SetWorkspacePrimarySource,
	SetWorkspaceSourceEnabled,
	SuppressWorkspaceBinding,
	UnadoptWorkspaceArtifact,
	UnsuppressWorkspaceBinding,
	UpdateWorkspace,
	UpdateWorkspaceAttachment,
} from '@/apis/wailsjs/go/main/WorkspaceStoreWrapper';

export class WailsWorkspaceStoreAPI implements IWorkspaceStoreAPI {
	async createFilesystemWorkspace(input: CreateFilesystemWorkspaceInput): Promise<WorkspaceView> {
		return requiredObject<WorkspaceView>(
			await CreateFilesystemWorkspace(input as Parameters<typeof CreateFilesystemWorkspace>[0]),
			'CreateFilesystemWorkspace'
		);
	}

	async createEmptyWorkspace(input: CreateEmptyWorkspaceInput): Promise<WorkspaceView> {
		return requiredObject<WorkspaceView>(
			await CreateEmptyWorkspace(input as Parameters<typeof CreateEmptyWorkspace>[0]),
			'CreateEmptyWorkspace'
		);
	}

	async getWorkspace(workspace: WorkspaceRef): Promise<WorkspaceView> {
		const response = await GetWorkspace({
			workspace,
		} as Parameters<typeof GetWorkspace>[0]);
		return requiredObject<WorkspaceView>(response.Body, 'GetWorkspace');
	}

	async listWorkspaces(): Promise<WorkspaceView[]> {
		const response = await ListWorkspaces({} as Parameters<typeof ListWorkspaces>[0]);
		const body = requireWailsBody(response.Body, 'ListWorkspaces');
		return wailsObjectArrayOrEmpty<WorkspaceView>(body.workspaces, 'ListWorkspaces.workspaces');
	}

	async updateWorkspace(workspace: WorkspaceRef, body: UpdateWorkspaceBody): Promise<WorkspaceView> {
		return requiredObject<WorkspaceView>(
			await UpdateWorkspace(
				workspace as Parameters<typeof UpdateWorkspace>[0],
				body as Parameters<typeof UpdateWorkspace>[1]
			),
			'UpdateWorkspace'
		);
	}

	async setWorkspacePrimarySource(
		workspace: WorkspaceRef,
		body: SetWorkspacePrimarySourceBody
	): Promise<WorkspaceView> {
		return requiredObject<WorkspaceView>(
			await SetWorkspacePrimarySource(
				workspace as Parameters<typeof SetWorkspacePrimarySource>[0],
				body as Parameters<typeof SetWorkspacePrimarySource>[1]
			),
			'SetWorkspacePrimarySource'
		);
	}

	async retireWorkspace(workspace: WorkspaceRef, expectedRevision: number): Promise<RetireWorkspaceResult> {
		return requiredObject<RetireWorkspaceResult>(
			await RetireWorkspace(workspace as Parameters<typeof RetireWorkspace>[0], expectedRevision),
			'RetireWorkspace'
		);
	}

	async purgeWorkspace(workspace: WorkspaceRef, expectedRevision: number): Promise<WorkspaceRef> {
		return requiredObject<WorkspaceRef>(
			await PurgeWorkspace(workspace as Parameters<typeof PurgeWorkspace>[0], expectedRevision),
			'PurgeWorkspace'
		);
	}

	async attachWorkspaceSource(workspace: WorkspaceRef, body: AttachWorkspaceSourceBody): Promise<WorkspaceView> {
		return requiredObject<WorkspaceView>(
			await AttachWorkspaceSource(
				workspace as Parameters<typeof AttachWorkspaceSource>[0],
				body as Parameters<typeof AttachWorkspaceSource>[1]
			),
			'AttachWorkspaceSource'
		);
	}

	async updateWorkspaceAttachment(
		workspace: WorkspaceRef,
		sourceID: ArtifactSourceID,
		body: UpdateWorkspaceAttachmentBody
	): Promise<WorkspaceView> {
		return requiredObject<WorkspaceView>(
			await UpdateWorkspaceAttachment(
				workspace as Parameters<typeof UpdateWorkspaceAttachment>[0],
				{
					...body,
					sourceID,
				} as Parameters<typeof UpdateWorkspaceAttachment>[1]
			),
			'UpdateWorkspaceAttachment'
		);
	}

	async detachWorkspaceSource(
		workspace: WorkspaceRef,
		sourceID: ArtifactSourceID,
		body: DetachWorkspaceSourceBody
	): Promise<WorkspaceView> {
		return requiredObject<WorkspaceView>(
			await DetachWorkspaceSource(
				workspace as Parameters<typeof DetachWorkspaceSource>[0],
				{
					...body,
					sourceID,
				} as Parameters<typeof DetachWorkspaceSource>[1]
			),
			'DetachWorkspaceSource'
		);
	}

	async refreshWorkspace(workspace: WorkspaceRef): Promise<WorkspaceRefreshResult> {
		return requiredObject<WorkspaceRefreshResult>(
			await RefreshWorkspace(workspace as Parameters<typeof RefreshWorkspace>[0]),
			'RefreshWorkspace'
		);
	}

	async getWorkspaceCatalog(workspace: WorkspaceRef): Promise<WorkspaceCatalogView> {
		const response = await GetWorkspaceCatalog({
			workspace,
		} as Parameters<typeof GetWorkspaceCatalog>[0]);
		return requiredObject<WorkspaceCatalogView>(response.Body, 'GetWorkspaceCatalog');
	}

	async getWorkspaceArtifact(workspace: WorkspaceRef, artifact: ArtifactRef): Promise<WorkspaceArtifactView> {
		const response = await GetWorkspaceArtifact({
			workspace,
			artifact,
		} as Parameters<typeof GetWorkspaceArtifact>[0]);
		return requiredObject<WorkspaceArtifactView>(response.Body, 'GetWorkspaceArtifact');
	}

	async listWorkspaceArtifacts(workspace: WorkspaceRef): Promise<WorkspaceArtifactView[]> {
		const response = await ListWorkspaceArtifacts({
			workspace,
		} as Parameters<typeof ListWorkspaceArtifacts>[0]);
		const body = requireWailsBody(response.Body, 'ListWorkspaceArtifacts');
		return wailsObjectArrayOrEmpty<WorkspaceArtifactView>(body.artifacts, 'ListWorkspaceArtifacts.artifacts');
	}

	async adoptWorkspaceOccurrence(
		workspace: WorkspaceRef,
		body: AdoptWorkspaceOccurrenceBody
	): Promise<WorkspaceArtifactView> {
		return requiredObject<WorkspaceArtifactView>(
			await AdoptWorkspaceOccurrence(
				workspace as Parameters<typeof AdoptWorkspaceOccurrence>[0],
				body as Parameters<typeof AdoptWorkspaceOccurrence>[1]
			),
			'AdoptWorkspaceOccurrence'
		);
	}

	async pinWorkspaceArtifact(workspace: WorkspaceRef, body: PinWorkspaceArtifactBody): Promise<WorkspaceArtifactView> {
		return requiredObject<WorkspaceArtifactView>(
			await PinWorkspaceArtifact(
				workspace as Parameters<typeof PinWorkspaceArtifact>[0],
				body as Parameters<typeof PinWorkspaceArtifact>[1]
			),
			'PinWorkspaceArtifact'
		);
	}

	async listWorkspaceSuppressions(workspace: WorkspaceRef): Promise<WorkspaceSuppressionView[]> {
		return wailsObjectArrayOrEmpty<WorkspaceSuppressionView>(
			await ListWorkspaceSuppressions(workspace as Parameters<typeof ListWorkspaceSuppressions>[0]),
			'ListWorkspaceSuppressions'
		);
	}

	async suppressWorkspaceBinding(
		workspace: WorkspaceRef,
		body: SuppressWorkspaceBindingBody
	): Promise<WorkspaceSuppressionView> {
		return requiredObject<WorkspaceSuppressionView>(
			await SuppressWorkspaceBinding(
				workspace as Parameters<typeof SuppressWorkspaceBinding>[0],
				body as Parameters<typeof SuppressWorkspaceBinding>[1]
			),
			'SuppressWorkspaceBinding'
		);
	}

	async unsuppressWorkspaceBinding(
		workspace: WorkspaceRef,
		binding: ArtifactSourceBinding,
		expectedRevision: number
	): Promise<UnsuppressWorkspaceBindingResult> {
		return requiredObject<UnsuppressWorkspaceBindingResult>(
			await UnsuppressWorkspaceBinding(
				workspace as Parameters<typeof UnsuppressWorkspaceBinding>[0],
				binding as Parameters<typeof UnsuppressWorkspaceBinding>[1],
				expectedRevision
			),
			'UnsuppressWorkspaceBinding'
		);
	}

	async setWorkspaceArtifactEnabled(
		workspace: WorkspaceRef,
		artifact: ArtifactRef,
		body: SetWorkspaceArtifactEnabledBody
	): Promise<WorkspaceArtifactView> {
		return requiredObject<WorkspaceArtifactView>(
			await SetWorkspaceArtifactEnabled(
				workspace as Parameters<typeof SetWorkspaceArtifactEnabled>[0],
				artifact as Parameters<typeof SetWorkspaceArtifactEnabled>[1],
				body.expectedRevision,
				body.enabled
			),
			'SetWorkspaceArtifactEnabled'
		);
	}

	async unadoptWorkspaceArtifact(
		workspace: WorkspaceRef,
		artifact: ArtifactRef,
		body: UnadoptWorkspaceArtifactBody
	): Promise<UnadoptWorkspaceArtifactResult> {
		return requiredObject<UnadoptWorkspaceArtifactResult>(
			await UnadoptWorkspaceArtifact(
				workspace as Parameters<typeof UnadoptWorkspaceArtifact>[0],
				artifact as Parameters<typeof UnadoptWorkspaceArtifact>[1],
				body as Parameters<typeof UnadoptWorkspaceArtifact>[2]
			),
			'UnadoptWorkspaceArtifact'
		);
	}

	async purgeWorkspaceArtifact(
		workspace: WorkspaceRef,
		artifact: ArtifactRef,
		expectedRevision: number
	): Promise<ArtifactRef> {
		return requiredObject<ArtifactRef>(
			await PurgeWorkspaceArtifact(
				workspace as Parameters<typeof PurgeWorkspaceArtifact>[0],
				artifact as Parameters<typeof PurgeWorkspaceArtifact>[1],
				expectedRevision
			),
			'PurgeWorkspaceArtifact'
		);
	}

	async listWorkspaceSourcesForManagement(): Promise<WorkspaceSourceSummary[]> {
		return wailsObjectArrayOrEmpty<WorkspaceSourceSummary>(
			await ListWorkspaceSourcesForManagement(),
			'ListWorkspaceSourcesForManagement'
		);
	}

	async registerWorkspaceDirectory(
		workspace: WorkspaceRef,
		input: RegisterWorkspaceDirectoryInput
	): Promise<WorkspaceDirectoryRegistrationResult> {
		return requiredObject<WorkspaceDirectoryRegistrationResult>(
			await RegisterWorkspaceDirectory(
				workspace as Parameters<typeof RegisterWorkspaceDirectory>[0],
				input as Parameters<typeof RegisterWorkspaceDirectory>[1]
			),
			'RegisterWorkspaceDirectory'
		);
	}

	async setWorkspaceSourceEnabled(
		sourceID: ArtifactSourceID,
		expectedRevision: number,
		enabled: boolean
	): Promise<WorkspaceSourceSummary> {
		return requiredObject<WorkspaceSourceSummary>(
			await SetWorkspaceSourceEnabled(sourceID, expectedRevision, enabled),
			'SetWorkspaceSourceEnabled'
		);
	}
}
