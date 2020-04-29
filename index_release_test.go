package versionbundle

import (
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/micrologger/microloggertest"
)

func Test_buildReleases(t *testing.T) {
	testCases := []struct {
		name             string
		indexReleases    []IndexRelease
		bundles          []Bundle
		expectedReleases []Release
		errorMatcher     func(error) bool
	}{
		{
			name: "case 0: build one release",
			indexReleases: []IndexRelease{
				{
					Active: true,
					Authorities: []Authority{
						{
							Name:    "app-controller",
							Version: "0.1.0",
						},
						{
							Name:    "appcatalog-controller",
							Version: "0.1.0",
						},
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "1.0.0",
				},
			},
			bundles: []Bundle{
				{
					Name:    "cert-operator",
					Version: "0.1.0",
				},
				{
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.1.0",
				},
				{
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.2.0",
				},
				{
					Name:     "cluster-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
				{
					Name:    "kvm-operator",
					Version: "1.2.0",
				},
				{
					Name:    "kvm-operator",
					Version: "1.4.2",
				},
				{
					Name:    "kvm-operator",
					Version: "2.2.1",
				},
				{
					Name:    "app-controller",
					Version: "0.1.0",
				},
				{
					Name:    "appcatalog-controller",
					Version: "0.1.0",
				},
			},
			expectedReleases: []Release{
				{
					active: true,
					bundles: []Bundle{
						{
							Name:    "app-controller",
							Version: "0.1.0",
						},
						{
							Name:    "appcatalog-controller",
							Version: "0.1.0",
						},
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					components: []Component{
						{
							Name:    "app-controller",
							Version: "0.1.0",
						},
						{
							Name:    "appcatalog-controller",
							Version: "0.1.0",
						},
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:    "cluster-operator",
							Version: "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					timestamp: time.Date(2018, time.April, 16, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: build two releases",
			indexReleases: []IndexRelease{
				{
					Active: true,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "1.0.0",
				},
				{
					Active: true,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 22, 12, 00, 0, 0, time.UTC),
					Version: "1.1.0",
				},
			},
			bundles: []Bundle{
				{
					Name:    "cert-operator",
					Version: "0.1.0",
				},
				{
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.1.0",
				},
				{
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.2.0",
				},
				{
					Name:     "cluster-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
				{
					Name:    "kvm-operator",
					Version: "1.2.0",
				},
				{
					Name:    "kvm-operator",
					Version: "1.4.2",
				},
				{
					Name:    "kvm-operator",
					Version: "2.2.1",
				},
			},
			expectedReleases: []Release{
				{
					active: true,
					bundles: []Bundle{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					components: []Component{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:    "cluster-operator",
							Version: "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					timestamp: time.Date(2018, time.April, 16, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
				{
					active: true,
					bundles: []Bundle{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					components: []Component{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:    "cluster-operator",
							Version: "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					timestamp: time.Date(2018, time.April, 22, 12, 0, 0, 0, time.UTC),
					version:   "1.1.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: try to build two release but miss one bundle for second one",
			indexReleases: []IndexRelease{
				{
					Active: true,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "1.0.0",
				},
				{
					Active: true,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.4.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 22, 12, 00, 0, 0, time.UTC),
					Version: "1.1.0",
				},
			},
			bundles: []Bundle{
				{
					Name:    "cert-operator",
					Version: "0.1.0",
				},
				{
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.1.0",
				},
				{
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.2.0",
				},
				{
					Name:     "cluster-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
				{
					Name:    "kvm-operator",
					Version: "1.2.0",
				},
				{
					Name:    "kvm-operator",
					Version: "1.4.2",
				},
				{
					Name:    "kvm-operator",
					Version: "2.2.1",
				},
			},
			expectedReleases: []Release{
				{
					active: true,
					bundles: []Bundle{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					components: []Component{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:    "cluster-operator",
							Version: "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					timestamp: time.Date(2018, time.April, 16, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
			},
			errorMatcher: nil,
		},
	}

	logger := microloggertest.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			releases, err := buildReleases(logger, tc.indexReleases, tc.bundles)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if !reflect.DeepEqual(releases, tc.expectedReleases) {
				t.Fatalf("releases don't match expectedReleases; got:\n%#v\n\nexpected:\n%#v\n\n", releases, tc.expectedReleases)
			}
		})
	}
}

func Test_findPreviousRelease(t *testing.T) {
	testCases := []struct {
		name            string
		r0              Release
		releases        []Release
		expectedRelease Release
	}{
		{
			name: "case 0: return empty Release when releases is empty",
			r0: Release{
				timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
			},
			releases:        []Release{},
			expectedRelease: Release{},
		},
		{
			name: "case 1: return empty Release when releases contains only current release",
			r0: Release{
				timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
			},
			releases: []Release{
				{
					timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
				},
			},
			expectedRelease: Release{},
		},
		{
			name: "case 2: return correct release when releases contains two releases",
			r0: Release{
				timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
			},
			releases: []Release{
				{
					timestamp: time.Date(2018, time.May, 23, 12, 0, 0, 0, time.UTC),
				},
				{
					timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
				},
			},
			expectedRelease: Release{
				timestamp: time.Date(2018, time.May, 23, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "case 3: return correct release when releases contains two older releases",
			r0: Release{
				timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
			},
			releases: []Release{
				{
					timestamp: time.Date(2018, time.May, 18, 12, 0, 0, 0, time.UTC),
				},
				{
					timestamp: time.Date(2018, time.May, 23, 12, 0, 0, 0, time.UTC),
				},
				{
					timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
				},
				{
					timestamp: time.Date(2018, time.May, 26, 12, 0, 0, 0, time.UTC),
				},
			},
			expectedRelease: Release{
				timestamp: time.Date(2018, time.May, 23, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			release := findPreviousRelease(tc.r0, tc.releases)

			if !reflect.DeepEqual(release, tc.expectedRelease) {
				t.Fatalf("\ngot:\n%v\n\nexpected:\n%v\n\n", release, tc.expectedRelease)
			}
		})
	}
}

func Test_validateReleaseAuthority(t *testing.T) {
	testCases := []struct {
		name         string
		releases     []IndexRelease
		errorMatcher func(error) bool
	}{
		{
			name: "case 0: success with only single release",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: success with multiple releases",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: failure with single release containing one with nil authorities",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active:      false,
					Authorities: nil,
					Date:        time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version:     "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 3: failure with single release containing one with empty authorities list",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active:      false,
					Authorities: make([]Authority, 0),
					Date:        time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version:     "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 6: failure with single release containing one authority without name",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 5: failure with single release containing one authority without version",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateReleaseAuthorities(tc.releases)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}
		})
	}
}

func Test_validateReleaseDates(t *testing.T) {
	testCases := []struct {
		name         string
		releases     []IndexRelease
		errorMatcher func(error) bool
	}{
		{
			name: "case 0: success with only single release",
			releases: []IndexRelease{
				{
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "1.0.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: success with multiple unique releases",
			releases: []IndexRelease{
				{
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "3.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 20, 13, 12, 00, 00, time.UTC),
					Version: "2.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 19, 13, 12, 00, 00, time.UTC),
					Version: "1.0.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: failure with one release that has empty date",
			releases: []IndexRelease{
				{
					Date:    time.Time{},
					Version: "1.0.0",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 3: failure with multiple releases including one with empty date",
			releases: []IndexRelease{
				{
					Date:    time.Time{},
					Version: "4.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "3.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 20, 13, 12, 00, 00, time.UTC),
					Version: "2.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 19, 13, 12, 00, 00, time.UTC),
					Version: "1.0.0",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateReleaseDates(tc.releases)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}
		})
	}
}

func Test_validateUniqueReleases(t *testing.T) {
	testCases := []struct {
		name         string
		releases     []IndexRelease
		errorMatcher func(error) bool
	}{
		{
			name: "case 0: success with only single release",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: success with multiple unique releases",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: failure with multiple releases including one duplicate version",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 3: failure with multiple releases including one with duplicate version contents",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1-duplicate",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 4: failure with multiple releases including multiple duplicate version",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.5",
						},
						{
							Name:    "kvm-operator",
							Version: "2.3.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 5: success with multiple unique releases with app diff only",
			releases: []IndexRelease{
				{
					Active: false,
					Apps: []App{
						{
							App:              "nginx-ingress-controller",
							ComponentVersion: "0.30.0",
							Version:          "1.6.0",
						},
					},
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Apps: []App{
						{
							App:              "nginx-ingress-controller",
							ComponentVersion: "0.29.0",
							Version:          "1.5.0",
						},
					},
					Authorities: []Authority{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateUniqueReleases(tc.releases)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}
		})
	}
}
