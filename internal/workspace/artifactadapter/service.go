package artifactadapter

import (
	"context"
	"fmt"
	"sort"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/compositionapi"
	"github.com/flexigpt/flexigpt-app/internal/workspace/collectiondata"
	"github.com/flexigpt/flexigpt-app/internal/workspace/spec"
)

type Service struct {
	collections     compositionapi.CollectionAPI
	sources         compositionapi.SourceAPI
	workspaceRootID root.RootID
}

func NewService(
	collections compositionapi.CollectionAPI,
	sources compositionapi.SourceAPI,
	workspaceRootID root.RootID,
) (*Service, error) {
	if collections == nil || sources == nil {
		return nil, fmt.Errorf(
			"%w: Workspace service dependencies are incomplete",
			spec.ErrInvalidWorkspace,
		)
	}
	if err := workspaceRootID.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		collections:     collections,
		sources:         sources,
		workspaceRootID: workspaceRootID,
	}, nil
}

func (s *Service) List(
	ctx context.Context,
	rootID root.RootID,
) ([]spec.Workspace, error) {
	if err := s.requireWorkspaceRoot(rootID); err != nil {
		return nil, err
	}

	collections, err := s.collections.ListByRoot(ctx, rootID)
	if err != nil {
		return nil, err
	}

	output := make([]spec.Workspace, 0)
	for _, value := range collections {
		if value.Kind != artifactbuiltin.WorkspaceCollectionV1Kind {
			continue
		}

		workspaceValue, err := s.Get(ctx, value.Ref())
		if err != nil {
			return nil, err
		}
		output = append(output, workspaceValue)
	}
	return output, nil
}

func (s *Service) Get(
	ctx context.Context,
	ref collection.CollectionRef,
) (spec.Workspace, error) {
	if err := s.requireWorkspaceRoot(ref.RootID); err != nil {
		return spec.Workspace{}, err
	}
	if err := ref.Validate(); err != nil {
		return spec.Workspace{}, err
	}

	value, err := s.collections.Get(ctx, ref)
	if err != nil {
		return spec.Workspace{}, err
	}
	if err := value.Validate(); err != nil {
		return spec.Workspace{}, fmt.Errorf(
			"%w: Collection reader returned an invalid Workspace Collection: %w",
			spec.ErrInvalidWorkspace,
			err,
		)
	}
	if value.Ref() != ref {
		return spec.Workspace{}, fmt.Errorf(
			"%w: Collection reader returned another Workspace",
			spec.ErrInvalidWorkspace,
		)
	}
	if value.Kind != artifactbuiltin.WorkspaceCollectionV1Kind {
		return spec.Workspace{}, fmt.Errorf(
			"%w: collection %q has kind %q",
			spec.ErrNotWorkspace,
			ref.CollectionID,
			value.Kind,
		)
	}

	data, err := collectiondata.DecodeCollectionData(value.Data)
	if err != nil {
		return spec.Workspace{}, fmt.Errorf(
			"%w: %w",
			spec.ErrInvalidWorkspace,
			err,
		)
	}

	attachments, err := s.collections.ListAttachments(ctx, ref)
	if err != nil {
		return spec.Workspace{}, err
	}
	sort.Slice(attachments, func(left, right int) bool {
		return attachments[left].SourceID < attachments[right].SourceID
	})

	sources := make([]source.Summary, 0, len(attachments))
	for _, attachment := range attachments {
		sourceValue, err := s.sources.Get(
			ctx,
			ref.RootID,
			attachment.SourceID,
		)
		if err != nil {
			return spec.Workspace{}, err
		}
		sources = append(sources, sourceValue)
	}

	mode, primarySourceID, err := validateWorkspaceState(
		value,
		data,
		attachments,
		sources,
	)
	if err != nil {
		return spec.Workspace{}, err
	}

	sort.Slice(sources, func(left, right int) bool {
		return sources[left].ID < sources[right].ID
	})

	return spec.Workspace{
		Collection:      value,
		Data:            data,
		Mode:            mode,
		PrimarySourceID: primarySourceID,
		Attachments:     attachments,
		Sources:         sources,
	}, nil
}

func (s *Service) requireWorkspaceRoot(
	rootID root.RootID,
) error {
	if err := rootID.Validate(); err != nil {
		return err
	}
	if rootID != s.workspaceRootID {
		return fmt.Errorf(
			"%w: Workspace Root must be %q",
			spec.ErrInvalidWorkspace,
			s.workspaceRootID,
		)
	}
	return nil
}
