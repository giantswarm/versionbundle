package versionbundle

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func Test_Bundles_Contains(t *testing.T) {
	testCases := []struct {
		Bundles        []Bundle
		Bundle         Bundle
		ExpectedResult bool
	}{
		// Test 0 ensures that a nil list and an empty bundle result in false.
		{
			Bundles:        nil,
			Bundle:         Bundle{},
			ExpectedResult: false,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Bundles:        []Bundle{},
			Bundle:         Bundle{},
			ExpectedResult: false,
		},

		// Test 2 ensures a list containing one version bundle and a matching
		// version bundle results in true.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
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
				Dependencies: []Dependency{},
				Deprecated:   false,
				Name:         "kubernetes-operator",
				Time:         time.Unix(10, 5),
				Version:      "0.1.0",
				WIP:          false,
			},
			ExpectedResult: true,
		},

		// Test 3 ensures a list containing two version bundle and a matching
		// version bundle results in true.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.2.0",
						},
						{
							Name:    "kube-dns",
							Version: "1.0.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.2.0",
					WIP:          false,
				},
			},
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.2.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Name:         "kubernetes-operator",
				Time:         time.Unix(10, 5),
				Version:      "0.2.0",
				WIP:          false,
			},
			ExpectedResult: true,
		},

		// Test 4 ensures a list containing one version bundle and a version bundle
		// that does not match results in false.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.2.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Name:         "kubernetes-operator",
				Time:         time.Unix(10, 5),
				Version:      "0.2.0",
				WIP:          false,
			},
			ExpectedResult: false,
		},

		// Test 5 ensures a list containing two version bundle and a version bundle
		// that does not match results in false.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.2.0",
						},
						{
							Name:    "kube-dns",
							Version: "1.0.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.2.0",
					WIP:          false,
				},
			},
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.3.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Name:         "kubernetes-operator",
				Time:         time.Unix(10, 5),
				Version:      "0.3.0",
				WIP:          false,
			},
			ExpectedResult: false,
		},
	}

	for i, tc := range testCases {
		result := Bundles(tc.Bundles).Contain(tc.Bundle)
		if result != tc.ExpectedResult {
			t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedResult, result)
		}
	}
}

func Test_Bundles_Copy(t *testing.T) {
	bundles := []Bundle{
		{
			Changelogs: []Changelog{},
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
			Name:       "kubernetes-operator",
			Time:       time.Unix(10, 5),
			Version:    "0.1.0",
			WIP:        false,
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
			Name:       "kubernetes-operator",
			Time:       time.Unix(20, 10),
			Version:    "0.0.9",
			WIP:        false,
		},
	}

	b1 := CopyBundles(bundles)
	b2 := CopyBundles(bundles)

	sort.Sort(SortBundlesByTime(b1))
	sort.Sort(SortBundlesByVersion(b2))

	if reflect.DeepEqual(b1, b2) {
		t.Fatalf("expected %#v got %#v", b1, b2)
	}
}

func Test_Bundles_GetBundlesByNameAndLabels(t *testing.T) {
	testCases := []struct {
		Bundles         []Bundle
		Labels          map[string]string
		Name            string
		ExpectedBundles []Bundle
		ErrorMatcher    func(err error) bool
	}{
		// Test 0 ensures that a nil list and an empty name throws an execution
		// failed error.
		{
			Bundles:         nil,
			Labels:          nil,
			Name:            "",
			ExpectedBundles: []Bundle{},
			ErrorMatcher:    IsExecutionFailed,
		},

		// Test 1 ensures that a nil list and a non-empty name throws an execution
		// failed error.
		{
			Bundles:         nil,
			Labels:          nil,
			Name:            "kubernetes-operator",
			ExpectedBundles: []Bundle{},
			ErrorMatcher:    IsExecutionFailed,
		},

		// Test 2 ensures that a non-empty list and an empty name throws an execution
		// failed error.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
			},
			Labels:          nil,
			Name:            "",
			ExpectedBundles: []Bundle{},
			ErrorMatcher:    IsExecutionFailed,
		},

		// Test 3 ensures that a non-empty list and an non-empty name throws a
		// not found errorn case the given name does not exist in the given list.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
			},
			Labels:          nil,
			Name:            "cert-operator",
			ExpectedBundles: []Bundle{},
			ErrorMatcher:    IsBundleNotFound,
		},

		// Test 4 is the same as 3 but with different version bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			Labels:          nil,
			Name:            "cert-operator",
			ExpectedBundles: []Bundle{},
			ErrorMatcher:    IsBundleNotFound,
		},

		// Test 5 ensures that a bundle can be found.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			Labels: nil,
			Name:   "cloud-config-operator",
			ExpectedBundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			ErrorMatcher: nil,
		},

		// Test 6 is the same as 5 but with different bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			Labels: nil,
			Name:   "cloud-config-operator",
			ExpectedBundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			ErrorMatcher: nil,
		},

		// Test 7 verifies that bundle can be found with name and labels.
		{
			Bundles: []Bundle{
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "aws",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(20, 15),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "azure",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(40, 35),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "kvm",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(10, 5),
					Version: "0.1.0",
					WIP:     false,
				},
			},
			Labels: map[string]string{
				"provider": "kvm",
			},
			Name: "cluster-operator",
			ExpectedBundles: []Bundle{
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "kvm",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(10, 5),
					Version: "0.1.0",
					WIP:     false,
				},
			},
			ErrorMatcher: nil,
		},

		// Test 8 verifies that bundleNotFoundError is returned when labels
		// don't match.
		{
			Bundles: []Bundle{
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "aws",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(20, 15),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "azure",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(40, 35),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "kvm",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(10, 5),
					Version: "0.1.0",
					WIP:     false,
				},
			},
			Labels: map[string]string{
				"provider": "bluemix",
			},
			Name:            "cluster-operator",
			ExpectedBundles: []Bundle{},
			ErrorMatcher:    IsBundleNotFound,
		},
	}

	for i, tc := range testCases {
		result, err := GetBundlesByNameAndLabels(tc.Bundles, tc.Name, tc.Labels)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if !reflect.DeepEqual(result, tc.ExpectedBundles) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedBundles, result)
			}
		}
	}
}

