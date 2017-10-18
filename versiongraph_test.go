package versiongraph

import (
	"reflect"
	"testing"
)

func Test_VersionGraph_Bundle(t *testing.T) {
	testCases := []struct {
		VersionGraph    Interface
		ExpectedBundles []Graph
		ErrorMatcher    func(err error) bool
	}{
		// Test 0 ensures that nil input results in nil output.
		{
			VersionGraph: &VersionGraph{
				Dependencies: nil,
			},
			ExpectedBundles: nil,
			ErrorMatcher:    nil,
		},

		// Test 1 ensures that an empty list of dependencies as input results in nil
		// output.
		{
			VersionGraph: &VersionGraph{
				Dependencies: []Dependency{},
			},
			ExpectedBundles: nil,
			ErrorMatcher:    nil,
		},

		// Test 2 ensures that a single dependency as input results in one graph
		// having one dependency.
		{
			VersionGraph: &VersionGraph{
				Dependencies: []Dependency{
					{
						Name:    "calico",
						Version: "1.0.0",

						Conflicts: nil,
						Requires:  nil,
					},
				},
			},
			ExpectedBundles: []Graph{
				{
					Dependencies: []Dependency{
						{
							Name:    "calico",
							Version: "1.0.0",

							Conflicts: nil,
							Requires:  nil,
						},
					},
				},
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures that a duplicated dependency as input results in an error.
		{
			VersionGraph: &VersionGraph{
				Dependencies: []Dependency{
					{
						Name:    "calico",
						Version: "1.0.0",

						Conflicts: nil,
						Requires:  nil,
					},
					{
						Name:    "calico",
						Version: "1.0.0",

						Conflicts: nil,
						Requires:  nil,
					},
				},
			},
			ExpectedBundles: nil,
			ErrorMatcher:    IsDuplicatedDependency,
		},
	}

	for i, tc := range testCases {
		if i != 3 {
			continue
		}
		bundles, err := tc.VersionGraph.Bundle()
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else {
			if !reflect.DeepEqual(bundles, tc.ExpectedBundles) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedBundles, bundles)
			}
		}
	}
}
