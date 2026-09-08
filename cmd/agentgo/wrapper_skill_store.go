package main

import (
	"context"
	"errors"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
	skillStore "github.com/flexigpt/flexigpt-app/internal/skill/store"
	skillBundle "github.com/flexigpt/flexigpt-app/internal/skill/store/bundle"
	"github.com/flexigpt/flexigpt-app/internal/skill/store/schemaadapter"
	"github.com/flexigpt/flexigpt-app/internal/skill/store/workspaceadapter"
)

type SkillStoreWrapper struct {
	api    *skillBundle.StoreAPI
	router *skillStore.ArtifactRouter

	builtInInstaller artifactbuiltin.HydrationInstaller
}

func InitSkillStoreWrapper(
	wrapper *SkillStoreWrapper,
	dependencies skillBundle.Dependencies,
	workspaceSkills *workspaceadapter.Adapter,
) error {
	if wrapper == nil || workspaceSkills == nil {
		return errors.New("skill store wrapper dependencies are incomplete")
	}
	if err := dependencies.Validate(); err != nil {
		return err
	}

	service, err := skillBundle.New(dependencies)
	if err != nil {
		return err
	}
	api, err := skillBundle.NewStoreAPI(service)
	if err != nil {
		return err
	}

	router, err := skillStore.NewArtifactRouter(
		dependencies.Artifacts,
		dependencies.Collections,
	)
	if err != nil {
		return err
	}

	workspaceResolver, err := workspaceadapter.NewStoreLoader(workspaceSkills)
	if err != nil {
		return err
	}
	bundleResolver, err := skillBundle.NewStoreLoader(service)
	if err != nil {
		return err
	}
	if err := router.Register(
		artifactbuiltin.WorkspaceCollectionV1Kind,
		workspaceResolver,
	); err != nil {
		return err
	}
	if err := router.Register(
		artifactbuiltin.SkillCollectionV1Kind,
		bundleResolver,
	); err != nil {
		return err
	}

	skillRegistry, err := schemaadapter.LoadRegistry()
	if err != nil {
		return err
	}
	packages, err := artifactbuiltin.EmbeddedSkillPackages()
	if err != nil {
		return err
	}
	builtIns, err := schemaadapter.NewInstaller(
		schemaadapter.InstallerDependencies{
			Skills:                 service,
			SkillRegistry:          skillRegistry,
			Packages:               packages,
			ShareableCanonicalizer: dependencies.Schemas,
		},
	)
	if err != nil {
		return err
	}

	wrapper.api = api
	wrapper.router = router
	wrapper.builtInInstaller = builtIns
	return nil
}

func (w *SkillStoreWrapper) CreateSkillBundle(
	request *skillBundle.CreateSkillBundleRequest,
) (*skillBundle.CreateSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.CreateSkillBundleResponse, error) {
			return w.api.CreateSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) GetSkillBundle(
	request *skillBundle.GetSkillBundleRequest,
) (*skillBundle.GetSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.GetSkillBundleResponse, error) {
			return w.api.GetSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) ListSkillBundles(
	request *skillBundle.ListSkillBundlesRequest,
) (*skillBundle.ListSkillBundlesResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.ListSkillBundlesResponse, error) {
			return w.api.ListSkillBundles(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) UpdateSkillBundle(
	request *skillBundle.UpdateSkillBundleRequest,
) (*skillBundle.UpdateSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.UpdateSkillBundleResponse, error) {
			return w.api.UpdateSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) RetireSkillBundle(
	request *skillBundle.RetireSkillBundleRequest,
) (*skillBundle.RetireSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.RetireSkillBundleResponse, error) {
			return w.api.RetireSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) PurgeSkillBundle(
	request *skillBundle.PurgeSkillBundleRequest,
) (*skillBundle.PurgeSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.PurgeSkillBundleResponse, error) {
			return w.api.PurgeSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) AttachSkillBundleSource(
	request *skillBundle.AttachSkillBundleSourceRequest,
) (*skillBundle.AttachSkillBundleSourceResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.AttachSkillBundleSourceResponse, error) {
			return w.api.AttachSkillBundleSource(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) RefreshSkillBundle(
	request *skillBundle.RefreshSkillBundleRequest,
) (*skillBundle.RefreshSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.RefreshSkillBundleResponse, error) {
			return w.api.RefreshSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) CreateManagedSkill(
	request *skillBundle.CreateManagedSkillStoreRequest,
) (*skillBundle.CreateManagedSkillStoreResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.CreateManagedSkillStoreResponse, error) {
			return w.api.CreateManagedSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) GetManagedSkillDocument(
	request *skillBundle.GetManagedSkillDocumentRequest,
) (*skillBundle.GetManagedSkillDocumentResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.GetManagedSkillDocumentResponse, error) {
			return w.api.GetManagedSkillDocument(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) AdoptSkill(
	request *skillBundle.AdoptSkillStoreRequest,
) (*skillBundle.AdoptSkillStoreResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.AdoptSkillStoreResponse, error) {
			return w.api.AdoptSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) PinSkill(
	request *skillBundle.PinSkillStoreRequest,
) (*skillBundle.PinSkillStoreResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.PinSkillStoreResponse, error) {
			return w.api.PinSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) GetSkill(
	request *skillBundle.GetSkillRequest,
) (*skillBundle.GetSkillResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.GetSkillResponse, error) {
			return w.api.GetSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) ListBundleSkills(
	request *skillBundle.ListBundleSkillsRequest,
) (*skillBundle.ListBundleSkillsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.ListBundleSkillsResponse, error) {
			return w.api.ListBundleSkills(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) SetSkillEnabled(
	request *skillBundle.SetSkillEnabledRequest,
) (*skillBundle.SetSkillEnabledResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.SetSkillEnabledResponse, error) {
			return w.api.SetSkillEnabled(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) UnadoptSkill(
	request *skillBundle.UnadoptSkillRequest,
) (*skillBundle.UnadoptSkillResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.UnadoptSkillResponse, error) {
			return w.api.UnadoptSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) PurgeSkill(
	request *skillBundle.PurgeSkillRequest,
) (*skillBundle.PurgeSkillResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillBundle.PurgeSkillResponse, error) {
			return w.api.PurgeSkill(ctx, request)
		},
	)
}

// ResolveArtifactSkill is a Skill aggregate bridge, not a generic Artifact
// Store API. Runtime-facing migration is handled by the later aggregate
// checkpoint.
func (w *SkillStoreWrapper) ResolveArtifactSkill(
	ref artifact.ArtifactRef,
) (skillStore.ResolvedArtifactSkill, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (skillStore.ResolvedArtifactSkill, error) {
			return w.router.ResolveArtifactSkill(ctx, ref)
		},
	)
}

func (w *SkillStoreWrapper) close() {
	if w == nil {
		return
	}
	w.builtInInstaller = nil
	w.router = nil
	w.api = nil
}
