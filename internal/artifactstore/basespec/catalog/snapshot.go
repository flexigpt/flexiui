package catalog

import (
	"fmt"
	"maps"
	"time"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
)

type Snapshot struct {
	RootID              basespec.RootID              `json:"rootID"`
	CollectionID        basespec.CollectionID        `json:"collectionID"`
	Revision            uint64                       `json:"revision"`
	CollectionRevision  uint64                       `json:"collectionRevision"`
	AttachmentRevisions map[basespec.SourceID]uint64 `json:"attachmentRevisions"`
	SourceRevisions     map[basespec.SourceID]uint64 `json:"sourceRevisions"`
	SourceGenerations   map[basespec.SourceID]string `json:"sourceGenerations"`
	PlanFingerprint     cryptoutil.Digest            `json:"planFingerprint"`
	DecoderFingerprint  cryptoutil.Digest            `json:"decoderFingerprint"`
	PublishedAt         time.Time                    `json:"publishedAt"`
	Diagnostics         []providerapi.Diagnostic     `json:"diagnostics,omitempty"`
	Occurrences         []Occurrence                 `json:"occurrences"`
}

func (s Snapshot) Validate() error {
	if err := basespec.ValidateRootID(s.RootID); err != nil {
		return err
	}
	if err := basespec.ValidateCollectionID(s.CollectionID); err != nil {
		return err
	}
	if s.Revision == 0 || s.CollectionRevision == 0 {
		return fmt.Errorf("%w: catalog revisions must be positive", basespec.ErrInvalid)
	}
	if err := cryptoutil.ValidateDigest(s.PlanFingerprint); err != nil {
		return fmt.Errorf("catalog plan fingerprint: %w", err)
	}
	if err := cryptoutil.ValidateDigest(s.DecoderFingerprint); err != nil {
		return fmt.Errorf("catalog decoder fingerprint: %w", err)
	}
	for sourceID, revision := range s.AttachmentRevisions {
		if err := basespec.ValidateSourceID(sourceID); err != nil {
			return err
		}
		if revision == 0 {
			return fmt.Errorf("%w: attachment revision must be positive", basespec.ErrInvalid)
		}
	}
	for sourceID := range s.AttachmentRevisions {
		if _, exists := s.SourceRevisions[sourceID]; !exists {
			return fmt.Errorf(
				"%w: collection attachment has no corresponding source revision",
				basespec.ErrInvalid,
			)
		}
	}
	for sourceID, revision := range s.SourceRevisions {
		if err := basespec.ValidateSourceID(sourceID); err != nil {
			return err
		}
		if revision == 0 {
			return fmt.Errorf("%w: source revision must be positive", basespec.ErrInvalid)
		}
		if _, attached := s.AttachmentRevisions[sourceID]; !attached {
			return fmt.Errorf(
				"%w: source revision has no collection attachment",
				basespec.ErrInvalid,
			)
		}
	}
	for sourceID, generation := range s.SourceGenerations {
		if err := basespec.ValidateSourceID(sourceID); err != nil {
			return err
		}
		if _, exists := s.SourceRevisions[sourceID]; !exists {
			return fmt.Errorf(
				"%w: source generation has no source revision",
				basespec.ErrInvalid,
			)
		}
		if err := basespec.ValidateSourceGeneration(generation); err != nil {
			return err
		}
	}
	if s.PublishedAt.IsZero() {
		return fmt.Errorf("%w: catalog publication time is required", basespec.ErrInvalid)
	}
	if err := providerapi.ValidateDiagnostics(s.Diagnostics); err != nil {
		return err
	}
	seenOccurrences := make(map[OccurrenceKey]struct{}, len(s.Occurrences))
	for index, occurrence := range s.Occurrences {
		if _, duplicate := seenOccurrences[occurrence.Key]; duplicate {
			return fmt.Errorf(
				"%w: duplicate occurrence %d",
				basespec.ErrInvalid,
				index,
			)
		}
		seenOccurrences[occurrence.Key] = struct{}{}
		if occurrence.RootID != s.RootID {
			return fmt.Errorf(
				"%w: occurrence %d belongs to another root",
				basespec.ErrInvalid,
				index,
			)
		}
		if occurrence.CollectionID != s.CollectionID ||
			occurrence.Key.CollectionID != s.CollectionID {
			return fmt.Errorf(
				"%w: occurrence %d belongs to another collection",
				basespec.ErrInvalid,
				index,
			)
		}
		if _, exists := s.SourceRevisions[occurrence.Key.SourceID]; !exists {
			return fmt.Errorf(
				"%w: occurrence %d has no source revision",
				basespec.ErrInvalid,
				index,
			)
		}
		if _, exists := s.SourceGenerations[occurrence.Key.SourceID]; !exists {
			return fmt.Errorf(
				"%w: occurrence %d has no source generation",
				basespec.ErrInvalid,
				index,
			)
		}
		if err := occurrence.Validate(); err != nil {
			return fmt.Errorf("occurrence %d: %w", index, err)
		}
	}
	return nil
}

