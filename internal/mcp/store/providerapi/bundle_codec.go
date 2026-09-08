package providerapi

import (
	"context"
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/schema"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
	mcpDomain "github.com/flexigpt/flexigpt-app/internal/mcp/store/domain"
)

type (
	BundleCodec struct{}
)

func NewBundleCodec() providerapi.SchemaCodec {
	return BundleCodec{}
}

func (BundleCodec) JSONSchema() []byte {
	return append([]byte(nil), artifactbuiltin.BundleV1JSONSchema...)
}

func (BundleCodec) Canonicalize(
	ctx context.Context,
	raw []byte,
) (schema.ParsedDocument, error) {
	if err := artifactbuiltin.CheckCodecContext(ctx); err != nil {
		return schema.ParsedDocument{}, err
	}
	value, canonical, err := parseBundle(raw)
	if err != nil {
		return schema.ParsedDocument{}, err
	}
	return schema.ParsedDocument{
		Key:    BundleCodec{}.Key(),
		Digest: value.Digest,
		Raw:    canonical,
	}, nil
}

func (BundleCodec) Key() schema.Key {
	return artifactbuiltin.MCPBundleSchemaKey
}

func parseBundle(
	raw []byte,
) (mcpDomain.BundleDocument, json.RawMessage, error) {
	value, err := jsonutil.DecodeCanonicalObject[mcpDomain.BundleDocument](raw, basespec.MaxDefinitionBytes)
	if err != nil {
		return mcpDomain.BundleDocument{}, nil, err
	}
	return mcpDomain.CanonicalizeBundle(value)
}
