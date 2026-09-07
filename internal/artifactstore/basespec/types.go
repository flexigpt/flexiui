package basespec

type StorageKey string

func (v StorageKey) Validate() error {
	return ValidateIdentifier(
		"storage key",
		string(v),
		MaxStorageKeyBytes,
	)
}

type LogicalName string

func (v LogicalName) Validate() error {
	return ValidateRequiredText(
		"logical name",
		string(v),
		MaxLogicalNameBytes,
	)
}

type LogicalVersion string

func (v LogicalVersion) Validate(optional bool) error {
	if v == "" && optional {
		return nil
	}
	return ValidateRequiredText(
		"logical version",
		string(v),
		MaxVersionBytes,
	)
}

type DecoderID string

func (v DecoderID) Validate() error {
	return ValidateIdentifier("decoder ID", string(v), MaxKindBytes)
}

type Locator string

func (v Locator) Validate(allowRoot bool) error {
	return validateRelativePath("locator", string(v), allowRoot)
}

type SubresourceLocator string

func (v SubresourceLocator) Validate() error {
	if v == "" {
		return nil
	}
	return validateRelativePath("subresource locator", string(v), false)
}
