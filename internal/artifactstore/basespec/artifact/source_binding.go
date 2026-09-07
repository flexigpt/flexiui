package artifact

import (
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
)

type SourceBinding struct {
	SourceID           source.SourceID             `json:"sourceID"`
	Locator            basespec.Locator            `json:"locator"`
	SubresourceLocator basespec.SubresourceLocator `json:"subresourceLocator,omitempty"`
	ExpectedKind       ArtifactKind                `json:"expectedKind"`
}

func (b SourceBinding) Validate() error {
	if err := b.SourceID.Validate(); err != nil {
		return err
	}
	if err := b.Locator.Validate(true); err != nil {
		return err
	}
	if err := b.SubresourceLocator.Validate(); err != nil {
		return err
	}
	return b.ExpectedKind.Validate()
}
