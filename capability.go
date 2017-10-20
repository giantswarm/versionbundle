package versionbundle

import (
	"reflect"

	"github.com/giantswarm/microerror"
)

type Capability struct {
	Bundles []Bundle `json:"bundles" yaml:"bundles"`
	Name    string   `json:"name" yaml:"name"`
}

func (c Capability) Validate() error {
	for _, b := range c.Bundles {
		err := b.Validate()
		if err != nil {
			return microerror.Mask(err)
		}
	}

	if c.Name == "" {
		return microerror.Maskf(invalidCapabilityError, "name must not be empty")
	}

	return nil
}

type SortCapabilitiesByName []Capability

func (c SortCapabilitiesByName) Len() int           { return len(c) }
func (c SortCapabilitiesByName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c SortCapabilitiesByName) Less(i, j int) bool { return c[i].Name < c[j].Name }

type ValidateCapabilities []Capability

func (c ValidateCapabilities) Validate() error {
	if c.hasDuplicates() {
		return microerror.Mask(duplicatedCapabilityError)
	}

	for _, capability := range c {
		err := capability.Validate()
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (c ValidateCapabilities) hasDuplicates() bool {
	for _, c1 := range c {
		var seen int

		for _, c2 := range c {
			if reflect.DeepEqual(c1, c2) {
				seen++

				if seen >= 2 {
					return true
				}
			}
		}
	}

	return false
}
