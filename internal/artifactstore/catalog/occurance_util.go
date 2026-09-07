package catalog

import (
	"maps"
	"sort"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
)

func SortOccurrences(values []Occurrence) {
	sort.Slice(values, func(left, right int) bool {
		if values[left].Key.CollectionID != values[right].Key.CollectionID {
			return values[left].Key.CollectionID < values[right].Key.CollectionID
		}
		if values[left].Key.SourceID != values[right].Key.SourceID {
			return values[left].Key.SourceID < values[right].Key.SourceID
		}
		if values[left].Key.Locator != values[right].Key.Locator {
			return values[left].Key.Locator < values[right].Key.Locator
		}
		return values[left].Key.SubresourceLocator <
			values[right].Key.SubresourceLocator
	})
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

	return equalOccurrences(left.Occurrences, right.Occurrences)
}

// equalOccurrences compares occurrence values independently of ordering.
func equalOccurrences(left, right []Occurrence) bool {
	if len(left) != len(right) {
		return false
	}

	byKey := make(map[OccurrenceKey]Occurrence, len(left))
	for _, value := range left {
		if _, duplicate := byKey[value.Key]; duplicate {
			return false
		}
		byKey[value.Key] = value
	}

	for _, value := range right {
		leftValue, found := byKey[value.Key]
		if !found || !equalOccurrence(leftValue, value) {
			return false
		}
	}
	return true
}

func equalOccurrence(left, right Occurrence) bool {
	if left.Definition == nil || right.Definition == nil {
		if left.Definition != nil || right.Definition != nil {
			return false
		}
	} else if left.Definition.Digest != right.Definition.Digest {
		return false
	}

	return left.RootID == right.RootID &&
		left.CollectionID == right.CollectionID &&
		left.Key == right.Key &&
		left.Kind == right.Kind &&
		left.LogicalName == right.LogicalName &&
		left.LogicalVersion == right.LogicalVersion &&
		cryptoutil.IsDigestEqual(left.DefinitionDigest, right.DefinitionDigest) &&
		cryptoutil.IsDigestEqual(left.SourceContentDigest, right.SourceContentDigest) &&
		left.DecoderID == right.DecoderID &&
		left.State == right.State &&
		left.ObservedAt.Equal(right.ObservedAt) &&
		providerapi.EqualDiagnostics(left.Diagnostics, right.Diagnostics)
}
