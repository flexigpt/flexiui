package providerapi

import (
	"context"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/collection"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/root"
)

type LifecycleOperation string

const (
	LifecycleCreateCollection         LifecycleOperation = "collection.create"
	LifecycleUpdateCollection         LifecycleOperation = "collection.update"
	LifecycleRetireCollection         LifecycleOperation = "collection.retire"
	LifecyclePurgeCollection          LifecycleOperation = "collection.purge"
	LifecycleAttachSource             LifecycleOperation = "attachment.attach"
	LifecycleUpdateAttachment         LifecycleOperation = "attachment.update"
	LifecycleDetachSource             LifecycleOperation = "attachment.detach"
	LifecycleReplaceAttachment        LifecycleOperation = "attachment.replace"
	LifecycleAdoptArtifact            LifecycleOperation = "artifact.adopt"
	LifecyclePinArtifact              LifecycleOperation = "artifact.pin"
	LifecycleSetArtifactEnabled       LifecycleOperation = "artifact.set-enabled"
	LifecycleSetArtifactName          LifecycleOperation = "artifact.set-name"
	LifecycleUpdateArtifactData       LifecycleOperation = "artifact.update-data"
	LifecycleUnadoptArtifact          LifecycleOperation = "artifact.unadopt"
	LifecyclePurgeArtifact            LifecycleOperation = "artifact.purge"
	LifecyclePurgeAndSuppressArtifact LifecycleOperation = "artifact.purge-suppress"
	LifecycleSuppressBinding          LifecycleOperation = "artifact.suppress"
	LifecycleUnsuppressBinding        LifecycleOperation = "artifact.unsuppress"
	LifecycleRefreshCollection        LifecycleOperation = "collection.refresh"
	LifecyclePublishManagedArtifact   LifecycleOperation = "managed-artifact.publish"
	LifecyclePublishManagedCollection LifecycleOperation = "managed-collection.publish"
	LifecycleRemoveManagedArtifact    LifecycleOperation = "managed-artifact.remove"
)

// LifecycleCommand identifies the aggregate whose lifecycle operation is
// about to be performed by Artifact Store.
//
// The command carries only identity and ownership information. Store-owned
// services retain validation and mutation of generic metadata.
type LifecycleCommand struct {
	Operation      LifecycleOperation
	RootID         root.RootID
	CollectionID   collection.CollectionID
	CollectionKind collection.CollectionKind
}

func (c LifecycleCommand) Validate() error {
	switch c.Operation {
	case LifecycleCreateCollection,
		LifecycleUpdateCollection,
		LifecycleRetireCollection,
		LifecyclePurgeCollection,
		LifecycleAttachSource,
		LifecycleUpdateAttachment,
		LifecycleDetachSource,
		LifecycleReplaceAttachment,
		LifecycleAdoptArtifact,
		LifecyclePinArtifact,
		LifecycleSetArtifactEnabled,
		LifecycleSetArtifactName,
		LifecycleUpdateArtifactData,
		LifecycleUnadoptArtifact,
		LifecyclePurgeArtifact,
		LifecyclePurgeAndSuppressArtifact,
		LifecycleSuppressBinding,
		LifecycleUnsuppressBinding,
		LifecycleRefreshCollection,
		LifecyclePublishManagedArtifact,
		LifecyclePublishManagedCollection,
		LifecycleRemoveManagedArtifact:
	default:
		return fmt.Errorf(
			"%w: unsupported provider lifecycle operation %q",
			basespec.ErrInvalid,
			c.Operation,
		)
	}
	if err := c.RootID.Validate(); err != nil {
		return err
	}
	if err := c.CollectionID.Validate(); err != nil {
		return err
	}
	return c.CollectionKind.Validate()
}

// CollectionLifecycle is the provider-owned validation hook for a registered
// Collection kind.
//
// Artifact Store calls this hook before invoking its own lifecycle service.
// Providers receive no mutation capability, repository, source runtime, or
// system Components reference.
type CollectionLifecycle interface {
	ValidateLifecycle(
		ctx context.Context,
		command LifecycleCommand,
	) error
}
