package versionbundle

import (
	"reflect"
	"sort"

	"github.com/giantswarm/microerror"
)

// Aggregate merges capabilities based on dependencies version bundles within
// the given capabilities define for their components.
func Aggregate(capabilities []Capability) (Aggregation, error) {
	err := ValidateCapabilities(capabilities).Validate()
	if err != nil {
		return Aggregation{}, microerror.Mask(err)
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
		for _, b1 := range c1.Bundles {
			c1.Bundles = []Bundle{
				b1,
			}

			newCapabilities := []Capability{
				c1,
			}

			for _, c2 := range capabilities {
				for _, b2 := range c2.Bundles {
					c2.Bundles = []Bundle{
						b2,
					}

					if reflect.DeepEqual(c1, c2) {
						continue
					}

					if capabilitiesConflictWithDependencies(c1, c2) {
						continue
					}

					if capabilitiesConflictWithDependencies(c2, c1) {
						continue
					}

					if containsCapabitlityWithBundleName(newCapabilities, c2) {
						continue
					}

					newCapabilities = append(newCapabilities, c2)
				}
			}

			sort.Sort(SortCapabilitiesByName(newCapabilities))
			for _, c := range newCapabilities {
				sort.Sort(SortBundlesByVersion(c.Bundles))
			}

			if containsAggregatedCapabilities(newAggregration.Capabilities, newCapabilities) {
				continue
			}

			if len(capabilities) != len(newCapabilities) {
				continue
			}

			newAggregration.Capabilities = append(newAggregration.Capabilities, newCapabilities)
		}
	}

	err = newAggregration.Validate()
	if err != nil {
		return Aggregation{}, microerror.Mask(err)
	}

	return newAggregration, nil
}

func capabilitiesConflictWithDependencies(c1, c2 Capability) bool {
	for _, b1 := range c1.Bundles {
		for _, d := range b1.Dependencies {
			for _, b2 := range c2.Bundles {
				for _, c := range b2.Components {
					if d.Name != c.Name {
						continue
					}

					if !d.Matches(c) {
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
		if reflect.DeepEqual(c, item) {
			return true
		}
	}

	return false
}

func containsCapabitlityWithBundleName(list []Capability, item Capability) bool {
	for _, c := range list {
		if c.Name == item.Name {
			return true
		}
	}

	return false
}
