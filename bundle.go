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
	Deprecated   bool         `json:"deprecated" yaml:"deprecated"`
	Time         time.Time    `json:"time" yaml:"time"`
	Version      string       `json:"version" yaml:"version"`
}

// TODO write tests
func (b Bundle) Validate() error {
	for _, c := range b.Changelogs {
		err := c.Validate()
		if err != nil {
			return microerror.Maskf(invalidBundleError, err.Error())
		}
	}

	for _, c := range b.Components {
		err := c.Validate()
		if err != nil {
			return microerror.Maskf(invalidBundleError, err.Error())
		}
	}

	for _, d := range b.Dependencies {
		err := d.Validate()
		if err != nil {
			return microerror.Maskf(invalidBundleError, err.Error())
		}
	}

	var emptyTime time.Time
	if b.Time == emptyTime {
		return microerror.Maskf(invalidBundleError, "time must not be empty")
	}

	versionSplit := strings.Split(b.Version, ".")
	if len(versionSplit) != 3 {
		return microerror.Maskf(invalidBundleError, "version format must be '<major>.<minor>.<patch>'")
	}

	if !isNumber(versionSplit[0]) {
		return microerror.Maskf(invalidBundleError, "major version must be int")
	}

	if !isNumber(versionSplit[1]) {
		return microerror.Maskf(invalidBundleError, "minor version must be int")
	}

	if !isNumber(versionSplit[2]) {
		return microerror.Maskf(invalidBundleError, "patch version must be int")
	}

	return nil
}

type ValidateBundles []Bundle

func (b ValidateBundles) Validate() error {
	if b.hasDuplicatedVersions() {
		return microerror.Mask(invalidBundleError)
	}

	for _, bundle := range b {
		err := bundle.Validate()
		if err != nil {
			return microerror.Maskf(invalidBundleError, err.Error())
		}
	}

	var deprecatedCount int
	for _, bundle := range b {
		if bundle.Deprecated {
			deprecatedCount++
		}
	}
	if deprecatedCount == len(b) {
		return microerror.Maskf(invalidBundleError, "at least one bundle must not be deprecated")
	}

	return nil
}

func (b ValidateBundles) hasDuplicatedVersions() bool {
	for _, b1 := range b {
		var seen int

		for _, b2 := range b {
			if b1.Version == b2.Version {
				seen++

				if seen >= 2 {
					return true
				}
			}
		}
	}

	return false
}