func (s Snapshot) Clone() Snapshot {
	output := s
	output.AttachmentRevisions = make(
		map[basespec.SourceID]uint64,
		len(s.AttachmentRevisions),
	)
	maps.Copy(output.AttachmentRevisions, s.AttachmentRevisions)
	output.SourceRevisions = make(
		map[basespec.SourceID]uint64,
		len(s.SourceRevisions),
	)
	maps.Copy(output.SourceRevisions, s.SourceRevisions)
	output.SourceGenerations = make(
		map[basespec.SourceID]string,
		len(s.SourceGenerations),
	)
	maps.Copy(output.SourceGenerations, s.SourceGenerations)
	output.Diagnostics = providerapi.CloneDiagnostics(s.Diagnostics)
	output.Occurrences = make([]Occurrence, len(s.Occurrences))
	for index, occurrence := range s.Occurrences {
		output.Occurrences[index] = occurrence.Clone()
	}
	return output
}

// DefinitionForOccurrence returns the SQLite-cached canonical definition for
// one current catalog occurrence. It does not perform CAS lookup.
func (s Snapshot) DefinitionForOccurrence(
	key OccurrenceKey,
) (providerapi.Definition, error) {
	for _, occurrence := range s.Occurrences {
		if occurrence.Key != key {
			continue
		}
		if occurrence.State != OccurrenceValid ||
			occurrence.Definition == nil {
			return providerapi.Definition{}, fmt.Errorf(
				"%w: occurrence %q has no available definition",
				basespec.ErrDefinitionNotFound,
				key.Locator,
			)
		}
		return occurrence.Definition.Clone(), nil
	}
	return providerapi.Definition{}, fmt.Errorf(
		"%w: occurrence %q",
		basespec.ErrDefinitionNotFound,
		key.Locator,
	)
}

// EqualSnapshot compares the semantic contents of two catalog snapshots.
// Occurrences are compared by occurrence key because their persisted ordering
// is not part of catalog identity.
func EqualSnapshot(left, right Snapshot) bool {
	if left.RootID != right.RootID ||
		left.CollectionID != right.CollectionID ||
		left.Revision != right.Revision ||
		left.CollectionRevision != right.CollectionRevision ||
		left.PlanFingerprint != right.PlanFingerprint ||
		left.DecoderFingerprint != right.DecoderFingerprint ||
		!left.PublishedAt.Equal(right.PublishedAt) ||
		!maps.Equal(left.AttachmentRevisions, right.AttachmentRevisions) ||
		!maps.Equal(left.SourceRevisions, right.SourceRevisions) ||
		!maps.Equal(left.SourceGenerations, right.SourceGenerations) ||
		!providerapi.EqualDiagnostics(left.Diagnostics, right.Diagnostics) {
		return false
	}

	return EqualOccurrences(left.Occurrences, right.Occurrences)
}
