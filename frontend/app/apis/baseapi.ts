// oxlint-disable import/no-mutable-exports
import { IS_WAILS_PLATFORM } from '@/lib/features';
import { setFrontendErrorLogger } from '@/lib/frontend_error_reporter';

import type {
	IAggregateAPI,
	IAssistantPresetStoreAPI,
	IAttachmentsDropAPI,
	IBackendAPI,
	IConversationStoreAPI,
	ILogger,
	IMCPAggregateAPI,
	IMCPRuntimeAPI,
	IMCPStoreAPI,
	IModelPresetStoreAPI,
	ISettingStoreAPI,
	ISkillAggregateAPI,
	ISkillRuntimeAPI,
	ISkillStoreAPI,
	IToolRuntimeAPI,
	IToolStoreAPI,
	IWorkspaceAggregateAPI,
	IWorkspaceRuntimeAPI,
	IWorkspaceStoreAPI,
} from '@/apis/interface';
import { MCPManagementAPI } from '@/apis/mcp_management';
import { SkillManagementAPI } from '@/apis/skill_management';
// oxlint-disable-next-line import/no-namespace
import * as wailsImpl from '@/apis/wailsapi';
import { WailsMCPAggregateAPI } from '@/apis/wailsapi/mcp_aggregate';
import { WailsMCPRuntimeAPI } from '@/apis/wailsapi/mcp_runtime';
import { WailsMCPStoreAPI } from '@/apis/wailsapi/mcp_store';
import { WailsSkillAggregateAPI } from '@/apis/wailsapi/skill_aggregate';
import { WailsSkillRuntimeAPI } from '@/apis/wailsapi/skill_runtime';
import { WailsSkillStoreAPI } from '@/apis/wailsapi/skill_store';
import { WailsWorkspaceAggregateAPI } from '@/apis/wailsapi/workspace_aggregate';
import { WailsWorkspaceRuntimeAPI } from '@/apis/wailsapi/workspace_runtime';
import { WailsWorkspaceStoreAPI } from '@/apis/wailsapi/workspace_store';
import { WorkspaceManagementAPI } from '@/apis/workspace_management';

export let log: ILogger;

export let assistantPresetStoreAPI: IAssistantPresetStoreAPI;
export let attachmentsDropAPI: IAttachmentsDropAPI;
export let backendAPI: IBackendAPI;
export let conversationStoreAPI: IConversationStoreAPI;
export let aggregateAPI: IAggregateAPI;
export let settingstoreAPI: ISettingStoreAPI;
export let modelPresetStoreAPI: IModelPresetStoreAPI;

let mcpStoreAPI: IMCPStoreAPI;
let mcpAggregateAPI: IMCPAggregateAPI;
export let mcpRuntimeAPI: IMCPRuntimeAPI;
export let mcpManagementAPI: MCPManagementAPI;

export let toolStoreAPI: IToolStoreAPI;
export let toolRuntimeAPI: IToolRuntimeAPI;

let skillStoreAPI: ISkillStoreAPI;
let skillAggregateAPI: ISkillAggregateAPI;
let skillRuntimeAPI: ISkillRuntimeAPI;
export let skillManagementAPI: SkillManagementAPI;

let workspaceStoreAPI: IWorkspaceStoreAPI;
let workspaceRuntimeAPI: IWorkspaceRuntimeAPI;
let workspaceAggregateAPI: IWorkspaceAggregateAPI;
export let workspaceManagementAPI: WorkspaceManagementAPI;

// Conditional initialization
if (IS_WAILS_PLATFORM) {
	// Initialize with Wails implementations
	log = new wailsImpl.WailsLogger();
	setFrontendErrorLogger(log);

	attachmentsDropAPI = new wailsImpl.WailsAttachmentsDropAPI();
	backendAPI = new wailsImpl.WailsBackendAPI();
	conversationStoreAPI = new wailsImpl.WailsConversationStoreAPI();
	aggregateAPI = new wailsImpl.WailsAggregateAPI();
	settingstoreAPI = new wailsImpl.WailsSettingStoreAPI();
	modelPresetStoreAPI = new wailsImpl.WailsModelPresetStoreAPI();
	mcpStoreAPI = new WailsMCPStoreAPI();
	mcpRuntimeAPI = new WailsMCPRuntimeAPI();
	mcpAggregateAPI = new WailsMCPAggregateAPI();
	mcpManagementAPI = new MCPManagementAPI(mcpStoreAPI, mcpAggregateAPI);
	toolStoreAPI = new wailsImpl.WailsToolStoreAPI();
	toolRuntimeAPI = new wailsImpl.WailsToolRuntimeAPI();
	assistantPresetStoreAPI = new wailsImpl.WailsAssistantPresetStoreAPI();
	workspaceStoreAPI = new WailsWorkspaceStoreAPI();
	workspaceRuntimeAPI = new WailsWorkspaceRuntimeAPI();
	workspaceAggregateAPI = new WailsWorkspaceAggregateAPI();
	workspaceManagementAPI = new WorkspaceManagementAPI(workspaceStoreAPI, workspaceRuntimeAPI, workspaceAggregateAPI);
	skillStoreAPI = new WailsSkillStoreAPI();
	skillAggregateAPI = new WailsSkillAggregateAPI();
	skillRuntimeAPI = new WailsSkillRuntimeAPI();
	skillManagementAPI = new SkillManagementAPI(skillStoreAPI, skillAggregateAPI, skillRuntimeAPI);
} else {
	// Error for unsupported platforms
	throw new Error('Unsupported platform');
}
