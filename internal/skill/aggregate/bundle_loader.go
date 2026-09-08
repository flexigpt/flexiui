package aggregate

import (
	"context"
	"errors"

	"github.com/flexigpt/agentskills-go/provider"
	"github.com/flexigpt/agentskills-go/provider/fs"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
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

type BundleLoader struct {
	source RuntimeSkillSource
}

func NewBundleLoader(
	source RuntimeSkillSource,
) (*BundleLoader, error) {
	if source == nil {
		return nil, errors.New("skill bundle runtime source is nil")
	}
	return &BundleLoader{source: source}, nil
}

func (r *BundleLoader) ResolveArtifactSkill(
	ctx context.Context,
	ref artifact.ArtifactRef,
) (ResolvedArtifactSkill, error) {
	value, err := r.source.ResolveSkill(ctx, ref)
	if err != nil {
		return ResolvedArtifactSkill{}, err
	}
	return resolvedArtifactSkillOf(value)
}

func (r *BundleLoader) ListCollectionSkills(
	ctx context.Context,
	ref collection.CollectionRef,
) ([]ResolvedArtifactSkill, error) {
	values, err := r.source.ListResolvedSkills(ctx, ref)
	if err != nil {
		return nil, err
	}

	output := make([]ResolvedArtifactSkill, 0, len(values))
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
) (ResolvedArtifactSkill, error) {
	output := ResolvedArtifactSkill{
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
		return ResolvedArtifactSkill{}, err
	}
	return output, nil
}
