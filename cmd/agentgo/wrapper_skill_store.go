package main

import (
	"context"
	"errors"
	"sort"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/middleware"
	skillBuiltin "github.com/flexigpt/flexigpt-app/internal/skill/store/builtin"
	skillConsumerAPI "github.com/flexigpt/flexigpt-app/internal/skill/store/consumerapi"
	skillDomain "github.com/flexigpt/flexigpt-app/internal/skill/store/domain"
)

func NewSkillBuiltInInstaller(
	skills skillConsumerAPI.BuiltinStore,
	schemas compositionapi.SchemaAPI,
) (artifactbuiltin.HydrationInstaller, error) {
	if skills == nil || schemas == nil {
		return nil, errors.New("skill built-in installer dependencies are incomplete")
	}

	registry, err := skillBuiltin.LoadRegistry()
	if err != nil {
		return nil, err
	}

	packages, err := artifactbuiltin.EmbeddedSkillPackages()
	if err != nil {
		return nil, err
	}

	return skillBuiltin.NewInstaller(
		skillBuiltin.InstallerDependencies{
			Skills:                 skills,
			SkillRegistry:          registry,
			Packages:               packages,
			ShareableCanonicalizer: schemas,
		},
	)
}

type SkillStoreWrapper struct {
	api   *skillConsumerAPI.API
	roots compositionapi.RootAPI
}

func InitSkillStoreWrapper(
	wrapper *SkillStoreWrapper,
	roots compositionapi.RootAPI,
	sources compositionapi.SourceAPI,
	collections compositionapi.CollectionAPI,
	artifacts compositionapi.ArtifactAPI,
	catalogs compositionapi.CatalogAPI,
	resources compositionapi.ResourceAPI,
	managedArtifacts compositionapi.ManagedArtifactAPI,
	protection compositionapi.ProtectionAPI,
) error {
	if wrapper == nil ||
		roots == nil {
		return errors.New("skill store wrapper dependencies are incomplete")
	}

	api, err := skillConsumerAPI.New(

		sources,
		collections,
		artifacts,
		catalogs,
		resources,
		managedArtifacts,
		protection,
	)
	if err != nil {
		return err
	}

	wrapper.api = api
	wrapper.roots = roots
	return nil
}

func (w *SkillStoreWrapper) CreateSkillBundle(
	request *skillConsumerAPI.CreateSkillBundleRequest,
) (*skillConsumerAPI.CreateSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.CreateSkillBundleResponse, error) {
			return w.api.CreateSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) GetSkillBundle(
	request *skillConsumerAPI.GetSkillBundleRequest,
) (*skillConsumerAPI.GetSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.GetSkillBundleResponse, error) {
			return w.api.GetSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) ListSkillBundles(
	request *skillConsumerAPI.ListSkillBundlesRequest,
) (*skillConsumerAPI.ListSkillBundlesResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.ListSkillBundlesResponse, error) {
			return w.api.ListSkillBundles(ctx, request)
		},
	)
}

// ListSkillBundlesForManagement hides Artifact Store root topology from the
// frontend while including user and protected built-in Skill Bundles.
func (w *SkillStoreWrapper) ListSkillBundlesForManagement() (
	[]skillDomain.SkillBundle,
	error,
) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() ([]skillDomain.SkillBundle, error) {
			if w == nil || w.api == nil || w.roots == nil {
				return nil, basespec.ErrClosed
			}

			roots, err := w.roots.List(ctx)
			if err != nil {
				return nil, err
			}

			bundles := make([]skillDomain.SkillBundle, 0)
			for _, rootValue := range roots {
				values, err := w.api.ListBundles(ctx, rootValue.ID)
				if err != nil {
					return nil, err
				}
				bundles = append(bundles, values...)
			}

			sort.Slice(bundles, func(left, right int) bool {
				if bundles[left].Collection.RootID != bundles[right].Collection.RootID {
					return bundles[left].Collection.RootID < bundles[right].Collection.RootID
				}
				return bundles[left].Collection.ID < bundles[right].Collection.ID
			})
			return bundles, nil
		},
	)
}

