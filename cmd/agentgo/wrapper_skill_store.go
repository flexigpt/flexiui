package main

import (
	"context"
	"errors"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	artifactConsumerAPI "github.com/flexigpt/flexigpt-app/internal/artifactstore/consumerapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
	skillStore "github.com/flexigpt/flexigpt-app/internal/skill/store"
	skillBundle "github.com/flexigpt/flexigpt-app/internal/skill/store/bundle"
	"github.com/flexigpt/flexigpt-app/internal/skill/store/schemaadapter"
	"github.com/flexigpt/flexigpt-app/internal/skill/store/workspaceadapter"
)

type SkillStoreWrapper struct {
	api    *skillBundle.API
	router *skillStore.ArtifactRouter

	builtInInstaller artifactbuiltin.HydrationInstaller
}

func InitSkillStoreWrapper(
	wrapper *SkillStoreWrapper,
	store *artifactConsumerAPI.API,
	workspaceSkills *workspaceadapter.Adapter,
) error {
	if wrapper == nil || store == nil || workspaceSkills == nil {
		return errors.New("skill store wrapper dependencies are incomplete")
	}

	api, err := skillBundle.New(
		skillBundle.Dependencies{Store: store},
	)
	if err != nil {
		return err
	}

	router, err := skillStore.NewArtifactRouter(store)
	if err != nil {
		return err
	}

	workspaceResolver, err := workspaceadapter.NewStoreLoader(workspaceSkills)
	if err != nil {
		return err
	}
	bundleResolver, err := skillBundle.NewStoreLoader(api)
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

	wrapper.api = api
	wrapper.router = router

	skillRegistry, err := schemaadapter.LoadRegistry()
	if err != nil {
		wrapper.close()
		return err
	}

	packages, err := artifactbuiltin.EmbeddedSkillPackages()
	if err != nil {
		wrapper.close()
		return err
	}
	builtIns, err := schemaadapter.NewInstaller(
		schemaadapter.InstallerDependencies{
			Skills:                 api,
			SkillRegistry:          skillRegistry,
			Packages:               packages,
			ShareableCanonicalizer: store,
		},
	)
	if err != nil {
		wrapper.close()
		return err
	}
	wrapper.builtInInstaller = builtIns

	return nil
}

// AttachSkillSource attaches a Source already created through Artifact Store
// administration. It intentionally accepts only Source identity and typed
// Skill Bundle role data, never a filesystem path or Source configuration.

func (w *SkillStoreWrapper) GetSkillBundle(
	ref collection.CollectionRef,
) (skillBundle.Bundle, error) {
	return middleware.WithRecoveryResp(func() (skillBundle.Bundle, error) {
		return w.api.GetBundle(context.Background(), ref)
	})
}

func (w *SkillStoreWrapper) ListSkillBundles(
	rootID root.RootID,
) ([]skillBundle.Bundle, error) {
	return middleware.WithRecoveryResp(func() ([]skillBundle.Bundle, error) {
		return w.api.ListBundles(context.Background(), rootID)
	})
}

func (w *SkillStoreWrapper) CreateManagedSkill(
	request *skillBundle.CreateManagedSkillRequest,
) (skillBundle.CreateManagedSkillResponse, error) {
	return middleware.WithRecoveryResp(
		func() (skillBundle.CreateManagedSkillResponse, error) {
			if request == nil {
				return skillBundle.CreateManagedSkillResponse{},
					errors.New("managed skill request is required")
			}
			value, err := w.api.CreateManagedSkill(context.Background(), *request)
			if err != nil {
				return value, err
			}

			return value, nil
		},
	)
}

func (w *SkillStoreWrapper) GetManagedSkillDocument(
	ref artifact.ArtifactRef,
) (skillBundle.ManagedSkillDocument, error) {
	return middleware.WithRecoveryResp(
		func() (skillBundle.ManagedSkillDocument, error) {
			return w.api.GetManagedSkillDocument(context.Background(), ref)
		},
	)
}

func (w *SkillStoreWrapper) ListBundleSkills(
	ref collection.CollectionRef,
) ([]artifact.Artifact, error) {
	return middleware.WithRecoveryResp(func() ([]artifact.Artifact, error) {
		return w.api.ListSkills(context.Background(), ref)
	})
}

// ResolveArtifactSkill performs the Store-owned ArtifactRef to Agent Skills
// translation. The returned Collection identifies the catalog that must be
// synchronized before Definition is used in a runtime session.
func (w *SkillStoreWrapper) ResolveArtifactSkill(
	ref artifact.ArtifactRef,
) (skillStore.ResolvedArtifactSkill, error) {
	return middleware.WithRecoveryResp(
		func() (skillStore.ResolvedArtifactSkill, error) {
			return w.router.ResolveArtifactSkill(
				context.Background(),
				ref,
			)
		},
	)
}

func (w *SkillStoreWrapper) close() {
	if w == nil {
		return
	}

	api := w.api
	w.builtInInstaller = nil

	w.router = nil
	w.api = nil

	if api != nil {
		_ = api.Close()
	}
}
