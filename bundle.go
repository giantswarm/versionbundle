package versionbundle

import (
	"strings"
	"time"

	"github.com/giantswarm/microerror"
)

type Bundle struct {
	Changelogs   []Changelog  `json:"changelogs" yaml:"changelogs"`
	Components   []Component  `json:"components" yaml:"components"`
	Dependencies []Dependency `json:"dependency" yaml:"dependency"`
	Time         time.Time    `json:"time" yaml:"time"`
	Version      string       `json:"version" yaml:"version"`
}

// TODO write tests
func (b Bundle) Validate() error {
	for _, c := range b.Changelogs {
		err := c.Validate()
		if err != nil {
			return microerror.Mask(err)
		}
	}

	for _, c := range b.Components {
		err := c.Validate()
		if err != nil {
			return microerror.Mask(err)
		}
	}

	for _, d := range b.Dependencies {
		err := d.Validate()
		if err != nil {
			return microerror.Mask(err)
		}
	}

	if b.Time == time.Time{} {
		return microerror.Maskf(invalidCapabilityError, "time must not be empty")
	}

	versionSplit := strings.Split(b.Version, ".")
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
