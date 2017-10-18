package versiongraph

import (
	"reflect"

	"github.com/giantswarm/microerror"
)

type VersionGraph struct {
	Dependencies []Dependency
}

func (g *VersionGraph) Bundle() ([]Graph, error) {
	if len(g.Dependencies) == 0 {
		return nil, nil
	}

	if hasDuplicatedDependecy(g.Dependencies) {
		return nil, microerror.Mask(duplicatedDependencyError)
	}

	var bundles []Graph

	for _, d := range g.Dependencies {
		g := Graph{
			Dependencies: []Dependency{d},
		}

		bundles = append(bundles, g)
	}

	return bundles, nil
}

func hasDuplicatedDependecy(list []Dependency) bool {
	for _, d1 := range list {
		var seen int

		for _, d2 := range list {
			if reflect.DeepEqual(d1, d2) {
				seen++

				if seen >= 2 {
					return true
				}
			}
		}
	}

	return false
}

func containsDependency(list []Dependency, item Dependency) bool {
	for _, d := range list {
		if reflect.DeepEqual(d, item) {
			return true
		}
	}

	return false
}
