package basespec

type StorageKey string

func (v StorageKey) Validate() error {
	return ValidateIdentifier(
		"storage key",
		string(v),
		MaxStorageKeyBytes,
	)
}

type (
	DecoderID string

	LogicalName    string
	LogicalVersion string

	Locator            string
	SubresourceLocator string
)

func ValidatePackageName(value LogicalName) error {
	return ValidatePortableName("package name", string(value))
}

func ValidatePackageVersion(value LogicalVersion) error {
	return ValidatePortableName("package version", string(value))
}

func ValidateLogicalName(value LogicalName) error {
	return ValidateRequiredText(
		"logical name",
		string(value),
		MaxLogicalNameBytes,
	)
}

func ValidateLogicalVersion(value LogicalVersion, optional bool) error {
	if value == "" && optional {
		return nil
	}
	return ValidateRequiredText(
		"logical version",
		string(value),
		MaxVersionBytes,
	)
}

func ValidateDecoderID(value DecoderID) error {
	return ValidateIdentifier("decoder ID", string(value), MaxKindBytes)
}
