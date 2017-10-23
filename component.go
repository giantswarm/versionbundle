package versionbundle

import (
	"strings"

	"github.com/giantswarm/microerror"
)

type Component struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

func (c Component) Validate() error {
	if c.Name == "" {
		return microerror.Maskf(invalidComponentError, "name must not be empty")
	}

	if c.Version == "" {
		return microerror.Maskf(invalidComponentError, "version must not be empty")
	}

	versionSplit := strings.Split(c.Version, ".")
	if len(versionSplit) != 3 {
		return microerror.Maskf(invalidComponentError, "version format must be '<major>.<minor>.<patch>'")
	}

	if !isPositiveNumber(versionSplit[0]) {
		return microerror.Maskf(invalidComponentError, "major version must be positive number")
	}

	if !isPositiveNumber(versionSplit[1]) {
		return microerror.Maskf(invalidComponentError, "minor version must be positive number")
	}

	if !isPositiveNumber(versionSplit[2]) {
		return microerror.Maskf(invalidComponentError, "patch version must be positive number")
	}

	return nil
}
