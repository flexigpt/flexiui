package bundle

import (
	"encoding/json"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
)

type CollectionData struct {
	SchemaVersion           string                  `json:"schemaVersion"`
	DiscoveryPolicyRevision string                  `json:"discoveryPolicyRevision"`
	LogicalName             basespec.LogicalName    `json:"logicalName"`
	LogicalVersion          basespec.LogicalVersion `json:"logicalVersion,omitempty"`
	Labels                  map[string]string       `json:"labels,omitempty"`
	ManagedSourceID         source.SourceID         `json:"managedSourceID,omitempty"`
}

func DecodeCollectionData(
	raw json.RawMessage,
) (CollectionData, error) {
	value, err := jsonutil.DecodeCanonicalObject[CollectionData](raw, basespec.MaxDefinitionBytes)
	if err != nil {
		return CollectionData{}, err
	}
	if err := value.Validate(); err != nil {
		return CollectionData{}, err
	}
	return value, nil
}

func EncodeCollectionData(
	value CollectionData,
) (json.RawMessage, error) {
	if err := value.Validate(); err != nil {
		return nil, err
	}
	return jsonutil.MarshalCanonicalObject(value, basespec.MaxDefinitionBytes)
}

func (value CollectionData) Validate() error {
	if value.SchemaVersion != artifactbuiltin.MCPSchemaVersion ||
		value.DiscoveryPolicyRevision != artifactbuiltin.DecoderRevision {
		return fmt.Errorf(
			"%w: invalid MCP Bundle Collection data",
			basespec.ErrInvalid,
		)
	}
	if err := value.LogicalName.Validate(); err != nil {
		return err
	}
	if err := value.LogicalVersion.Validate(true); err != nil {
		return err
	}
	if err := validateCollectionLabels(value.Labels); err != nil {
		return err
	}
	if value.ManagedSourceID != "" {
		if err := value.ManagedSourceID.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func validateCollectionLabels(values map[string]string) error {
	if len(values) > basespec.MaxLabels {
		return fmt.Errorf(
			"%w: MCP Bundle labels exceed %d entries",
			basespec.ErrInvalid,
			basespec.MaxLabels,
		)
	}
	for key, value := range values {
		if err := basespec.ValidateIdentifier(
			"MCP Bundle label key",
			key,
			basespec.MaxKindBytes,
		); err != nil {
			return err
		}
		if err := basespec.ValidateRequiredText(
			"MCP Bundle label value",
			value,
			basespec.MaxLabelValueBytes,
		); err != nil {
			return err
		}
	}
	return nil
}
