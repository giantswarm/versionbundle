package versionbundle

import (
	"reflect"
	"testing"
	"time"
)

func Test_Release_Components(t *testing.T) {
	testCases := []struct {
		Bundles            []Bundle
		ExpectedComponents []Component
		ErrorMatcher       func(err error) bool
	}{
		// Test 0 ensures creating a release with a nil slice of bundles throws
		// an error when creating a new release type.
		{
			Bundles:            nil,
			ExpectedComponents: nil,
			ErrorMatcher:       IsInvalidConfig,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Bundles:            []Bundle{},
			ExpectedComponents: nil,
			ErrorMatcher:       IsInvalidConfig,
		},

		// Test 2 ensures computing the release components when having a list
		// of one bundle given works as expected.
		{
			Bundles: []Bundle{
				{
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
					Name:    "kubernetes-operator",
					Version: "0.0.1",
				},
			},
			ExpectedComponents: []Component{
				{
					Name:    "calico",
					Version: "1.1.0",
				},
				{
					Name:    "kube-dns",
					Version: "1.0.0",
				},
				{
					Name:    "kubernetes-operator",
					Version: "0.0.1",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 3 is the same as 2 but with a different components.
		{
			Bundles: []Bundle{
				{
					Components: []Component{
						{
							Name:    "kube-dns",
							Version: "1.17.0",
						},
						{
							Name:    "calico",
							Version: "3.1.0",
						},
					},
					Name:    "kubernetes-operator",
					Version: "11.4.1",
				},
			},
			ExpectedComponents: []Component{
				{
					Name:    "calico",
					Version: "3.1.0",
				},
				{
					Name:    "kube-dns",
					Version: "1.17.0",
				},
				{
					Name:    "kubernetes-operator",
					Version: "11.4.1",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures computing the release components when having a list of
		// two bundles given works as expected.
		{
			Bundles: []Bundle{
				{
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
					Name:    "kubernetes-operator",
					Version: "0.1.0",
				},
				{
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
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
			},
			ExpectedComponents: []Component{
				{
					Name:    "calico",
					Version: "1.1.0",
				},
				{
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
				{
					Name:    "etcd",
					Version: "3.2.0",
				},
				{
					Name:    "kube-dns",
					Version: "1.0.0",
				},
				{
					Name:    "kubernetes",
					Version: "1.7.1",
				},
				{
					Name:    "kubernetes-operator",
					Version: "0.1.0",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 5 is like 4 but with version bundles being flipped.
		{
			Bundles: []Bundle{
				{
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
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
				{
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
					Name:    "kubernetes-operator",
					Version: "0.1.0",
				},
			},
			ExpectedComponents: []Component{
				{
					Name:    "calico",
					Version: "1.1.0",
				},
				{
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
				{
					Name:    "etcd",
					Version: "3.2.0",
				},
				{
					Name:    "kube-dns",
					Version: "1.0.0",
				},
				{
					Name:    "kubernetes",
					Version: "1.7.1",
				},
				{
					Name:    "kubernetes-operator",
					Version: "0.1.0",
				},
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		config := ReleaseConfig{
			Bundles: tc.Bundles,
		}

		r, err := NewRelease(config)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}

		c := r.Components()
		if !reflect.DeepEqual(c, tc.ExpectedComponents) {
			t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedComponents, c)
		}
	}
}

func Test_Releases_GetNewestRelease(t *testing.T) {
	testCases := []struct {
		Releases        []Release
		ExpectedRelease Release
		ErrorMatcher    func(err error) bool
	}{
		// Test 0 ensures that a nil list throws an execution failed error.
		{
			Releases:        nil,
			ExpectedRelease: Release{},
			ErrorMatcher:    IsExecutionFailed,
		},

		// Test 1 ensures that the newest release can be found.
		{
			Releases: []Release{
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
			},
			ExpectedRelease: Release{
				bundles: []Bundle{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
				version:   "0.1.0",
			},
			ErrorMatcher: nil,
		},

		// Test 2 is the same as 1 but with different releases.
		{
			Releases: []Release{
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
					version:   "0.2.0",
				},
			},
			ExpectedRelease: Release{
				bundles: []Bundle{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
				version:   "0.2.0",
			},
			ErrorMatcher: nil,
		},

		// Test 3 is the same as 1 but with different releases.
		{
			Releases: []Release{
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
					version:   "0.2.0",
				},
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
			},
			ExpectedRelease: Release{
				bundles: []Bundle{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
				version:   "0.2.0",
			},
			ErrorMatcher: nil,
		},

		// Test 4 is the same as 1 but with different releases.
		{
			Releases: []Release{
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
					version:   "0.2.0",
				},
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 40, 0, time.UTC),
					version:   "2.3.12",
				},
			},
			ExpectedRelease: Release{
				bundles: []Bundle{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 40, 0, time.UTC),
				version:   "2.3.12",
			},
			ErrorMatcher: nil,
		},

		// Test 5 is the same as 1 but with different releases.
		{
			Releases: []Release{
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
					version:   "0.2.0",
				},
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 40, 0, time.UTC),
					version:   "2.3.12",
				},
				{
					bundles: []Bundle{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
			},
			ExpectedRelease: Release{
				bundles: []Bundle{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 40, 0, time.UTC),
				version:   "2.3.12",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		result, err := GetNewestRelease(tc.Releases)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if !reflect.DeepEqual(result, tc.ExpectedRelease) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedRelease, result)
			}
		}
	}
}
