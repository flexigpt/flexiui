package bundle

import (
	"context"
	"errors"

	"github.com/flexigpt/agentskills-go/provider"
	"github.com/flexigpt/agentskills-go/provider/fs"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	skillStore "github.com/flexigpt/flexigpt-app/internal/skill/store"
)

type RuntimeSkillSource interface {
	ResolveSkill(
		ctx context.Context,
		ref artifact.ArtifactRef,
	) (ResolvedSkill, error)

	ListResolvedSkills(
		ctx context.Context,
		ref collection.CollectionRef,
	) ([]ResolvedSkill, error)
}

type StoreLoader struct {
	source RuntimeSkillSource
}

func NewStoreLoader(
	source RuntimeSkillSource,
) (*StoreLoader, error) {
	if source == nil {
		return nil, errors.New("skill bundle runtime source is nil")
	}
	return &StoreLoader{source: source}, nil
}

func (r *StoreLoader) ResolveArtifactSkill(
	ctx context.Context,
	ref artifact.ArtifactRef,
) (skillStore.ResolvedArtifactSkill, error) {
	value, err := r.source.ResolveSkill(ctx, ref)
	if err != nil {
		return skillStore.ResolvedArtifactSkill{}, err
	}
	return resolvedArtifactSkillOf(value)
}

func (r *StoreLoader) ListCollectionSkills(
	ctx context.Context,
	ref collection.CollectionRef,
) ([]skillStore.ResolvedArtifactSkill, error) {
	values, err := r.source.ListResolvedSkills(ctx, ref)
	if err != nil {
		return nil, err
	}

	output := make([]skillStore.ResolvedArtifactSkill, 0, len(values))
	for _, value := range values {
		projected, err := resolvedArtifactSkillOf(value)
		if err != nil {
			return nil, err
		}
		output = append(output, projected)
	}
	return output, nil
}

func resolvedArtifactSkillOf(
	value ResolvedSkill,
) (skillStore.ResolvedArtifactSkill, error) {
	output := skillStore.ResolvedArtifactSkill{
		Artifact:   value.Artifact,
		Collection: value.Collection,
		Definition: provider.SkillDef{
			Type:     fs.Type,
			Name:     value.Name,
			Location: value.Location,
		},
		Version: value.Version,
	}
	if err := output.Validate(); err != nil {
		return skillStore.ResolvedArtifactSkill{}, err
	}
	return output, nil
}
