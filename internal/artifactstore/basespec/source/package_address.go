package source

import (
	"fmt"
	"path"
	"strings"

	"github.com/flexigpt/flexigpt-app/internal/artifactstore/basespec"
)

// ManagedPackageAddress is the generic semantic address of one complete
// managed package.
//
// Artifact Store owns only this three-segment address shape:
//
//	<kind>/<name>/<version>
//
// Artifact families own the values of Kind, Name, Version, all primary file
// names, and all package-relative resource conventions.
type ManagedPackageAddress struct {
	Kind    basespec.PackageKind    `json:"kind"`
	Name    basespec.LogicalName    `json:"name"`
	Version basespec.LogicalVersion `json:"version"`
}

func NewManagedPackageAddress(
	kind basespec.PackageKind,
	name basespec.LogicalName,
	version basespec.LogicalVersion,
) (ManagedPackageAddress, error) {
	value := ManagedPackageAddress{
		Kind:    kind,
		Name:    name,
		Version: version,
	}
	if err := value.Validate(); err != nil {
		return ManagedPackageAddress{}, err
	}
	return value, nil
}

// ParseManagedPackageAddressDirectory decodes an address previously derived
// through Directory. It accepts no extra namespace or implementation segments.
func ParseManagedPackageAddressDirectory(
	directory basespec.Locator,
) (ManagedPackageAddress, error) {
	if err := basespec.ValidatePortableLocator(directory, false); err != nil {
		return ManagedPackageAddress{}, err
	}

	segments := strings.Split(string(directory), "/")
	if len(segments) != 3 {
		return ManagedPackageAddress{}, fmt.Errorf(
			"%w: managed package directory %q must contain kind, name, and version",
			basespec.ErrInvalid,
			directory,
		)
	}

	return NewManagedPackageAddress(
		basespec.PackageKind(segments[0]),
		basespec.LogicalName(segments[1]),
		basespec.LogicalVersion(segments[2]),
	)
}

func (a ManagedPackageAddress) Validate() error {
	if err := basespec.ValidatePackageKind(a.Kind); err != nil {
		return err
	}
	if err := basespec.ValidatePackageName(a.Name); err != nil {
		return err
	}
	return basespec.ValidatePackageVersion(a.Version)
}

// Directory returns the source-relative directory used by MapStore and normal
// filesystem users. It is derived from semantic package identity and never
// caller-supplied as an arbitrary directory.
func (a ManagedPackageAddress) Directory() (basespec.Locator, error) {
	if err := a.Validate(); err != nil {
		return "", err
	}
	value := basespec.Locator(path.Join(
		string(a.Kind),
		string(a.Name),
		string(a.Version),
	))
	if err := basespec.ValidatePortableLocator(value, false); err != nil {
		return "", err
	}
	return value, nil
}

// FileLocator returns a source-relative locator for one package-relative
// regular file.
func (a ManagedPackageAddress) FileLocator(
	relative basespec.Locator,
) (basespec.Locator, error) {
	if err := basespec.ValidatePortableLocator(relative, false); err != nil {
		return "", err
	}
	directory, err := a.Directory()
	if err != nil {
		return "", err
	}
	value := basespec.Locator(path.Join(
		string(directory),
		string(relative),
	))
	if err := basespec.ValidatePortableLocator(value, false); err != nil {
		return "", err
	}
	return value, nil
}
