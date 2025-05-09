package consts

import "github.com/flexigpt/flexiui/pkg/aiprovider/spec"

const (
	OpenAICompatibleAPIKeyHeaderKey          = "Authorization"
	OpenAICompatibleChatCompletionPathPrefix = "/v1/chat/completions"
)

var OpenAICompatibleDefaultHeaders = map[string]string{"content-type": "application/json"}

var InbuiltProviders = map[spec.ProviderName]spec.ProviderInfo{
	ProviderNameAnthropic:   AnthropicProviderInfo,
	ProviderNameDeepseek:    DeepseekProviderInfo,
	ProviderNameGoogle:      GoogleProviderInfo,
	ProviderNameHuggingFace: HuggingfaceProviderInfo,
	ProviderNameLlamaCPP:    LlamacppProviderInfo,
	ProviderNameOpenAI:      OpenAIProviderInfo,
}

var InbuiltProviderModels = map[spec.ProviderName]map[spec.ModelName]spec.ModelParams{
	ProviderNameAnthropic:   AnthropicModels,
	ProviderNameDeepseek:    DeepseekModels,
	ProviderNameGoogle:      GoogleModels,
	ProviderNameHuggingFace: HuggingfaceModels,
	ProviderNameLlamaCPP:    LlamacppModels,
	ProviderNameOpenAI:      OpenAIModels,
}
