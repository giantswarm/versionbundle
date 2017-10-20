package versionbundle

import (
	"testing"
	"time"
)

func Test_Bundle_Validate(t *testing.T) {
	testCases := []struct {
		Bundle       Bundle
		ErrorMatcher func(err error) bool
	}{
		// Test 0 ensures that an empty capability is not valid.
		{
			Bundle:       Bundle{},
			ErrorMatcher: IsInvalidBundleError,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Bundle: Bundle{
				Changelogs:   []Changelog{},
				Components:   []Component{},
				Dependencies: []Dependency{},
				Deprecated:   false,
				Time:         time.Time{},
				Version:      "",
			},
			ErrorMatcher: IsInvalidBundleError,
		},
	}

	for i, tc := range testCases {
		err := tc.Bundle.Validate()
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}
	}
}