func Test_Bundles_GetNewestBundle(t *testing.T) {
	testCases := []struct {
		Bundles        []Bundle
		Labels         map[string]string
		ExpectedBundle Bundle
		ErrorMatcher   func(err error) bool
	}{
		// Test 0 ensures that a nil list throws an execution failed error.
		{
			Bundles:        nil,
			Labels:         nil,
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 1 ensures that the newest bundle can be found.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			Labels: nil,
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Name:         "cloud-config-operator",
				Time:         time.Unix(10, 5),
				Version:      "0.1.0",
				WIP:          false,
			},
			ErrorMatcher: nil,
		},

		// Test 2 is the same as 1 but with different bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(20, 15),
					Version:      "0.2.0",
					WIP:          false,
				},
			},
			Labels: nil,
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Name:         "cloud-config-operator",
				Time:         time.Unix(20, 15),
				Version:      "0.2.0",
				WIP:          false,
			},
			ErrorMatcher: nil,
		},

		// Test 3 is the same as 1 but with different bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(20, 15),
					Version:      "0.2.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			Labels: nil,
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Name:         "cloud-config-operator",
				Time:         time.Unix(20, 15),
				Version:      "0.2.0",
				WIP:          false,
			},
			ErrorMatcher: nil,
		},

		// Test 4 is the same as 1 but with different bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(20, 15),
					Version:      "0.2.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(40, 35),
					Version:      "2.3.12",
					WIP:          false,
				},
			},
			Labels: nil,
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Name:         "cloud-config-operator",
				Time:         time.Unix(40, 35),
				Version:      "2.3.12",
				WIP:          false,
			},
			ErrorMatcher: nil,
		},

		// Test 5 is the same as 1 but with different bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(20, 15),
					Version:      "0.2.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(40, 35),
					Version:      "2.3.12",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			Labels: nil,
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Name:         "cloud-config-operator",
				Time:         time.Unix(40, 35),
				Version:      "2.3.12",
				WIP:          false,
			},
			ErrorMatcher: nil,
		},

		// Test 6 verifies that correct bundle with given labels is found
		{
			Bundles: []Bundle{
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "aws",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(20, 15),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "azure",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(40, 35),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "kvm",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(10, 5),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "aws",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(120, 15),
					Version: "0.3.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "azure",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(140, 35),
					Version: "0.2.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "kvm",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(110, 5),
					Version: "0.4.0",
					WIP:     false,
				},
			},
			Labels: map[string]string{
				"provider": "azure",
			},
			ExpectedBundle: Bundle{
				Changelogs:   []Changelog{},
				Components:   []Component{},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Labels: map[string]string{
					"provider": "azure",
				},
				Name:    "cluster-operator",
				Time:    time.Unix(140, 35),
				Version: "0.2.0",
				WIP:     false,
			},
			ErrorMatcher: nil,
		},

		// Test 7 verifies that bundle with given labels is not found
		{
			Bundles: []Bundle{
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "aws",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(20, 15),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "azure",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(40, 35),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "kvm",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(10, 5),
					Version: "0.1.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "aws",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(120, 15),
					Version: "0.3.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "azure",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(140, 35),
					Version: "0.2.0",
					WIP:     false,
				},
				{
					Changelogs:   []Changelog{},
					Components:   []Component{},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Labels: map[string]string{
						"provider": "kvm",
					},
					Name:    "cluster-operator",
					Time:    time.Unix(110, 5),
					Version: "0.4.0",
					WIP:     false,
				},
			},
			Labels: map[string]string{
				"provider": "bluemix",
			},
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsBundleNotFound,
		},
	}

	for i, tc := range testCases {
		result, err := GetNewestBundle(tc.Bundles, tc.Labels)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if !reflect.DeepEqual(result, tc.ExpectedBundle) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedBundle, result)
			}
		}
	}
}

