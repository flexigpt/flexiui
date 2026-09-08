package providerapi

import "github.com/flexigpt/flexigpt-app/internal/artifactstore/providerapi"

const artifactProviderName = "agent-skill"

// Provider registers the Agent Skill Artifact Store plugin.
type Provider struct {
	descriptor providerapi.Descriptor
}

func NewProvider() (*Provider, error) {
	decoder := NewDecoder()
	descriptor := providerapi.Descriptor{
		Name: artifactProviderName,
		CollectionBehaviors: []providerapi.CollectionBehavior{
			NewCollectionBehavior(),
		},
		Schemas: []providerapi.SchemaCodec{
			NewShareableCodec(),
		},
		Decoders: []providerapi.Decoder{
			decoder,
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
