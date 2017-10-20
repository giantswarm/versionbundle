package versionbundle

import (
	"github.com/giantswarm/microerror"
)

type Aggregation struct {
	Capabilities [][]Capability `json:"capabilities" yaml:"capabilities"`
}

func (a Aggregation) Validate() error {
	err := ValidateBundledCapabilities(a.Capabilities).Validate()
	if err != nil {
		return microerror.Mask(err)
	}

	for _, bundle := range a.Capabilities {
		for _, c := range bundle {
			err := c.Validate()
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	if len(a.Capabilities) != 0 {
		l := len(a.Capabilities[0])
		for _, bundle := range a.Capabilities {
			if l != len(bundle) {
				return microerror.Mask(invalidAggregationError)
			}
		}
	}

	return nil
}
