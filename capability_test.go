package versionbundle

import (
	"testing"
	"time"
)

func Test_Capability_Validate(t *testing.T) {
	testCases := []struct {
		Capability   Capability
		ErrorMatcher func(err error) bool
	}{
		// Test 0 ensures that an empty capability is not valid.
		{
			Capability:   Capability{},
			ErrorMatcher: IsInvalidCapability,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Capability: Capability{
				Bundles: []Bundle{},
				Name:    "",
			},
			ErrorMatcher: IsInvalidCapability,
		},

		// Test 2 ensures that a valid capability does not throw an error.
		{
			Capability: Capability{
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
			ErrorMatcher: nil,
		},

		// Test 3 is the same as 2 but with multiple bundles.
		{
			Capability: Capability{
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
								Version: "1.5.0",
							},
							{
								Name:    "kube-dns",
								Version: "2.0.0",
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
						Version:    "0.2.0",
					},
				},
				Name: "kubernetes-operator",
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures that a capability with only having one deprecated bundle
		// is not valid.
		{
			Capability: Capability{
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
						Deprecated: true,
						Time:       time.Unix(10, 5),
						Version:    "0.1.0",
					},
				},
				Name: "kubernetes-operator",
			},
			ErrorMatcher: IsInvalidCapability,
		},

		// Test 5 is the same as 4 but with multiple bundles.
		{
			Capability: Capability{
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
						Deprecated: true,
						Time:       time.Unix(10, 5),
						Version:    "0.1.0",
					},
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
								Version: "1.5.0",
							},
							{
								Name:    "kube-dns",
								Version: "2.0.0",
							},
						},
						Dependencies: []Dependency{
							{
								Name:    "kubernetes",
								Version: "<= 1.7.x",
							},
						},
						Deprecated: true,
						Time:       time.Unix(10, 5),
						Version:    "0.2.0",
					},
				},
				Name: "kubernetes-operator",
			},
			ErrorMatcher: IsInvalidCapability,
		},

		// Test 6 ensures that deprecated bundles are allowed as soon as at least
		// one bundle is not deprecated.
		{
			Capability: Capability{
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
								Version: "1.5.0",
							},
							{
								Name:    "kube-dns",
								Version: "2.0.0",
							},
						},
						Dependencies: []Dependency{
							{
								Name:    "kubernetes",
								Version: "<= 1.7.x",
							},
						},
						Deprecated: true,
						Time:       time.Unix(10, 5),
						Version:    "0.2.0",
					},
				},
				Name: "kubernetes-operator",
			},
			ErrorMatcher: nil,
		},

		// Test 7 ensures that a bundles within a capability cannot have the same
		// version.
		{
			Capability: Capability{
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
								Version: "1.5.0",
							},
							{
								Name:    "kube-dns",
								Version: "2.0.0",
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
			ErrorMatcher: IsInvalidCapability,
		},
	}

	for i, tc := range testCases {
		err := tc.Capability.Validate()
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}
	}
}
