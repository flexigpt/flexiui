package workspaceadapter

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/flexigpt/agentskills-go/provider"
	"github.com/flexigpt/agentskills-go/provider/fs"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
	skillAggregate "github.com/flexigpt/flexigpt-app/internal/skill/aggregate"
)

// StoreLoader adapts Workspace-owned Artifact projection to the generic
// Artifact Skill router. Runtime does not import Workspace and cannot infer
// Workspace ownership from a reference shape.
type StoreLoader struct {
	adapter *Adapter
}

func NewStoreLoader(adapter *Adapter) (*StoreLoader, error) {
	if adapter == nil {
		return nil, errors.New("workspace skill store loader adapter is nil")
	}
	return &StoreLoader{adapter: adapter}, nil
}

func (r *StoreLoader) ResolveArtifactSkill(
	ctx context.Context,
	ref artifact.ArtifactRef,
) (skillAggregate.ResolvedArtifactSkill, error) {
	value, err := r.adapter.LoadArtifact(ctx, ref)
	if err != nil {
		return skillAggregate.ResolvedArtifactSkill{}, err
	}
	return workspaceResolvedSkill(value)
}

func (r *StoreLoader) ListCollectionSkills(
	ctx context.Context,
	workspace collection.CollectionRef,
) ([]skillAggregate.ResolvedArtifactSkill, error) {
	plan, err := r.adapter.LoadAll(ctx, workspace)
	if err != nil {
		return nil, err
	}
	if len(plan.Skills) == 0 {
		return []skillAggregate.ResolvedArtifactSkill{}, nil
	}

	output := make([]skillAggregate.ResolvedArtifactSkill, 0, len(plan.Skills))
	for _, value := range plan.Skills {
		projected, err := workspaceResolvedSkill(value)
		if err != nil {
			return nil, err
		}
		output = append(output, projected)
	}
	return output, nil
}

func workspaceResolvedSkill(
	value WorkspaceSkill,
) (skillAggregate.ResolvedArtifactSkill, error) {
	if !value.ProjectionValid ||
		value.RuntimeLocation == "" ||
		value.SourceContentDigest == "" ||
		value.SourceGeneration == "" ||
		value.ArtifactRevision == 0 {
		return skillAggregate.ResolvedArtifactSkill{}, fmt.Errorf(
			"%w: Workspace Skill is missing prepared runtime material",
			basespec.ErrReferenceUnresolved,
		)
	}

	versionInput := string(value.DefinitionDigest) + "\x00" +
		string(value.SourceContentDigest) + "\x00" +
		value.SourceGeneration + "\x00" +
		strconv.FormatUint(value.ArtifactRevision, 10)

	output := skillAggregate.ResolvedArtifactSkill{
		Artifact:   value.Artifact,
		Collection: value.Workspace,
		Definition: provider.SkillDef{
			Type:     fs.Type,
			Name:     value.Skill.Name,
			Location: value.RuntimeLocation,
		},
		Version: "workspace:" + string(
			cryptoutil.DigestBytes([]byte(versionInput)),
		),
	}
	if err := output.Validate(); err != nil {
		return skillAggregate.ResolvedArtifactSkill{}, err
	}
	return output, nil
}
