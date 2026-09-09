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
	SetWorkspaceArtifactRuntimeDisabledBody,
	SetWorkspacePrimarySourceBody,
	SuppressWorkspaceBindingBody,
	UnadoptWorkspaceArtifactBody,
	UnadoptWorkspaceArtifactResult,
	UnsuppressWorkspaceBindingResult,
	UpdateWorkspaceAttachmentBody,
	UpdateWorkspaceBody,
	WorkspaceArtifactView,
	WorkspaceCatalogView,
	WorkspaceContextInspectionView,
	WorkspaceContextLoadPlan,
	WorkspaceContextView,
	WorkspaceDirectoryRegistrationResult,
	WorkspaceRef,
	WorkspaceRefreshResult,
	WorkspaceSkillLoadView,
	WorkspaceSkillView,
	WorkspaceSourceSummary,
	WorkspaceSuppressionView,
	WorkspaceView,
} from '@/spec/workspace';

import type { IWorkspaceAggregateAPI, IWorkspaceRuntimeAPI, IWorkspaceStoreAPI } from '@/apis/interface';

export class WorkspaceManagementAPI {
	constructor(
		// oxlint-disable-next-line typescript/parameter-properties
		private readonly store: IWorkspaceStoreAPI,
		// oxlint-disable-next-line typescript/parameter-properties
		private readonly runtime: IWorkspaceRuntimeAPI,
		// oxlint-disable-next-line typescript/parameter-properties
		private readonly aggregate: IWorkspaceAggregateAPI
	) {}

	createFilesystemWorkspace(input: CreateFilesystemWorkspaceInput): Promise<WorkspaceView> {
		return this.store.createFilesystemWorkspace(input);
	}

	createEmptyWorkspace(input: CreateEmptyWorkspaceInput): Promise<WorkspaceView> {
		return this.store.createEmptyWorkspace(input);
	}

	getWorkspace(workspace: WorkspaceRef): Promise<WorkspaceView> {
		return this.store.getWorkspace(workspace);
	}

	listWorkspaces(): Promise<WorkspaceView[]> {
		return this.store.listWorkspaces();
	}

	updateWorkspace(workspace: WorkspaceRef, body: UpdateWorkspaceBody): Promise<WorkspaceView> {
		return this.store.updateWorkspace(workspace, body);
	}

	setWorkspacePrimarySource(workspace: WorkspaceRef, body: SetWorkspacePrimarySourceBody): Promise<WorkspaceView> {
		return this.store.setWorkspacePrimarySource(workspace, body);
	}

	retireWorkspace(workspace: WorkspaceRef, expectedRevision: number): Promise<RetireWorkspaceResult> {
		return this.store.retireWorkspace(workspace, expectedRevision);
	}

	purgeWorkspace(workspace: WorkspaceRef, expectedRevision: number): Promise<WorkspaceRef> {
		return this.store.purgeWorkspace(workspace, expectedRevision);
	}

	attachWorkspaceSource(workspace: WorkspaceRef, body: AttachWorkspaceSourceBody): Promise<WorkspaceView> {
		return this.store.attachWorkspaceSource(workspace, body);
	}

	updateWorkspaceAttachment(
		workspace: WorkspaceRef,
		sourceID: ArtifactSourceID,
		body: UpdateWorkspaceAttachmentBody
	): Promise<WorkspaceView> {
		return this.store.updateWorkspaceAttachment(workspace, sourceID, body);
	}

	detachWorkspaceSource(
		workspace: WorkspaceRef,
		sourceID: ArtifactSourceID,
		body: DetachWorkspaceSourceBody
	): Promise<WorkspaceView> {
		return this.store.detachWorkspaceSource(workspace, sourceID, body);
	}

	refreshWorkspace(workspace: WorkspaceRef): Promise<WorkspaceRefreshResult> {
		return this.store.refreshWorkspace(workspace);
	}

	getWorkspaceCatalog(workspace: WorkspaceRef): Promise<WorkspaceCatalogView> {
		return this.store.getWorkspaceCatalog(workspace);
	}

	getWorkspaceArtifact(workspace: WorkspaceRef, artifact: ArtifactRef): Promise<WorkspaceArtifactView> {
		return this.store.getWorkspaceArtifact(workspace, artifact);
	}

	listWorkspaceArtifacts(workspace: WorkspaceRef): Promise<WorkspaceArtifactView[]> {
		return this.store.listWorkspaceArtifacts(workspace);
	}

