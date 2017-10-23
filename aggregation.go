package versionbundle

import (
	"github.com/giantswarm/microerror"
)

// Aggregation describes the aggregated version bundles which are wrapped by
// Capabilities. Each Capability must be within a merged Aggregation must only
// contain one version bundle. The set of capabilities being bundled within the
// same list represent the version bundles being able to work together based on
// the dependencies version bundles defined for their exposed components.
type Aggregation struct {
	Capabilities [][]Capability `json:"capabilities" yaml:"capabilities"`
}

func (a Aggregation) Validate() error {
	err := ValidateBundledCapabilities(a.Capabilities).Validate()
	if err != nil {
		return microerror.Maskf(invalidAggregationError, err.Error())
	}

	for _, capabilitiesList := range a.Capabilities {
		for _, c := range capabilitiesList {
			err := c.Validate()
			if err != nil {
				return microerror.Maskf(invalidAggregationError, err.Error())
			}
		}
	}

	if len(a.Capabilities) != 0 {
		l := len(a.Capabilities[0])
		for _, capabilitiesList := range a.Capabilities {
			if l != len(capabilitiesList) {
				return microerror.Mask(invalidAggregationError)
			}
		}
	}

	for _, capabilitiesList := range a.Capabilities {
		for _, c := range capabilitiesList {
			if len(c.Bundles) != 1 {
				return microerror.Maskf(invalidAggregationError, "there must be one capability bundle")
			}
		}
	}

	return nil
}
