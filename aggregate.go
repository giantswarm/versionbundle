package versionbundle

func Aggregate(capabilities []Capability) (Aggregation, error) {
	var newAggregration Aggregation

	if len(capabilities) == 0 {
		return newAggregration, nil
	}

	if len(capabilities) == 1 {
		newAggregration.BundledCapabilities = append(newAggregration.BundledCapabilities, []Capability{capabilities[0]})
	}

	// TODO compute

	return newAggregration, nil
}