	adoptWorkspaceOccurrence(
		workspace: WorkspaceRef,
		body: AdoptWorkspaceOccurrenceBody
	): Promise<WorkspaceArtifactView> {
		return this.store.adoptWorkspaceOccurrence(workspace, body);
	}

	pinWorkspaceArtifact(workspace: WorkspaceRef, body: PinWorkspaceArtifactBody): Promise<WorkspaceArtifactView> {
		return this.store.pinWorkspaceArtifact(workspace, body);
	}

	listWorkspaceSuppressions(workspace: WorkspaceRef): Promise<WorkspaceSuppressionView[]> {
		return this.store.listWorkspaceSuppressions(workspace);
	}

	suppressWorkspaceBinding(
		workspace: WorkspaceRef,
		body: SuppressWorkspaceBindingBody
	): Promise<WorkspaceSuppressionView> {
		return this.store.suppressWorkspaceBinding(workspace, body);
	}

	unsuppressWorkspaceBinding(
		workspace: WorkspaceRef,
		binding: ArtifactSourceBinding,
		expectedRevision: number
	): Promise<UnsuppressWorkspaceBindingResult> {
		return this.store.unsuppressWorkspaceBinding(workspace, binding, expectedRevision);
	}

	setWorkspaceArtifactEnabled(
		workspace: WorkspaceRef,
		artifact: ArtifactRef,
		body: SetWorkspaceArtifactEnabledBody
	): Promise<WorkspaceArtifactView> {
		return this.store.setWorkspaceArtifactEnabled(workspace, artifact, body);
	}

	unadoptWorkspaceArtifact(
		workspace: WorkspaceRef,
		artifact: ArtifactRef,
		body: UnadoptWorkspaceArtifactBody
	): Promise<UnadoptWorkspaceArtifactResult> {
		return this.store.unadoptWorkspaceArtifact(workspace, artifact, body);
	}

	purgeWorkspaceArtifact(
		workspace: WorkspaceRef,
		artifact: ArtifactRef,
		expectedRevision: number
	): Promise<ArtifactRef> {
		return this.store.purgeWorkspaceArtifact(workspace, artifact, expectedRevision);
	}

	setWorkspaceArtifactRuntimeDisabled(
		workspace: WorkspaceRef,
		artifact: ArtifactRef,
		body: SetWorkspaceArtifactRuntimeDisabledBody
	): Promise<WorkspaceArtifactView> {
		return this.aggregate.setWorkspaceArtifactRuntimeDisabled(workspace, artifact, body);
	}

	listWorkspaceContexts(workspace: WorkspaceRef): Promise<WorkspaceContextView[]> {
		return this.runtime.listWorkspaceContexts(workspace);
	}

	loadWorkspaceContexts(workspace: WorkspaceRef, artifacts?: ArtifactRef[]): Promise<WorkspaceContextInspectionView> {
		return this.runtime.loadWorkspaceContexts(workspace, artifacts);
	}

	composeWorkspaceContext(workspace: WorkspaceRef, artifacts?: ArtifactRef[]): Promise<WorkspaceContextLoadPlan> {
		return this.runtime.composeWorkspaceContext(workspace, artifacts);
	}

	listWorkspaceSkills(workspace: WorkspaceRef): Promise<WorkspaceSkillView[]> {
		return this.runtime.listWorkspaceSkills(workspace);
	}

	loadWorkspaceSkills(workspace: WorkspaceRef, artifacts: ArtifactRef[]): Promise<WorkspaceSkillLoadView> {
		return this.runtime.loadWorkspaceSkills(workspace, artifacts);
	}

	listWorkspaceSourcesForManagement(): Promise<WorkspaceSourceSummary[]> {
		return this.store.listWorkspaceSourcesForManagement();
	}

	registerWorkspaceDirectory(
		workspace: WorkspaceRef,
		input: RegisterWorkspaceDirectoryInput
	): Promise<WorkspaceDirectoryRegistrationResult> {
		return this.store.registerWorkspaceDirectory(workspace, input);
	}

	setWorkspaceSourceEnabled(
		sourceID: ArtifactSourceID,
		expectedRevision: number,
		enabled: boolean
	): Promise<WorkspaceSourceSummary> {
		return this.store.setWorkspaceSourceEnabled(sourceID, expectedRevision, enabled);
	}
}