func (w *SkillStoreWrapper) UpdateSkillBundle(
	request *skillConsumerAPI.UpdateSkillBundleRequest,
) (*skillConsumerAPI.UpdateSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.UpdateSkillBundleResponse, error) {
			return w.api.UpdateSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) RetireSkillBundle(
	request *skillConsumerAPI.RetireSkillBundleRequest,
) (*skillConsumerAPI.RetireSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.RetireSkillBundleResponse, error) {
			return w.api.RetireSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) PurgeSkillBundle(
	request *skillConsumerAPI.PurgeSkillBundleRequest,
) (*skillConsumerAPI.PurgeSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.PurgeSkillBundleResponse, error) {
			return w.api.PurgeSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) AttachSkillBundleSource(
	request *skillConsumerAPI.AttachSkillBundleSourceRequest,
) (*skillConsumerAPI.AttachSkillBundleSourceResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.AttachSkillBundleSourceResponse, error) {
			return w.api.AttachSkillBundleSource(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) RegisterSkillBundleDirectory(
	request *skillConsumerAPI.RegisterSkillBundleDirectoryRequest,
) (*skillConsumerAPI.RegisterSkillBundleDirectoryResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.RegisterSkillBundleDirectoryResponse, error) {
			return w.api.RegisterSkillBundleDirectory(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) RefreshSkillBundle(
	request *skillConsumerAPI.RefreshSkillBundleRequest,
) (*skillConsumerAPI.RefreshSkillBundleResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.RefreshSkillBundleResponse, error) {
			return w.api.RefreshSkillBundle(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) CreateManagedSkill(
	request *skillConsumerAPI.CreateManagedSkillRequest,
) (*skillConsumerAPI.CreateManagedSkillStoreResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.CreateManagedSkillStoreResponse, error) {
			return w.api.CreateManagedSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) GetManagedSkillDocument(
	request *skillConsumerAPI.GetManagedSkillDocumentRequest,
) (*skillConsumerAPI.GetManagedSkillDocumentResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.GetManagedSkillDocumentResponse, error) {
			return w.api.GetManagedSkillDocument(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) AdoptSkill(
	request *skillConsumerAPI.AdoptSkillRequest,
) (*skillConsumerAPI.AdoptSkillResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.AdoptSkillResponse, error) {
			return w.api.AdoptSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) PinSkill(
	request *skillConsumerAPI.PinSkillRequest,
) (*skillConsumerAPI.PinSkillResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.PinSkillResponse, error) {
			return w.api.PinSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) GetSkill(
	request *skillConsumerAPI.GetSkillRequest,
) (*skillConsumerAPI.GetSkillResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.GetSkillResponse, error) {
			return w.api.GetSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) ListBundleSkills(
	request *skillConsumerAPI.ListBundleSkillsRequest,
) (*skillConsumerAPI.ListBundleSkillsResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.ListBundleSkillsResponse, error) {
			return w.api.ListBundleSkills(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) SetSkillEnabled(
	request *skillConsumerAPI.SetSkillEnabledRequest,
) (*skillConsumerAPI.SetSkillEnabledResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.SetSkillEnabledResponse, error) {
			return w.api.SetSkillEnabled(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) UnadoptSkill(
	request *skillConsumerAPI.UnadoptSkillRequest,
) (*skillConsumerAPI.UnadoptSkillResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.UnadoptSkillResponse, error) {
			return w.api.UnadoptSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) PurgeSkill(
	request *skillConsumerAPI.PurgeSkillRequest,
) (*skillConsumerAPI.PurgeSkillResponse, error) {
	ctx := context.Background()

	return middleware.WithRecoveryResp(
		func() (*skillConsumerAPI.PurgeSkillResponse, error) {
			return w.api.PurgeSkill(ctx, request)
		},
	)
}

func (w *SkillStoreWrapper) close() {
	if w == nil {
		return
	}
	w.api = nil
	w.roots = nil
}