func Test_Bundles_Validate(t *testing.T) {
	testCases := []struct {
		Bundles      []Bundle
		ErrorMatcher func(err error) bool
	}{
		// Test 0 ensures that a nil list is not valid.
		{
			Bundles:      nil,
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Bundles:      []Bundle{},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 2 ensures validation of a list of version bundles where any version
		// bundle has no changelogs throws an error.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 3 is the same as 2 but with multiple version bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.2.0",
					WIP:        false,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 4 ensures validation of a list of version bundles where any version
		// bundle has no components throws an error.
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
					Components: []Component{},
					Dependencies: []Dependency{
						{
							Name:    "kubernetes",
							Version: "<= 1.7.x",
						},
					},
					Deprecated: false,
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 5 is the same as 4 but with multiple version bundles.
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
					Components: []Component{},
					Dependencies: []Dependency{
						{
							Name:    "kubernetes",
							Version: "<= 1.7.x",
						},
					},
					Deprecated: false,
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.2.0",
					WIP:        false,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 6 ensures validation of a list of version bundles where any version
		// bundle has no dependency does not throw an error.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			ErrorMatcher: nil,
		},

		// Test 7 is the same as 6 but with multiple version bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.2.0",
					WIP:        false,
				},
			},
			ErrorMatcher: nil,
		},

		// Test 8 ensures validation of a list of version bundles not having the
		// same name throws an error.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Name:       "ingress-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.2.0",
					WIP:        false,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 9 ensures validation of a list of version bundles having duplicated
		// version bundles throws an error.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 10 ensures validation of a list of version bundles having the same
		// version throws an error.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kube-dns",
							Description: "Kube-DNS version updated.",
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
							Version: "1.1.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(20, 10),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 11 ensures validation of a list of version bundles in which a newer
		// version bundle (time) has a lower version number throws an error.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kube-dns",
							Description: "Kube-DNS version updated.",
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
							Version: "1.1.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(20, 10),
					Version:      "0.0.9",
					WIP:          false,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 12 ensures validation of a list of version bundles where any version
		// bundle is deprecated and WIP throws an error.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   true,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          true,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 13 is the same as 12 but with multiple version bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kubernetes-operator",
					Time:         time.Unix(10, 5),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
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
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.2.0",
					WIP:        true,
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},
	}

	for i, tc := range testCases {
		err := Bundles(tc.Bundles).Validate()
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}
	}
}

func Test_isSubset(t *testing.T) {
	testCases := []struct {
		name           string
		subset         map[string]string
		superset       map[string]string
		expectedResult bool
	}{
		{
			name: "case 0: subset is subset of superset",
			subset: map[string]string{
				"foo": "1",
				"bar": "2",
				"baz": "3",
			},
			superset: map[string]string{
				"foo":   "1",
				"bar":   "2",
				"baz":   "3",
				"quux":  "4",
				"alice": "bob",
			},
			expectedResult: true,
		},
		{
			name: "case 1: subset equals superset",
			subset: map[string]string{
				"foo": "1",
				"bar": "2",
				"baz": "3",
			},
			superset: map[string]string{
				"foo": "1",
				"bar": "2",
				"baz": "3",
			},
			expectedResult: true,
		},
		{
			name: "case 2: subset has differing values from superset",
			subset: map[string]string{
				"foo": "1",
				"bar": "2",
				"baz": "3",
			},
			superset: map[string]string{
				"foo":   "0",
				"bar":   "2",
				"baz":   "3",
				"quux":  "4",
				"alice": "bob",
			},
			expectedResult: false,
		},
		{
			name: "case 3: subset has keys that are missing from superset",
			subset: map[string]string{
				"foo": "1",
				"bar": "2",
				"baz": "3",
			},
			superset: map[string]string{
				"foo":   "1",
				"baz":   "3",
				"quux":  "4",
				"alice": "bob",
			},
			expectedResult: false,
		},
		{
			name:   "case 4: subset is empty",
			subset: map[string]string{},
			superset: map[string]string{
				"foo":   "1",
				"baz":   "3",
				"quux":  "4",
				"alice": "bob",
			},
			expectedResult: true,
		},
		{
			name:   "case 5: subset is nil",
			subset: nil,
			superset: map[string]string{
				"foo":   "1",
				"baz":   "3",
				"quux":  "4",
				"alice": "bob",
			},
			expectedResult: true,
		},
		{
			name: "case 6: superset is empty",
			subset: map[string]string{
				"foo": "1",
				"bar": "2",
				"baz": "3",
			},
			superset:       map[string]string{},
			expectedResult: false,
		},
		{
			name: "case 7: subset is nil",
			subset: map[string]string{
				"foo": "1",
				"bar": "2",
				"baz": "3",
			},
			superset:       nil,
			expectedResult: false,
		},
		{
			name:           "case 8: subset and superset are empty",
			subset:         map[string]string{},
			superset:       map[string]string{},
			expectedResult: true,
		},

		{
			name:           "case 9: subset and superset are nil",
			subset:         nil,
			superset:       nil,
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret := isSubset(tc.subset, tc.superset)

			if ret != tc.expectedResult {
				t.Fatalf("isSubset() == %v, want %v", ret, tc.expectedResult)
			}
		})
	}
}
