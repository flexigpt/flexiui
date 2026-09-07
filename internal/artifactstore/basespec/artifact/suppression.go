package artifact

import (
	"fmt"
	"time"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
)

type Suppression struct {
	RootID       root.RootID             `json:"rootID"`
	CollectionID collection.CollectionID `json:"collectionID"`
	Binding      SourceBinding           `json:"binding"`
	Revision     uint64                  `json:"revision"`
	CreatedAt    time.Time               `json:"createdAt"`
	ModifiedAt   time.Time               `json:"modifiedAt"`
}

func (s Suppression) Validate() error {
	if err := s.RootID.Validate(); err != nil {
		return err
	}
	if err := s.CollectionID.Validate(); err != nil {
		return err
	}
	if err := s.Binding.Validate(); err != nil {
		return err
	}
	if s.Revision == 0 {
		return fmt.Errorf(
			"%w: suppression revision must be positive",
			basespec.ErrInvalid,
		)
	}
	if s.CreatedAt.IsZero() || s.ModifiedAt.IsZero() {
		return fmt.Errorf(
			"%w: suppression timestamps are required",
			basespec.ErrInvalid,
		)
	}

	if s.ModifiedAt.Before(s.CreatedAt) {
		return fmt.Errorf(
			"%w: suppression modified time precedes creation",
			basespec.ErrInvalid,
		)
	}
	return nil
}
