package versionbundle

import (
	"reflect"
	"testing"
	"time"
)

func Test_Aggregate(t *testing.T) {
	testCases := []struct {
		Capabilities        []Capability
		ExpectedAggregation Aggregation
		ErrorMatcher        func(err error) bool
	}{
		// Test 0 ensures that nil input results in empty output.
		{
			Capabilities:        nil,
			ExpectedAggregation: Aggregation{},
			ErrorMatcher:        nil,
		},

		// Test 1 is the same as 0 but with an empty capabilities list.
		{
			Capabilities:        []Capability{},
			ExpectedAggregation: Aggregation{},
			ErrorMatcher:        nil,
		},

		// Test 2 ...
		{
			Capabilities: []Capability{
				{
					Bundles: []Bundle{
						{
							Changelogs: []Changelog{
								{
									Component:   "calico",
									Description: "Calico version updated.",
									Kind:        "changed",
								},
								{
									Component:   "kubernetes",
									Description: "Kubernetes version requirements changed due to calico update.",
									Kind:        "changed",
								},
							},
							Components: []Component{
								{
									Name:    "calico",
									Version: "1.1.0",
								},
								{
									Name:    "kube-dns",
									Version: "1.0.0",
								},
							},
							Dependencies: []Dependency{
								{
									Name:    "kubernetes",
									Version: "<= 1.7.x",
								},
							},
							Time:    time.Unix(10, 5),
							Version: "0.1.0",
						},
					},
					Name: "kubernetes-operator",
				},
			},
			ExpectedAggregation: Aggregation{
				BundledCapabilities: [][]Capability{
					{
						{
							Bundles: []Bundle{
								{
									Changelogs: []Changelog{
										{
											Component:   "calico",
											Description: "Calico version updated.",
											Kind:        "changed",
										},
										{
											Component:   "kubernetes",
											Description: "Kubernetes version requirements changed due to calico update.",
											Kind:        "changed",
										},
									},
									Components: []Component{
										{
											Name:    "calico",
											Version: "1.1.0",
										},
										{
											Name:    "kube-dns",
											Version: "1.0.0",
										},
									},
									Dependencies: []Dependency{
										{
											Name:    "kubernetes",
											Version: "<= 1.7.x",
										},
									},
									Time:    time.Unix(10, 5),
									Version: "0.1.0",
								},
							},
							Name: "kubernetes-operator",
						},
					},
				},
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		aggregation, err := Aggregate(tc.Capabilities)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else {
			if !reflect.DeepEqual(aggregation, tc.ExpectedAggregation) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedAggregation, aggregation)
			}
		}
	}
}
