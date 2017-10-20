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
	versionSplit := strings.Split(c.Version, ".")
	if len(versionSplit) != 3 {
		return microerror.Maskf(invalidDependencyError, "version format must be '<major>.<minor>.<patch>'")
	}

	if !isNumber(versionSplit[0]) {
		return microerror.Maskf(invalidDependencyError, "major version must be int")
	}

	if !isNumber(versionSplit[1]) {
		return microerror.Maskf(invalidDependencyError, "minor version must be int")
	}

	if !isNumber(versionSplit[2]) {
		return microerror.Maskf(invalidDependencyError, "patch version must be int")
	}

	return nil
}
