package versionbundle

import "github.com/giantswarm/microerror"

type Aggregation struct {
	Capabilities [][]Capability `json:"capabilities" yaml:"capabilities"`
}

// TODO write tests
func (a Aggregation) Validate() error {
	for _, bundle := range c.Capabilities {
		for _, c := range bundle {
			err := c.Validate()
			if err != nil {
				return microerror.Mask(err)
			}
		}
	}

	// TODO ensure same length

	return nil
}
