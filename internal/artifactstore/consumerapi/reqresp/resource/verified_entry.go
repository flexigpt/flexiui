package resource

import (
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
)

// VerifiedEntry contains bytes read from an exact current Collection Source
// generation. The bytes are owned by the returned value.
type VerifiedEntry struct {
	Collection       collection.CollectionRef
	SourceID         source.SourceID
	CatalogRevision  uint64
	SourceRevision   uint64
	SourceGeneration string
	Content          []byte
	Digest           cryptoutil.Digest
}

func (e VerifiedEntry) Validate() error {
	if err := e.Collection.Validate(); err != nil {
		return err
	}
	if err := e.SourceID.Validate(); err != nil {
		return err
	}
	if e.CatalogRevision == 0 || e.SourceRevision == 0 {
		return fmt.Errorf(
			"%w: verified entry revisions are required",
			basespec.ErrInvalid,
		)
	}
	if err := basespec.ValidateSourceGeneration(e.SourceGeneration); err != nil {
		return err
	}
	if err := cryptoutil.ValidateDigest(e.Digest); err != nil {
		return err
	}
	if cryptoutil.DigestBytes(e.Content) != e.Digest {
		return fmt.Errorf(
			"%w: verified entry content does not match its digest",
			basespec.ErrDigestMismatch,
		)
	}
	return nil
}

func (e VerifiedEntry) Clone() VerifiedEntry {
	output := e
	output.Content = append([]byte(nil), e.Content...)
	return output
}
