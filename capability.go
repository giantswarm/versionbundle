package versionbundle

import (
	"reflect"

	"github.com/giantswarm/microerror"
)

// Capability represents the data format being exposed by authorities like
// microservices or operators.
type Capability struct {
	// Bundles are the version bundles being exposed by an authority. Over the
	// curse of a lifetime of an authority bundles are added and removed based on
	// the implementation and support of exposed component versions.
	Bundles []Bundle `json:"bundles" yaml:"bundles"`
	// Name is the name of the authority, e.g. the name of the microservice or
	// operator.
	Name string `json:"name" yaml:"name"`
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
