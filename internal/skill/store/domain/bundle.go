package domain

import (
	"github.com/flexigpt/agentskills-go/document"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
)

type AttachmentDraft struct {
	SourceID              source.SourceID                        `json:"sourceId"`
	Role                  collection.AttachmentRole              `json:"role"`
	Enabled               bool                                   `json:"enabled"`
	DiscoveryRoot         basespec.Locator                       `json:"discoveryRoot"`
	ExpectedMemberDigests map[basespec.Locator]cryptoutil.Digest `json:"expectedMemberDigests"`
}

type SkillBundle struct {
	Collection  collection.Collection   `json:"collection"`
	Data        CollectionData          `json:"data"`
	Attachments []collection.Attachment `json:"attachments"`
	Sources     []source.Summary        `json:"sources"`
}

type ManagedSkillDocument struct {
	Artifact artifact.Artifact      `json:"artifact"`
	Document document.SkillDocument `json:"document"`
}

// BuiltInBundleTopology is trusted bootstrap input. SkillStoreWrapper does not
// expose this type through Wails.
type BuiltInBundleTopology struct {
	RootID                root.RootID                            `json:"rootId"`
	CollectionID          collection.CollectionID                `json:"collectionId"`
	SourceID              source.SourceID                        `json:"sourceId"`
	LogicalName           basespec.LogicalName                   `json:"logicalName"`
	LogicalVersion        basespec.LogicalVersion                `json:"logicalVersion"`
	DisplayName           string                                 `json:"displayName"`
	Description           string                                 `json:"description"`
	Labels                map[string]string                      `json:"labels"`
	Enabled               bool                                   `json:"enabled"`
	DiscoveryRoot         basespec.Locator                       `json:"discoveryRoot"`
	ExpectedMemberDigests map[basespec.Locator]cryptoutil.Digest `json:"expectedMemberDigests"`
}

type BuiltInCollectionSkill struct {
	ArtifactID artifact.ArtifactID `json:"artifactId"`
	Member     basespec.Locator    `json:"member"`
	Enabled    bool                `json:"enabled"`
}
