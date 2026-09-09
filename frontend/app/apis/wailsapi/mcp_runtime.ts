import type {
	InvokeMCPToolRequestBody,
	InvokeMCPToolResponseBody,
	MCPApprovalEvaluation,
	MCPApprovalResolution,
	MCPApprovalResolutionResult,
	MCPCompletionRefType,
	MCPCompletionResult,
	MCPGetPromptResponseBody,
	MCPGlobalSettings,
	MCPOAuthAuthorization,
	MCPPromptRef,
	MCPProviderToolMapping,
	MCPReadResourceResponseBody,
	MCPResourceRef,
	MCPResourceTemplateRef,
	MCPRuntimeServerID,
	MCPServerRuntimeSnapshot,
	MCPToolCapability,
} from '@/spec/mcp_artifact';

import type { IMCPRuntimeAPI } from '@/apis/interface';
import {
	requiredObject,
	requireWailsBoolean,
	requireWailsFiniteNumber,
	wailsObjectArrayOrEmpty,
} from '@/apis/wailsapi/transport';
import {
	CancelPendingMCPOAuthAuthorization,
	CompleteMCPArgument,
	DisconnectMCPServer,
	EvaluateMappedMCPToolCall,
	EvaluateMCPToolCall,
	GetMCPGlobalSettings,
	GetMCPPrompt,
	GetMCPServerStatus,
	InvokeMappedMCPTool,
	InvokeMCPTool,
	ListMCPServerPrompts,
	ListMCPServerResources,
	ListMCPServerResourceTemplates,
	ListMCPServerTools,
	ListPendingMCPOAuthAuthorizations,
	ReadMCPResource,
	RefreshMCPServer,
	ResolveMCPApproval,
	StartMCPServerConnect,
	UpdateMCPGlobalSettings,
} from '@/apis/wailsjs/go/main/MCPRuntimeWrapper';

export class WailsMCPRuntimeAPI implements IMCPRuntimeAPI {
	async connectMCPServer(server: MCPRuntimeServerID): Promise<MCPServerRuntimeSnapshot> {
		return requiredObject<MCPServerRuntimeSnapshot>(
			await StartMCPServerConnect(server as Parameters<typeof StartMCPServerConnect>[0]),
			'StartMCPServerConnect'
		);
	}

	async disconnectMCPServer(server: MCPRuntimeServerID): Promise<void> {
		await DisconnectMCPServer(server as Parameters<typeof DisconnectMCPServer>[0]);
	}

	async refreshMCPServer(server: MCPRuntimeServerID): Promise<MCPServerRuntimeSnapshot> {
		return requiredObject<MCPServerRuntimeSnapshot>(
			await RefreshMCPServer(server as Parameters<typeof RefreshMCPServer>[0]),
			'RefreshMCPServer'
		);
	}

	async getMCPServerStatus(server: MCPRuntimeServerID): Promise<MCPServerRuntimeSnapshot> {
		return requiredObject<MCPServerRuntimeSnapshot>(
			await GetMCPServerStatus(server as Parameters<typeof GetMCPServerStatus>[0]),
			'GetMCPServerStatus'
		);
	}

	async listMCPServerTools(server: MCPRuntimeServerID): Promise<MCPToolCapability[]> {
		return wailsObjectArrayOrEmpty<MCPToolCapability>(
			await ListMCPServerTools(server as Parameters<typeof ListMCPServerTools>[0]),
			'ListMCPServerTools'
		);
	}

	async listMCPServerResources(server: MCPRuntimeServerID): Promise<MCPResourceRef[]> {
		return wailsObjectArrayOrEmpty<MCPResourceRef>(
			await ListMCPServerResources(server as Parameters<typeof ListMCPServerResources>[0]),
			'ListMCPServerResources'
		);
	}

	async listMCPServerResourceTemplates(server: MCPRuntimeServerID): Promise<MCPResourceTemplateRef[]> {
		return wailsObjectArrayOrEmpty<MCPResourceTemplateRef>(
			await ListMCPServerResourceTemplates(server as Parameters<typeof ListMCPServerResourceTemplates>[0]),
			'ListMCPServerResourceTemplates'
		);
	}

	async listMCPServerPrompts(server: MCPRuntimeServerID): Promise<MCPPromptRef[]> {
		return wailsObjectArrayOrEmpty<MCPPromptRef>(
			await ListMCPServerPrompts(server as Parameters<typeof ListMCPServerPrompts>[0]),
			'ListMCPServerPrompts'
		);
	}

	async readMCPResource(server: MCPRuntimeServerID, uri: string): Promise<MCPReadResourceResponseBody> {
		return requiredObject<MCPReadResourceResponseBody>(
			await ReadMCPResource(server as Parameters<typeof ReadMCPResource>[0], uri),
			'ReadMCPResource'
		);
	}

