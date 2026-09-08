package bundle

import (
	"github.com/flexigpt/agentskills-go/document"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/artifact"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/cryptoutil"
)

const (
	DiscoveryPolicyRevision = "skill.bundle.discovery.v1"

	RoleExternal collection.AttachmentRole = "external"
	RoleLibrary  collection.AttachmentRole = "library"
)

type AttachmentDraft struct {
	SourceID              source.SourceID
	Role                  collection.AttachmentRole
	Enabled               bool
	DiscoveryRoot         basespec.Locator
	ExpectedMemberDigests map[basespec.Locator]cryptoutil.Digest
}

type BuiltInBundleTopology struct {
	RootID                root.RootID
	CollectionID          collection.CollectionID
	SourceID              source.SourceID
	LogicalName           basespec.LogicalName
	LogicalVersion        basespec.LogicalVersion
	DisplayName           string
	Description           string
	Labels                map[string]string
	Enabled               bool
	DiscoveryRoot         basespec.Locator
	ExpectedMemberDigests map[basespec.Locator]cryptoutil.Digest
}

type ManagedSkillDocument struct {
	Artifact artifact.Artifact
	Document document.SkillDocument
}

type Bundle struct {
	Collection  collection.Collection
	Data        CollectionData
	Attachments []collection.Attachment
	Sources     []source.Summary
}
