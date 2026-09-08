package collectiondata

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
	"github.com/flexigpt/flexigpt-app/internal/jsonutil"
	workspaceDomain "github.com/flexigpt/flexigpt-app/internal/workspace/store/domain"
)

func EncodeCollectionData(value workspaceDomain.CollectionData) (json.RawMessage, error) {
	if err := ValidateCollectionData(value); err != nil {
		return nil, err
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	canonical, err := jsonutil.CanonicalizeObject(
		raw,
		basespec.MaxLocalDataBytes,
	)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(canonical), nil
}

func DecodeCollectionData(raw json.RawMessage) (workspaceDomain.CollectionData, error) {
	canonical, err := jsonutil.CanonicalizeObject(
		raw,
		basespec.MaxLocalDataBytes,
	)
	if err != nil {
		return workspaceDomain.CollectionData{}, err
	}

	decoder := json.NewDecoder(bytes.NewReader(canonical))
	decoder.DisallowUnknownFields()
	var value workspaceDomain.CollectionData
	if err := decoder.Decode(&value); err != nil {
		return workspaceDomain.CollectionData{}, err
	}
	if err := ValidateCollectionData(value); err != nil {
		return workspaceDomain.CollectionData{}, err
	}
	return value, nil
}

func ValidateCollectionData(value workspaceDomain.CollectionData) error {
	if err := workspaceDomain.ValidateDiscoveryPreferences(value.Discovery); err != nil {
		return fmt.Errorf("%w: %w", workspaceDomain.ErrInvalidWorkspace, err)
	}
	return nil
}
