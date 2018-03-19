package versionbundle

import (
	"reflect"
	"sort"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type AggregatorConfig struct {
	Logger micrologger.Logger
}

type Aggregator struct {
	logger micrologger.Logger
}

func NewAggregator(config AggregatorConfig) (*Aggregator, error) {
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	a := &Aggregator{
		logger: config.Logger,
	}

	return a, nil
}

// Aggregate merges version bundles based on dependencies each version bundle
// within the given version bundles define for their own components.
func (a *Aggregator) Aggregate(bundles []Bundle) ([][]Bundle, error) {
	if len(bundles) == 0 {
		return nil, nil
	}

	var aggregatedBundles [][]Bundle

	if len(bundles) == 1 {
		aggregatedBundles = append(aggregatedBundles, bundles)
		return aggregatedBundles, nil
	}

	for _, b1 := range bundles {
		newGroup := []Bundle{
			b1,
		}

		for _, b2 := range bundles {
			if reflect.DeepEqual(b1, b2) {
				continue
			}

			if a.bundlesConflictWithDependencies(b1, b2) {
				continue
			}

			if a.bundlesConflictWithDependencies(b2, b1) {
				continue
			}

			if a.containsBundleByName(newGroup, b2) {
				continue
			}

			newGroup = append(newGroup, b2)
		}

		sort.Sort(SortBundlesByVersion(newGroup))
		sort.Stable(SortBundlesByName(newGroup))

		if a.containsAggregatedBundle(aggregatedBundles, newGroup) {
			continue
		}

		if distinctCount(bundles) != len(newGroup) {
			continue
		}

		aggregatedBundles = append(aggregatedBundles, newGroup)
	}

	err := AggregatedBundles(aggregatedBundles).Validate()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return aggregatedBundles, nil
}

func (a *Aggregator) bundlesConflictWithDependencies(b1, b2 Bundle) bool {
	for _, d := range b1.Dependencies {
		for _, c := range b2.Components {
			if d.Name != c.Name {
				continue
			}

			if !d.Matches(c) {
				return true
			}
		}
	}

	return false
}

func (a *Aggregator) containsAggregatedBundle(list [][]Bundle, item []Bundle) bool {
	for _, grouped := range list {
		if reflect.DeepEqual(grouped, item) {
			return true
		}
	}

	return false
}

func (a *Aggregator) containsBundleByName(list []Bundle, item Bundle) bool {
	for _, b := range list {
		if b.Name == item.Name {
			return true
		}
	}

	return false
}

func distinctCount(list []Bundle) int {
	m := map[string]struct{}{}

	for _, b := range list {
		m[b.Name] = struct{}{}
	}

	return len(m)
}
