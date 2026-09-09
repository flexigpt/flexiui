import type { ArtifactRef } from '@/spec/artifact';
import type {
	WorkspaceContextInspectionView,
	WorkspaceContextLoadPlan,
	WorkspaceContextView,
	WorkspaceRef,
	WorkspaceSkillLoadView,
	WorkspaceSkillView,
} from '@/spec/workspace';

import type { IWorkspaceRuntimeAPI } from '@/apis/interface';
import { requireWailsBody, wailsObjectArrayOrEmpty } from '@/apis/wailsapi/transport';
import {
	ComposeWorkspaceContext,
	ListWorkspaceContexts,
	ListWorkspaceSkills,
	LoadWorkspaceContexts,
	LoadWorkspaceSkills,
} from '@/apis/wailsjs/go/main/WorkspaceRuntimeWrapper';

export class WailsWorkspaceRuntimeAPI implements IWorkspaceRuntimeAPI {
	async listWorkspaceContexts(workspace: WorkspaceRef): Promise<WorkspaceContextView[]> {
		const response = await ListWorkspaceContexts({
			workspace,
		} as Parameters<typeof ListWorkspaceContexts>[0]);
		const body = requireWailsBody(response.Body, 'ListWorkspaceContexts');
		return wailsObjectArrayOrEmpty<WorkspaceContextView>(body.contexts, 'ListWorkspaceContexts.contexts');
	}

	async loadWorkspaceContexts(
		workspace: WorkspaceRef,
		artifacts?: ArtifactRef[]
	): Promise<WorkspaceContextInspectionView> {
		const response = await LoadWorkspaceContexts({
			workspace,
			Body: {
				artifacts,
			},
		} as Parameters<typeof LoadWorkspaceContexts>[0]);
		return requireWailsBody(response.Body, 'LoadWorkspaceContexts') as WorkspaceContextInspectionView;
	}

	async composeWorkspaceContext(workspace: WorkspaceRef, artifacts?: ArtifactRef[]): Promise<WorkspaceContextLoadPlan> {
		const response = await ComposeWorkspaceContext({
			workspace,
			Body: {
				artifacts,
			},
		} as Parameters<typeof ComposeWorkspaceContext>[0]);
		return requireWailsBody(response.Body, 'ComposeWorkspaceContext') as WorkspaceContextLoadPlan;
	}

	async listWorkspaceSkills(workspace: WorkspaceRef): Promise<WorkspaceSkillView[]> {
		const response = await ListWorkspaceSkills({
			workspace,
		} as Parameters<typeof ListWorkspaceSkills>[0]);
		const body = requireWailsBody(response.Body, 'ListWorkspaceSkills');
		return wailsObjectArrayOrEmpty<WorkspaceSkillView>(body.skills, 'ListWorkspaceSkills.skills');
	}

	async loadWorkspaceSkills(workspace: WorkspaceRef, artifacts: ArtifactRef[]): Promise<WorkspaceSkillLoadView> {
		const response = await LoadWorkspaceSkills({
			workspace,
			Body: {
				artifacts,
			},
		} as Parameters<typeof LoadWorkspaceSkills>[0]);
		return requireWailsBody(response.Body, 'LoadWorkspaceSkills') as WorkspaceSkillLoadView;
	}
}
