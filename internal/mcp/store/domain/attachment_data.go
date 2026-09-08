package domain

import (
	"encoding/json"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/source"
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
)

type AttachmentData struct {
	SchemaVersion  string                       `json:"schemaVersion"`
	PackageAddress source.ManagedPackageAddress `json:"packageAddress"`
}

func DecodeAttachmentData(
	raw json.RawMessage,
) (AttachmentData, error) {
	value, err := jsonutil.DecodeCanonicalObject[AttachmentData](
		raw,
		basespec.MaxDefinitionBytes,
	)
	if err != nil {
		return AttachmentData{}, err
	}
	if err := value.Validate(); err != nil {
		return AttachmentData{}, err
	}
	return value, nil
}

func (d AttachmentData) Validate() error {
	if d.SchemaVersion != artifactbuiltin.MCPSchemaVersion {
		return fmt.Errorf("%w: invalid MCP Bundle attachment data schema", basespec.ErrInvalid)
	}
	return ValidateBundlePackageAddress(d.PackageAddress)
}

func EncodeAttachmentData(
	value AttachmentData,
) (json.RawMessage, error) {
	if err := value.Validate(); err != nil {
		return nil, err
	}
	if _, err := DocumentLocatorForPackage(value.PackageAddress); err != nil {
		return nil, err
	}
	return jsonutil.MarshalCanonicalObject(
		value,
		basespec.MaxDefinitionBytes,
	)
}

func (d AttachmentData) DocumentLocator() (basespec.Locator, error) {
	if err := d.Validate(); err != nil {
		return "", err
	}
	return DocumentLocatorForPackage(d.PackageAddress)
}
