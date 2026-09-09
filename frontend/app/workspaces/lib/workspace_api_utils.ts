import type { ArtifactRef } from '@/spec/artifact';
import type {
	CreateEmptyWorkspaceInput,
	CreateFilesystemWorkspaceInput,
	WorkspaceRef,
	WorkspaceView,
} from '@/spec/workspace';

import { workspaceManagementAPI } from '@/apis/baseapi';

import { sortWorkspaces } from '@/workspaces/lib/workspace_utils';

export type { CreateEmptyWorkspaceInput, CreateFilesystemWorkspaceInput };

export function workspaceRefKey(workspace: WorkspaceRef): string {
	return `${workspace.rootID}:${workspace.collectionID}`;
}

export function artifactRefKey(artifact: ArtifactRef): string {
	return `${artifact.rootID}:${artifact.artifactID}`;
}

export function workspaceRefsEqual(
	left: WorkspaceRef | null | undefined,
	right: WorkspaceRef | null | undefined
): boolean {
	return left?.rootID === right?.rootID && left?.collectionID === right?.collectionID;
}

export async function listAllWorkspaces(): Promise<WorkspaceView[]> {
	return sortWorkspaces(await workspaceManagementAPI.listWorkspaces());
}

export async function createFilesystemWorkspaceCollection(
	body: CreateFilesystemWorkspaceInput
): Promise<WorkspaceView> {
	return workspaceManagementAPI.createFilesystemWorkspace(body);
}

export async function createEmptyWorkspaceCollection(body: CreateEmptyWorkspaceInput): Promise<WorkspaceView> {
	return workspaceManagementAPI.createEmptyWorkspace(body);
}
