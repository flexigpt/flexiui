import {
	type AddProviderRequest,
	DefaultModelParams,
	type ModelDefaults,
	type ModelName,
	type ModelParams,
	type ProviderName,
	type ReasoningParams,
} from '@/models/aiprovidermodel';
import {
	type AISetting,
	type AISettingAttrs,
	type ChatOptions,
	DefaultChatOptions,
	DefaultModelSetting,
	type ModelSetting,
} from '@/models/settingmodel';

import { providerSetAPI, settingstoreAPI } from '@/apis/baseapi';

export async function SetAppSettings(defaultProvider: ProviderName) {
	await settingstoreAPI.setAppSettings(defaultProvider);
	await providerSetAPI.setDefaultProvider(defaultProvider);
}

export async function AddAISetting(providerName: ProviderName, aiSetting: AISetting) {
	// Persist to setting store
	await settingstoreAPI.addAISetting(providerName, aiSetting);
	const req: AddProviderRequest = {
		provider: providerName,
		apiKey: aiSetting.apiKey,
		origin: aiSetting.origin,
		chatCompletionPathPrefix: aiSetting.chatCompletionPathPrefix,
	};
	await providerSetAPI.addProvider(req);
}

/**
 * @public
 */
export async function DeleteAISetting(providerName: ProviderName) {
	await settingstoreAPI.deleteAISetting(providerName);
	await providerSetAPI.deleteProvider(providerName);
}

export async function SetAISettingAPIKey(providerName: ProviderName, apiKey: string) {
	await settingstoreAPI.setAISettingAPIKey(providerName, apiKey);
	await providerSetAPI.setProviderAPIKey(providerName, apiKey);
}

export async function SetAISettingAttrs(providerName: ProviderName, aiSettingAttrs: AISettingAttrs) {
	await settingstoreAPI.setAISettingAttrs(providerName, aiSettingAttrs);
	if (typeof aiSettingAttrs.origin !== 'undefined' || typeof aiSettingAttrs.chatCompletionPathPrefix !== 'undefined') {
		await providerSetAPI.setProviderAttribute(
			providerName,
			aiSettingAttrs.origin,
			aiSettingAttrs.chatCompletionPathPrefix
		);
	}
}

function mergeDefaultsModelSettingAndInbuilt(
	providerName: ProviderName,
	modelName: ModelName,
	inbuiltProviderModels: Record<ProviderName, Record<ModelName, ModelParams>>,
	modelSetting: ModelSetting
): ModelParams {
	let inbuiltModelParams: ModelParams | undefined = undefined;
	if (providerName in inbuiltProviderModels && modelName in inbuiltProviderModels[providerName]) {
		inbuiltModelParams = inbuiltProviderModels[providerName][modelName];
	}

	let temperature = DefaultModelParams.temperature ?? 0.1;
	if (typeof modelSetting.temperature !== 'undefined') {
		temperature = modelSetting.temperature;
	} else if (inbuiltModelParams) {
		temperature = inbuiltModelParams.temperature ?? 0.1;
	}

	let stream = DefaultModelParams.stream;
	if (typeof modelSetting.stream !== 'undefined') {
		stream = modelSetting.stream;
	} else if (inbuiltModelParams) {
		stream = inbuiltModelParams.stream;
	}

	let maxPromptLength = DefaultModelParams.maxPromptLength;
	if (typeof modelSetting.maxPromptLength !== 'undefined') {
		maxPromptLength = modelSetting.maxPromptLength;
	} else if (inbuiltModelParams) {
		maxPromptLength = inbuiltModelParams.maxPromptLength;
	}

	let maxOutputLength = DefaultModelParams.maxOutputLength;
	if (typeof modelSetting.maxOutputLength !== 'undefined') {
		maxOutputLength = modelSetting.maxOutputLength;
	} else if (inbuiltModelParams) {
		maxOutputLength = inbuiltModelParams.maxOutputLength;
	}

	let reasoning: ReasoningParams | undefined = undefined;
	if (typeof modelSetting.reasoning !== 'undefined') {
		reasoning = modelSetting.reasoning;
	} else if (inbuiltModelParams) {
		reasoning = inbuiltModelParams.reasoning;
	}

	let systemPrompt = DefaultModelParams.systemPrompt;
	if (typeof modelSetting.systemPrompt !== 'undefined') {
		systemPrompt = modelSetting.systemPrompt;
	} else if (inbuiltModelParams) {
		systemPrompt = inbuiltModelParams.systemPrompt;
	}

	let timeout = DefaultModelParams.timeout;
	if (typeof modelSetting.timeout !== 'undefined') {
		timeout = modelSetting.timeout;
	} else if (inbuiltModelParams) {
		timeout = inbuiltModelParams.timeout;
	}

	let additionalParameters = { ...DefaultModelParams.additionalParameters };
	if (modelSetting.additionalParameters) {
		additionalParameters = { ...additionalParameters, ...modelSetting.additionalParameters };
	}

	return {
		name: modelName,
		stream: stream,
		maxPromptLength: maxPromptLength,
		maxOutputLength: maxOutputLength,
		temperature: temperature,
		reasoning: reasoning,
		systemPrompt: systemPrompt,
		timeout: timeout,
		additionalParameters: additionalParameters,
	};
}

