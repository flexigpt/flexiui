package providerapi

import "github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"

const artifactProviderName = "mcp-bundle"

// Provider registers the MCP Bundle collection behavior, portable schemas,
// and bundle decoder as one Artifact Store plugin.
type Provider struct {
	descriptor providerapi.Descriptor
}

func NewProvider() (*Provider, error) {
	descriptor := providerapi.Descriptor{
		Name: artifactProviderName,
		CollectionBehaviors: []providerapi.CollectionBehavior{
			NewCollectionBehavior(),
		},
		Schemas: []providerapi.SchemaCodec{
			NewBundleCodec(),
			NewServerCodec(),
			NewPolicyCodec(),
		},
		Decoders: []providerapi.Decoder{
			NewDecoder(),
		},
	}
	if err := descriptor.Validate(); err != nil {
		return nil, err
	}

	return &Provider{
		descriptor: descriptor.Clone(),
	}, nil
}

func (p *Provider) Descriptor() providerapi.Descriptor {
	if p == nil {
		return providerapi.Descriptor{}
	}
	return p.descriptor.Clone()
}
