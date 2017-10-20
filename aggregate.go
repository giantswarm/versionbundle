package versionbundle

import (
	"reflect"
	"sort"
	"strings"

	"github.com/giantswarm/microerror"
)

func Aggregate(capabilities []Capability) (Aggregation, error) {
	// TODO validate depdenency version
	// TODO dispatch validation to types (implement validation in types)

	if hasDuplicatedCapabilities(capabilities) {
		return Aggregation{}, microerror.Mask(duplicatedCapabilityError)
	}

	var newAggregration Aggregation

	if len(capabilities) == 0 {
		return newAggregration, nil
	}

	if len(capabilities) == 1 {
		newAggregration.Capabilities = append(newAggregration.Capabilities, []Capability{capabilities[0]})
		return newAggregration, nil
	}

	for _, c1 := range capabilities {
		newCapabilities := []Capability{
			c1,
		}

		for _, c2 := range capabilities {
			if reflect.DeepEqual(c1, c2) {
				continue
			}

			if capabilitiesConflictWithDependencies(c1, c2) {
				continue
			}

			if capabilitiesConflictWithDependencies(c2, c1) {
				continue
			}

			newCapabilities = append(newCapabilities, c2)
		}

		if containsAggregatedCapabilities(newAggregration.Capabilities, newCapabilities) {
			continue
		}

		if len(capabilities) != len(newCapabilities) {
			continue
		}

		newAggregration.Capabilities = append(newAggregration.Capabilities, newCapabilities)
	}

	return newAggregration, nil
}

func hasDuplicatedCapabilities(list []Capability) bool {
	for _, c1 := range list {
		var seen int

		for _, c2 := range list {
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

func capabilitiesConflictWithDependencies(c1, c2 Capability) bool {
	for _, b1 := range c1.Bundles {
		for _, d := range b1.Dependencies {
			for _, b2 := range c2.Bundles {
				for _, c := range b2.Components {
					if c.Name != d.Name {
						continue
					}

					if !versionRequirementMatches(c.Version, d.Version) {
						return true
					}
				}
			}
		}
	}

	return false
}

func containsAggregatedCapabilities(list [][]Capability, item []Capability) bool {
	for _, c := range list {
		sort.Sort(SortCapabilitiesByName(c))
		sort.Sort(SortCapabilitiesByName(item))

		if reflect.DeepEqual(c, item) {
			return true
		}
	}

	return false
}

func versionRequirementMatches(componentVersion string, dependencyVersion string) bool {
	split := strings.Split(dependencyVersion, " ")
	operator := split[0]
	dependencyVersion = split[1]
	componentVersion, dependencyVersion = alignWildcardVersion(componentVersion, dependencyVersion)

	// TODO support more operators
	if operator == "<=" {
		if componentVersion <= dependencyVersion {
			return true
		}
	}

	return false
}

func alignWildcardVersion(componentVersion, dependencyVersion string) (string, string) {
	i := strings.Index(dependencyVersion, "x")

	if i == -1 {
		return componentVersion, dependencyVersion
	}

	componentVersion = componentVersion[:i-1]
	dependencyVersion = dependencyVersion[:i-1]

	return componentVersion, dependencyVersion
}