export async function GetChatInputOptions(): Promise<{ allOptions: ChatOptions[]; default: ChatOptions }> {
	try {
		// Fetch configuration info and settings
		const info = await providerSetAPI.getConfigurationInfo();
		if (info.defaultProvider === '' || Object.keys(info.configuredProviders).length === 0) {
			return { allOptions: [DefaultChatOptions], default: DefaultChatOptions };
		}
		const configDefaultProvider = info.defaultProvider;
		const providerInfoDict = info.configuredProviders;
		const inbuiltProviderModels = info.inbuiltProviderModels;

		const onlySettings = await settingstoreAPI.getAllSettings();
		const mergedSettings = MergeInbuiltModelsWithSettings(
			onlySettings.aiSettings,
			inbuiltProviderModels,
			info.inbuiltProviderModelDefaults
		);
		// Initialize default option and input models array
		let defaultOption: ChatOptions | undefined;
		const inputModels: ChatOptions[] = [];

		for (const providerName of Object.keys(providerInfoDict)) {
			const aiSetting = mergedSettings[providerName];
			if (aiSetting.isEnabled) {
				const settingsDefaultModelName = aiSetting.defaultModel;
				for (const [modelName, modelSetting] of Object.entries(aiSetting.modelSettings)) {
					if (!modelSetting.isEnabled) {
						continue;
					}
					const mergedModelParam = mergeDefaultsModelSettingAndInbuilt(
						providerName,
						modelName,
						inbuiltProviderModels,
						modelSetting
					);
					const chatOption: ChatOptions = {
						title: modelSetting.displayName,
						provider: providerName,
						name: modelName,
						stream: mergedModelParam.stream,
						maxPromptLength: mergedModelParam.maxPromptLength,
						maxOutputLength: mergedModelParam.maxOutputLength,
						temperature: mergedModelParam.temperature,
						reasoning: mergedModelParam.reasoning,
						systemPrompt: mergedModelParam.systemPrompt,
						timeout: mergedModelParam.timeout,
						additionalParameters: mergedModelParam.additionalParameters,
						disablePreviousMessages: false,
					};
					inputModels.push(chatOption);
					if (modelName === settingsDefaultModelName && providerName === configDefaultProvider) {
						defaultOption = chatOption;
					}
				}
			}
		}

		if (defaultOption === undefined || inputModels.length === 0) {
			throw Error('No input option found !!!');
		}
		return { allOptions: inputModels, default: defaultOption };
	} catch (error) {
		console.error('Error fetching chat input options:', error);
		return { allOptions: [DefaultChatOptions], default: DefaultChatOptions };
	}
}

// Create a Model setting from default or inbuilt ModelParams if available
export async function PopulateModelSettingDefaults(
	providerName: ProviderName,
	modelName: ModelName,
	existingData?: ModelSetting
): Promise<ModelSetting> {
	const info = await providerSetAPI.getConfigurationInfo();
	if (
		info.defaultProvider === '' ||
		Object.keys(info.configuredProviders).length === 0 ||
		!(providerName in info.configuredProviders)
	) {
		return DefaultModelSetting;
	}

	let modelSetting = existingData;
	if (typeof modelSetting === 'undefined') {
		modelSetting = {
			displayName: modelName,
			isEnabled: true,
		};
	}

	const mergedModelParam = mergeDefaultsModelSettingAndInbuilt(
		providerName,
		modelName,
		info.inbuiltProviderModels,
		modelSetting
	);

	return {
		displayName: modelSetting.displayName,
		isEnabled: modelSetting.isEnabled,
		stream: mergedModelParam.stream,
		maxPromptLength: mergedModelParam.maxPromptLength,
		maxOutputLength: mergedModelParam.maxOutputLength,
		temperature: mergedModelParam.temperature,
		reasoning: mergedModelParam.reasoning,
		systemPrompt: mergedModelParam.systemPrompt,
		timeout: mergedModelParam.timeout,
	};
}

export function MergeInbuiltModelsWithSettings(
	aiSettings: Record<ProviderName, AISetting>,
	inbuiltProviderModels: Record<ProviderName, Record<ModelName, ModelParams>>,
	inbuiltProviderModelDefaults: Record<ProviderName, Record<ModelName, ModelDefaults>>
): Record<ProviderName, AISetting> {
	const newSettings = { ...aiSettings };

	for (const provider in inbuiltProviderModels) {
		// If provider not in aiSettings, add it
		if (!(provider in newSettings)) {
			continue;
		}

		// For each model in inbuiltProviderInfo, ensure it exists in modelSettings
		for (const model in inbuiltProviderModels[provider]) {
			if (model in newSettings[provider].modelSettings) {
				continue;
			}

			newSettings[provider].modelSettings[model] = {
				displayName: model,
				isEnabled: false,
				stream: inbuiltProviderModels[provider][model].stream,
				maxPromptLength: inbuiltProviderModels[provider][model].maxPromptLength,
				maxOutputLength: inbuiltProviderModels[provider][model].maxOutputLength,
				temperature: inbuiltProviderModels[provider][model].temperature,
				reasoning: inbuiltProviderModels[provider][model].reasoning,
				systemPrompt: inbuiltProviderModels[provider][model].systemPrompt,
				timeout: inbuiltProviderModels[provider][model].timeout,
				additionalParameters: inbuiltProviderModels[provider][model].additionalParameters,
			};
			if (provider in inbuiltProviderModelDefaults && model in inbuiltProviderModelDefaults[provider]) {
				newSettings[provider].modelSettings[model].displayName =
					inbuiltProviderModelDefaults[provider][model].displayName;
				newSettings[provider].modelSettings[model].isEnabled = inbuiltProviderModelDefaults[provider][model].isEnabled;
			}
		}
	}

	return newSettings;
}
