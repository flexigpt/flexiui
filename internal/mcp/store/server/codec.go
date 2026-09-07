package server

import (
	"context"
	"encoding/json"

	"github.com/flexigpt/flexigpt-app/internal/artifactbuiltin"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec/schema"
	"github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
)

type ServerCodec struct{}

func NewServerCodec() providerapi.SchemaCodec {
	return ServerCodec{}
}

func (ServerCodec) JSONSchema() []byte {
	return append([]byte(nil), artifactbuiltin.ServerV1JSONSchema...)
}

func (ServerCodec) Canonicalize(
	ctx context.Context,
	raw []byte,
) (schema.ParsedDocument, error) {
	if err := artifactbuiltin.CheckCodecContext(ctx); err != nil {
		return schema.ParsedDocument{}, err
	}
	value, canonical, err := parseServer(raw)
	if err != nil {
		return schema.ParsedDocument{}, err
	}
	return schema.ParsedDocument{
		Key:    ServerCodec{}.Key(),
		Digest: value.Digest,
		Raw:    canonical,
	}, nil
}

func (ServerCodec) Key() schema.Key {
	return artifactbuiltin.MCPServerSchemaKey
}

func parseServer(
	raw []byte,
) (ServerDocument, json.RawMessage, error) {
	value, err := jsonutil.DecodeCanonicalObject[ServerDocument](raw, basespec.MaxDefinitionBytes)
	if err != nil {
		return ServerDocument{}, nil, err
	}
	return CanonicalizeServer(value)
}