	async getMCPPrompt(
		server: MCPRuntimeServerID,
		promptName: string,
		promptArguments?: Record<string, string>
	): Promise<MCPGetPromptResponseBody> {
		return requiredObject<MCPGetPromptResponseBody>(
			await GetMCPPrompt(server as Parameters<typeof GetMCPPrompt>[0], promptName, promptArguments ?? {}),
			'GetMCPPrompt'
		);
	}

	async completeMCPArgument(
		server: MCPRuntimeServerID,
		refType: MCPCompletionRefType,
		name: string,
		argumentName: string,
		argumentValue?: string,
		context?: Record<string, string>
	): Promise<MCPCompletionResult> {
		return requiredObject<MCPCompletionResult>(
			await CompleteMCPArgument(
				server as Parameters<typeof CompleteMCPArgument>[0],
				{
					refType,
					name,
					argumentName,
					argumentValue,
					context,
				} as Parameters<typeof CompleteMCPArgument>[1]
			),
			'CompleteMCPArgument'
		);
	}

	async evaluateMCPToolCall(
		server: MCPRuntimeServerID,
		request: InvokeMCPToolRequestBody
	): Promise<MCPApprovalEvaluation> {
		return requiredObject<MCPApprovalEvaluation>(
			await EvaluateMCPToolCall(
				server as Parameters<typeof EvaluateMCPToolCall>[0],
				request as Parameters<typeof EvaluateMCPToolCall>[1]
			),
			'EvaluateMCPToolCall'
		);
	}

	async evaluateMappedMCPToolCall(
		mapping: MCPProviderToolMapping,
		request: InvokeMCPToolRequestBody
	): Promise<MCPApprovalEvaluation> {
		return requiredObject<MCPApprovalEvaluation>(
			await EvaluateMappedMCPToolCall(
				mapping as Parameters<typeof EvaluateMappedMCPToolCall>[0],
				request as Parameters<typeof EvaluateMappedMCPToolCall>[1]
			),
			'EvaluateMappedMCPToolCall'
		);
	}

	async invokeMCPTool(
		server: MCPRuntimeServerID,
		request: InvokeMCPToolRequestBody
	): Promise<InvokeMCPToolResponseBody> {
		return requiredObject<InvokeMCPToolResponseBody>(
			await InvokeMCPTool(
				server as Parameters<typeof InvokeMCPTool>[0],
				request as Parameters<typeof InvokeMCPTool>[1]
			),
			'InvokeMCPTool'
		);
	}

	async invokeMappedMCPTool(
		mapping: MCPProviderToolMapping,
		request: InvokeMCPToolRequestBody
	): Promise<InvokeMCPToolResponseBody> {
		return requiredObject<InvokeMCPToolResponseBody>(
			await InvokeMappedMCPTool(
				mapping as Parameters<typeof InvokeMappedMCPTool>[0],
				request as Parameters<typeof InvokeMappedMCPTool>[1]
			),
			'InvokeMappedMCPTool'
		);
	}

	async resolveMCPApproval(
		approvalID: string,
		resolution: MCPApprovalResolution
	): Promise<MCPApprovalResolutionResult> {
		return requiredObject<MCPApprovalResolutionResult>(
			await ResolveMCPApproval(approvalID, resolution),
			'ResolveMCPApproval'
		);
	}

	async listPendingMCPOAuthAuthorizations(): Promise<MCPOAuthAuthorization[]> {
		return wailsObjectArrayOrEmpty<MCPOAuthAuthorization>(
			await ListPendingMCPOAuthAuthorizations(),
			'ListPendingMCPOAuthAuthorizations'
		);
	}

	async cancelPendingMCPOAuthAuthorization(server: MCPRuntimeServerID): Promise<boolean> {
		return requireWailsBoolean(
			await CancelPendingMCPOAuthAuthorization(server as Parameters<typeof CancelPendingMCPOAuthAuthorization>[0]),
			'CancelPendingMCPOAuthAuthorization'
		);
	}

	async getMCPGlobalSettings(): Promise<MCPGlobalSettings> {
		return requiredObject<MCPGlobalSettings>(await GetMCPGlobalSettings(), 'GetMCPGlobalSettings');
	}

	async updateMCPGlobalSettings(expectedRevision: number, oauthLoopbackListenAddr?: string): Promise<number> {
		return requireWailsFiniteNumber(
			await UpdateMCPGlobalSettings(expectedRevision, {
				oauthLoopbackListenAddr,
			} as Parameters<typeof UpdateMCPGlobalSettings>[1]),
			'UpdateMCPGlobalSettings'
		);
	}
}
