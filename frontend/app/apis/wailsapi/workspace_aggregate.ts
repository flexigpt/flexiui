import type { ArtifactRef } from '@/spec/artifact';
import type { SetWorkspaceArtifactRuntimeDisabledBody, WorkspaceArtifactView, WorkspaceRef } from '@/spec/workspace';

import type { IWorkspaceAggregateAPI } from '@/apis/interface';
import { requireWailsBody } from '@/apis/wailsapi/transport';
import { SetWorkspaceArtifactRuntimeDisabled } from '@/apis/wailsjs/go/main/WorkspaceAggregateWrapper';

export class WailsWorkspaceAggregateAPI implements IWorkspaceAggregateAPI {
	async setWorkspaceArtifactRuntimeDisabled(
		workspace: WorkspaceRef,
		artifact: ArtifactRef,
		body: SetWorkspaceArtifactRuntimeDisabledBody
	): Promise<WorkspaceArtifactView> {
		const response = await SetWorkspaceArtifactRuntimeDisabled({
			workspace,
			artifact,
			Body: body,
		} as Parameters<typeof SetWorkspaceArtifactRuntimeDisabled>[0]);

		return requireWailsBody(response.Body, 'SetWorkspaceArtifactRuntimeDisabled') as WorkspaceArtifactView;
	}
}
