package bundle

import (
	"context"
	"errors"
)

// StoreAPI is the domain-specific request/response facade for Skill Bundles.
//
// It does not expose Artifact Store as a frontend API. It accepts Skill Bundle
// terms, performs request-shape validation, and delegates domain semantics to
// the Skill Bundle service.
type StoreAPI struct {
	service *API
}

func NewStoreAPI(service *API) (*StoreAPI, error) {
	if service == nil {
		return nil, wrapStoreError(
			"facade initialization",
			errors.New("skill bundle service is required"),
		)
	}
	return &StoreAPI{service: service}, nil
}

func (a *StoreAPI) CreateSkillBundle(
	ctx context.Context,
	request *CreateSkillBundleRequest,
) (*CreateSkillBundleResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Skill Bundle creation",
	); err != nil {
		return nil, err
	}

	value, err := a.service.CreateBundle(ctx, *request.Body)
	if err != nil {
		return nil, wrapStoreError("create", err)
	}
	return &CreateSkillBundleResponse{Body: &value}, nil
}

func (a *StoreAPI) GetSkillBundle(
	ctx context.Context,
	request *GetSkillBundleRequest,
) (*GetSkillBundleResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Skill Bundle get",
	); err != nil {
		return nil, err
	}

	value, err := a.service.GetBundle(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("get", err)
	}
	return &GetSkillBundleResponse{Body: &value}, nil
}

func (a *StoreAPI) ListSkillBundles(
	ctx context.Context,
	request *ListSkillBundlesRequest,
) (*ListSkillBundlesResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Skill Bundle list",
	); err != nil {
		return nil, err
	}

	values, err := a.service.ListBundles(ctx, request.RootID)
	if err != nil {
		return nil, wrapStoreError("list", err)
	}
	return &ListSkillBundlesResponse{
		Body: &ListSkillBundlesResponseBody{
			Bundles: values,
		},
	}, nil
}

func (a *StoreAPI) UpdateSkillBundle(
	ctx context.Context,
	request *UpdateSkillBundleRequest,
) (*UpdateSkillBundleResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Skill Bundle update",
	); err != nil {
		return nil, err
	}

	value, err := a.service.UpdateBundle(ctx, *request.Body)
	if err != nil {
		return nil, wrapStoreError("update", err)
	}
	return &UpdateSkillBundleResponse{Body: &value}, nil
}

func (a *StoreAPI) RetireSkillBundle(
	ctx context.Context,
	request *RetireSkillBundleRequest,
) (*RetireSkillBundleResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Skill Bundle retirement",
	); err != nil {
		return nil, err
	}

	value, err := a.service.RetireBundle(
		ctx,
		request.Bundle,
		request.ExpectedRevision,
	)
	if err != nil {
		return nil, wrapStoreError("retire", err)
	}
	return &RetireSkillBundleResponse{Body: &value}, nil
}

func (a *StoreAPI) PurgeSkillBundle(
	ctx context.Context,
	request *PurgeSkillBundleRequest,
) (*PurgeSkillBundleResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Skill Bundle purge",
	); err != nil {
		return nil, err
	}

	if err := a.service.PurgeBundle(
		ctx,
		request.Bundle,
		request.ExpectedRevision,
	); err != nil {
		return nil, wrapStoreError("purge", err)
	}
	return &PurgeSkillBundleResponse{
		Bundle: request.Bundle,
	}, nil
}

func (a *StoreAPI) AttachSkillBundleSource(
	ctx context.Context,
	request *AttachSkillBundleSourceRequest,
) (*AttachSkillBundleSourceResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Skill Bundle Source attachment",
	); err != nil {
		return nil, err
	}

	value, err := a.service.AttachSource(
		ctx,
		request.Body.Bundle,
		request.Body.ExpectedCollectionRevision,
		request.Body.Attachment,
	)
	if err != nil {
		return nil, wrapStoreError("attach Source", err)
	}
	return &AttachSkillBundleSourceResponse{Body: &value}, nil
}

func (a *StoreAPI) RefreshSkillBundle(
	ctx context.Context,
	request *RefreshSkillBundleRequest,
) (*RefreshSkillBundleResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Skill Bundle refresh",
	); err != nil {
		return nil, err
	}

	value, err := a.service.RefreshBundle(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("refresh", err)
	}
	return &RefreshSkillBundleResponse{Body: &value}, nil
}

