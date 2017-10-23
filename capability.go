package versionbundle

import (
	"encoding/json"
	"reflect"

	"github.com/giantswarm/microerror"
)

type Capability struct {
	Bundles []Bundle `json:"bundles" yaml:"bundles"`
	Name    string   `json:"name" yaml:"name"`
}

func (c Capability) Copy() Capability {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	var copy Capability
	err = json.Unmarshal(b, &copy)
	if err != nil {
		panic(err)
	}

	return copy
}

func (c Capability) Validate() error {
	if len(c.Bundles) == 0 {
		return microerror.Maskf(invalidCapabilityError, "bundles must not be empty")
	}

	if c.Name == "" {
		return microerror.Maskf(invalidCapabilityError, "name must not be empty")
	}

	err := ValidateBundles(c.Bundles).Validate()
	if err != nil {
		return microerror.Maskf(invalidCapabilityError, err.Error())
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
		return microerror.Maskf(invalidCapabilityError, "capabilities must not be duplicated")
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

type ValidateBundledCapabilities [][]Capability

func (c ValidateBundledCapabilities) Validate() error {
	if c.hasDuplicates() {
		return microerror.Maskf(invalidCapabilityError, "capabilities must not be duplicated")
	}

	for _, capability := range c {
		err := ValidateCapabilities(capability).Validate()
		if err != nil {
			return microerror.Mask(err)
		}
	}

	return nil
}

func (c ValidateBundledCapabilities) hasDuplicates() bool {
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
