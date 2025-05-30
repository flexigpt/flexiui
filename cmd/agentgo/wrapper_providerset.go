package main

import (
	"context"
	"errors"

	"github.com/ppipada/flexigpt-app/pkg/aiprovider"
	aiproviderSpec "github.com/ppipada/flexigpt-app/pkg/aiprovider/spec"
	"github.com/ppipada/flexigpt-app/pkg/middleware"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ProviderSetWrapper struct {
	providersetAPI *aiprovider.ProviderSetAPI
	appContext     context.Context
}

// NewProviderSetWrapper creates a new ProviderSet with the specified default provider.
func InitProviderSetWrapper(
	ps *ProviderSetWrapper,
	defaultInbuiltProvider aiproviderSpec.ProviderName,
) error {
	p, err := aiprovider.NewProviderSetAPI(defaultInbuiltProvider, false)
	if err != nil {
		return errors.Join(err, errors.New("invalid default provider"))
	}
	ps.providersetAPI = p
	return nil
}

func SetWrappedProviderAppContext(w *ProviderSetWrapper, ctx context.Context) {
	w.appContext = ctx
}

func (w *ProviderSetWrapper) SetDefaultProvider(
	req *aiproviderSpec.SetDefaultProviderRequest,
) (*aiproviderSpec.SetDefaultProviderResponse, error) {
	return middleware.WithRecoveryResp(func() (*aiproviderSpec.SetDefaultProviderResponse, error) {
		return w.providersetAPI.SetDefaultProvider(context.Background(), req)
	})
}

func (w *ProviderSetWrapper) GetConfigurationInfo(
	req *aiproviderSpec.GetConfigurationInfoRequest,
) (*aiproviderSpec.GetConfigurationInfoResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*aiproviderSpec.GetConfigurationInfoResponse, error) {
			return w.providersetAPI.GetConfigurationInfo(context.Background(), req)
		},
	)
}

func (w *ProviderSetWrapper) AddProvider(
	req *aiproviderSpec.AddProviderRequest,
) (*aiproviderSpec.AddProviderResponse, error) {
	return middleware.WithRecoveryResp(func() (*aiproviderSpec.AddProviderResponse, error) {
		return w.providersetAPI.AddProvider(context.Background(), req)
	})
}

func (w *ProviderSetWrapper) DeleteProvider(
	req *aiproviderSpec.DeleteProviderRequest,
) (*aiproviderSpec.DeleteProviderResponse, error) {
	return middleware.WithRecoveryResp(func() (*aiproviderSpec.DeleteProviderResponse, error) {
		return w.providersetAPI.DeleteProvider(context.Background(), req)
	})
}

func (w *ProviderSetWrapper) SetProviderAPIKey(
	req *aiproviderSpec.SetProviderAPIKeyRequest,
) (*aiproviderSpec.SetProviderAPIKeyResponse, error) {
	return middleware.WithRecoveryResp(func() (*aiproviderSpec.SetProviderAPIKeyResponse, error) {
		return w.providersetAPI.SetProviderAPIKey(context.Background(), req)
	})
}

func (w *ProviderSetWrapper) SetProviderAttribute(
	req *aiproviderSpec.SetProviderAttributeRequest,
) (*aiproviderSpec.SetProviderAttributeResponse, error) {
	return middleware.WithRecoveryResp(
		func() (*aiproviderSpec.SetProviderAttributeResponse, error) {
			return w.providersetAPI.SetProviderAttribute(context.Background(), req)
		},
	)
}

// FetchCompletion handles the completion request and streams data back to the frontend.
func (w *ProviderSetWrapper) FetchCompletion(
	provider string,
	prompt string,
	modelParams aiproviderSpec.ModelParams,
	prevMessages []aiproviderSpec.ChatCompletionRequestMessage,
	callbackID string,
) (*aiproviderSpec.FetchCompletionResponse, error) {
	return middleware.WithRecoveryResp(func() (*aiproviderSpec.FetchCompletionResponse, error) {
		onStreamData := func(data string) error {
			runtime.EventsEmit(w.appContext, callbackID, data)
			return nil
		}

		req := &aiproviderSpec.FetchCompletionRequest{
			Body: &aiproviderSpec.FetchCompletionRequestBody{
				Provider:     aiproviderSpec.ProviderName(provider),
				Prompt:       prompt,
				ModelParams:  modelParams,
				PrevMessages: prevMessages,
				OnStreamData: onStreamData,
			},
		}
		resp, err := w.providersetAPI.FetchCompletion(
			context.Background(),
			req,
		)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})
}
