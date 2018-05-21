package versionbundle

import (
	"net/url"
	"testing"
	"time"
)

func urlMustParse(v string) *url.URL {
	u, err := url.Parse(v)
	if err != nil {
		panic(err)
	}

	return u
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
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
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
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
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
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
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
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1-duplicate",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: IsInvalidRelease,
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
