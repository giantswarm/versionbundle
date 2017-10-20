package versionbundle

import (
	"testing"
	"time"
)

func Test_Aggregation_Validate(t *testing.T) {
	testCases := []struct {
		Aggregation  Aggregation
		ErrorMatcher func(err error) bool
	}{
		// Test 0 ensures that an empty aggregation is valid.
		{
			Aggregation:  Aggregation{},
			ErrorMatcher: nil,
		},

		// Test 1 is the same as 0 but with an empty list of capabilities.
		{
			Aggregation: Aggregation{
				Capabilities: [][]Capability{},
			},
			ErrorMatcher: nil,
		},

		// Test 2 ensures that an aggregation with one bundled capabilities list is
		// valid.
		{
			Aggregation: Aggregation{
				Capabilities: [][]Capability{
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
									Deprecated: false,
									Time:       time.Unix(10, 5),
									Version:    "0.1.0",
								},
							},
							Name: "kubernetes-operator",
						},
					},
				},
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures that an aggregation with two bundled capabilities list is
		// valid.
		{
			Aggregation: Aggregation{
				Capabilities: [][]Capability{
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
									Deprecated: false,
									Time:       time.Unix(10, 5),
									Version:    "0.1.0",
								},
							},
							Name: "kubernetes-operator",
						},
						{
							Bundles: []Bundle{
								{
									Changelogs: []Changelog{
										{
											Component:   "etcd",
											Description: "Etcd version updated.",
											Kind:        "changed",
										},
										{
											Component:   "kubernetes",
											Description: "Kubernetes version updated.",
											Kind:        "changed",
										},
									},
									Components: []Component{
										{
											Name:    "etcd",
											Version: "3.2.0",
										},
										{
											Name:    "kubernetes",
											Version: "1.7.1",
										},
									},
									Dependencies: []Dependency{},
									Deprecated:   false,
									Time:         time.Unix(20, 15),
									Version:      "0.2.0",
								},
							},
							Name: "cloud-config-operator",
						},
					},
				},
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures that an aggregation with bundled capabilities lists having
		// different lengths is not valid.
		{
			Aggregation: Aggregation{
				Capabilities: [][]Capability{
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
									Deprecated: false,
									Time:       time.Unix(10, 5),
									Version:    "0.1.0",
								},
							},
							Name: "kubernetes-operator",
						},
						{
							Bundles: []Bundle{
								{
									Changelogs: []Changelog{
										{
											Component:   "etcd",
											Description: "Etcd version updated.",
											Kind:        "changed",
										},
										{
											Component:   "kubernetes",
											Description: "Kubernetes version updated.",
											Kind:        "changed",
										},
									},
									Components: []Component{
										{
											Name:    "etcd",
											Version: "3.2.0",
										},
										{
											Name:    "kubernetes",
											Version: "1.7.1",
										},
									},
									Dependencies: []Dependency{},
									Deprecated:   false,
									Time:         time.Unix(20, 15),
									Version:      "0.2.0",
								},
							},
							Name: "cloud-config-operator",
						},
					},
					{
						{
							Bundles: []Bundle{
								{
									Changelogs: []Changelog{
										{
											Component:   "etcd",
											Description: "Etcd version updated.",
											Kind:        "changed",
										},
										{
											Component:   "kubernetes",
											Description: "Kubernetes version updated.",
											Kind:        "changed",
										},
									},
									Components: []Component{
										{
											Name:    "etcd",
											Version: "3.2.0",
										},
										{
											Name:    "kubernetes",
											Version: "1.7.1",
										},
									},
									Dependencies: []Dependency{},
									Deprecated:   false,
									Time:         time.Unix(20, 15),
									Version:      "0.2.0",
								},
							},
							Name: "cloud-config-operator",
						},
					},
				},
			},
			ErrorMatcher: IsInvalidAggregationError,
		},

		// Test 5 ensures that an aggregation with bundled capabilities lists having
		// the same capabilities is not valid.
		{
			Aggregation: Aggregation{
				Capabilities: [][]Capability{
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
									Deprecated: false,
									Time:       time.Unix(10, 5),
									Version:    "0.1.0",
								},
							},
							Name: "kubernetes-operator",
						},
						{
							Bundles: []Bundle{
								{
									Changelogs: []Changelog{
										{
											Component:   "etcd",
											Description: "Etcd version updated.",
											Kind:        "changed",
										},
										{
											Component:   "kubernetes",
											Description: "Kubernetes version updated.",
											Kind:        "changed",
										},
									},
									Components: []Component{
										{
											Name:    "etcd",
											Version: "3.2.0",
										},
										{
											Name:    "kubernetes",
											Version: "1.7.1",
										},
									},
									Dependencies: []Dependency{},
									Deprecated:   false,
									Time:         time.Unix(20, 15),
									Version:      "0.2.0",
								},
							},
							Name: "cloud-config-operator",
						},
					},
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
									Deprecated: false,
									Time:       time.Unix(10, 5),
									Version:    "0.1.0",
								},
							},
							Name: "kubernetes-operator",
						},
						{
							Bundles: []Bundle{
								{
									Changelogs: []Changelog{
										{
											Component:   "etcd",
											Description: "Etcd version updated.",
											Kind:        "changed",
										},
										{
											Component:   "kubernetes",
											Description: "Kubernetes version updated.",
											Kind:        "changed",
										},
									},
									Components: []Component{
										{
											Name:    "etcd",
											Version: "3.2.0",
										},
										{
											Name:    "kubernetes",
											Version: "1.7.1",
										},
									},
									Dependencies: []Dependency{},
									Deprecated:   false,
									Time:         time.Unix(20, 15),
									Version:      "0.2.0",
								},
							},
							Name: "cloud-config-operator",
						},
					},
				},
			},
			ErrorMatcher: IsInvalidAggregationError,
		},
	}

	for i, tc := range testCases {
		err := tc.Aggregation.Validate()
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}
	}
}