func (a *StoreAPI) CreateManagedSkill(
	ctx context.Context,
	request *CreateManagedSkillStoreRequest,
) (*CreateManagedSkillStoreResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"managed Skill creation",
	); err != nil {
		return nil, err
	}

	value, err := a.service.CreateManagedSkill(ctx, *request.Body)
	if err != nil {
		return nil, wrapStoreError("create managed Skill", err)
	}
	return &CreateManagedSkillStoreResponse{Body: &value}, nil
}

func (a *StoreAPI) GetManagedSkillDocument(
	ctx context.Context,
	request *GetManagedSkillDocumentRequest,
) (*GetManagedSkillDocumentResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"managed Skill document get",
	); err != nil {
		return nil, err
	}

	value, err := a.service.GetManagedSkillDocument(ctx, request.Artifact)
	if err != nil {
		return nil, wrapStoreError("get managed Skill document", err)
	}
	return &GetManagedSkillDocumentResponse{Body: &value}, nil
}

func (a *StoreAPI) AdoptSkill(
	ctx context.Context,
	request *AdoptSkillStoreRequest,
) (*AdoptSkillStoreResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Skill adoption",
	); err != nil {
		return nil, err
	}

	value, err := a.service.AdoptSkill(ctx, *request.Body)
	if err != nil {
		return nil, wrapStoreError("adopt Skill", err)
	}
	return &AdoptSkillStoreResponse{Body: &value}, nil
}

func (a *StoreAPI) PinSkill(
	ctx context.Context,
	request *PinSkillStoreRequest,
) (*PinSkillStoreResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Skill pin",
	); err != nil {
		return nil, err
	}

	value, err := a.service.PinSkill(ctx, *request.Body)
	if err != nil {
		return nil, wrapStoreError("pin Skill", err)
	}
	return &PinSkillStoreResponse{Body: &value}, nil
}

func (a *StoreAPI) GetSkill(
	ctx context.Context,
	request *GetSkillRequest,
) (*GetSkillResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Skill get",
	); err != nil {
		return nil, err
	}

	value, err := a.service.GetSkill(ctx, request.Artifact)
	if err != nil {
		return nil, wrapStoreError("get Skill", err)
	}
	return &GetSkillResponse{Body: &value}, nil
}

func (a *StoreAPI) ListBundleSkills(
	ctx context.Context,
	request *ListBundleSkillsRequest,
) (*ListBundleSkillsResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Skill list",
	); err != nil {
		return nil, err
	}

	values, err := a.service.ListSkills(ctx, request.Bundle)
	if err != nil {
		return nil, wrapStoreError("list Skills", err)
	}
	return &ListBundleSkillsResponse{
		Body: &ListBundleSkillsResponseBody{
			Skills: values,
		},
	}, nil
}

func (a *StoreAPI) SetSkillEnabled(
	ctx context.Context,
	request *SetSkillEnabledRequest,
) (*SetSkillEnabledResponse, error) {
	if err := requireRequestBody(
		request,
		request != nil && request.Body != nil,
		true,
		"Skill enabled update",
	); err != nil {
		return nil, err
	}

	value, err := a.service.SetSkillEnabled(
		ctx,
		request.Body.Artifact,
		request.Body.ExpectedRevision,
		request.Body.Enabled,
	)
	if err != nil {
		return nil, wrapStoreError("set Skill enabled", err)
	}
	return &SetSkillEnabledResponse{Body: &value}, nil
}

func (a *StoreAPI) UnadoptSkill(
	ctx context.Context,
	request *UnadoptSkillRequest,
) (*UnadoptSkillResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Skill unadopt",
	); err != nil {
		return nil, err
	}

	if err := a.service.UnadoptSkill(
		ctx,
		request.Artifact,
		request.ExpectedRevision,
		request.Suppress,
	); err != nil {
		return nil, wrapStoreError("unadopt Skill", err)
	}
	return &UnadoptSkillResponse{
		Artifact: request.Artifact,
	}, nil
}

func (a *StoreAPI) PurgeSkill(
	ctx context.Context,
	request *PurgeSkillRequest,
) (*PurgeSkillResponse, error) {
	if err := requireRequestBody(
		request,
		false,
		false,
		"Skill purge",
	); err != nil {
		return nil, err
	}

	if err := a.service.PurgeSkill(
		ctx,
		request.Artifact,
		request.ExpectedRevision,
	); err != nil {
		return nil, wrapStoreError("purge Skill", err)
	}
	return &PurgeSkillResponse{
		Artifact: request.Artifact,
	}, nil
}
