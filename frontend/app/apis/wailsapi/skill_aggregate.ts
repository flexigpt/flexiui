import type { ArtifactRef } from '@/spec/artifact';
import type { ResolvedSkillRuntime, SkillBundleRef, SkillRuntimeCatalogID } from '@/spec/skill';

import type { ISkillAggregateAPI } from '@/apis/interface';
import { requireNonBlankString, requireWailsBody } from '@/apis/wailsapi/transport';
import { ResolveArtifactSkill, RuntimeCatalogIDForCollection } from '@/apis/wailsjs/go/main/SkillAggregateWrapper';

export class WailsSkillAggregateAPI implements ISkillAggregateAPI {
	async runtimeCatalogIDForCollection(bundle: SkillBundleRef): Promise<SkillRuntimeCatalogID> {
		const catalogID = await RuntimeCatalogIDForCollection(
			bundle as Parameters<typeof RuntimeCatalogIDForCollection>[0]
		);

		return requireNonBlankString(catalogID, 'RuntimeCatalogIDForCollection');
	}

	async resolveArtifactSkill(artifact: ArtifactRef): Promise<ResolvedSkillRuntime> {
		return requireWailsBody(
			await ResolveArtifactSkill(artifact as Parameters<typeof ResolveArtifactSkill>[0]),
			'ResolveArtifactSkill'
		) as ResolvedSkillRuntime;
	}
}
