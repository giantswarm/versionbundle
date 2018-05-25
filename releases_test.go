package versionbundle

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_deduplicateReleases(t *testing.T) {
	testCases := []struct {
		name           string
		input          []Release
		expectedOutput []Release
	}{
		{
			name: "case 0: test single Release",
			input: []Release{
				{
					timestamp: time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:   "1.0.1",
				},
			},
			expectedOutput: []Release{
				{
					timestamp: time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:   "1.0.1",
				},
			},
		},
		{
			name: "case 1: test with two Releases",
			input: []Release{
				{
					timestamp: time.Date(2018, time.February, 8, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
				},
				{
					timestamp:  time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:    "1.0.1",
					deprecated: true,
				},
			},
			expectedOutput: []Release{
				{
					timestamp:  time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: time.Date(2018, time.February, 8, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
				},
			},
		},
		{
			name: "case 2: test with two Releases with same version",
			input: []Release{
				{
					timestamp: time.Date(2018, time.February, 8, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
				},
				{
					timestamp:  time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:    "1.1.1",
					deprecated: true,
				},
			},
			expectedOutput: []Release{
				{
					timestamp: time.Date(2018, time.February, 8, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
				},
			},
		},
		{
			name: "case 3: test with three same versions in multiple Releases",
			input: []Release{
				{
					timestamp:  time.Date(2018, time.February, 8, 0, 0, 00, 0, time.UTC),
					version:    "1.1.1",
					deprecated: true,
				},
				{
					timestamp: time.Date(2018, time.March, 22, 0, 0, 00, 0, time.UTC),
					version:   "1.2.0",
					wip:       true,
				},
				{
					timestamp:  time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: time.Date(2018, time.February, 18, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
				},
				{
					timestamp:  time.Date(2018, time.February, 12, 0, 0, 00, 0, time.UTC),
					version:    "1.1.1",
					deprecated: true,
				},
			},
			expectedOutput: []Release{
				{
					timestamp:  time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: time.Date(2018, time.February, 18, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
				},
				{
					timestamp: time.Date(2018, time.March, 22, 0, 0, 00, 0, time.UTC),
					version:   "1.2.0",
					wip:       true,
				},
			},
		},
		{
			name: "case 4: test with two Releases with same version where older is active",
			input: []Release{
				{
					active:    false,
					timestamp: time.Date(2018, time.February, 8, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
				},
				{
					active:    true,
					timestamp: time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
				},
			},
			expectedOutput: []Release{
				{
					active:    true,
					timestamp: time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
				},
			},
		},
		{
			name: "case 5: test with three same versions in multiple Releases where middle one of duplicates is active",
			input: []Release{
				{
					timestamp:  time.Date(2018, time.February, 8, 0, 0, 00, 0, time.UTC),
					version:    "1.1.1",
					deprecated: true,
					active:     false,
				},
				{
					timestamp: time.Date(2018, time.March, 22, 0, 0, 00, 0, time.UTC),
					version:   "1.2.0",
					wip:       true,
				},
				{
					timestamp:  time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: time.Date(2018, time.February, 18, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
					wip:       true,
					active:    false,
				},
				{
					timestamp: time.Date(2018, time.February, 12, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
					active:    true,
				},
			},
			expectedOutput: []Release{
				{
					timestamp:  time.Date(2018, time.February, 2, 0, 0, 00, 0, time.UTC),
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: time.Date(2018, time.February, 12, 0, 0, 00, 0, time.UTC),
					version:   "1.1.1",
					active:    true,
				},
				{
					timestamp: time.Date(2018, time.March, 22, 0, 0, 00, 0, time.UTC),
					version:   "1.2.0",
					wip:       true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := deduplicateReleases(tc.input)

			if !reflect.DeepEqual(output, tc.expectedOutput) {
				fmt.Printf("%s: This is what I got: \n", tc.name)
				for _, r := range output {
					fmt.Printf("%s: %s, deprecated: %v, wip: %v, active: %v\n", r.Version(), r.Timestamp(), r.deprecated, r.wip, r.Active())
				}

				fmt.Printf("%s: This is what I want: \n", tc.name)
				for _, r := range tc.expectedOutput {
					fmt.Printf("%s: %s, deprecated: %v, wip: %v, active: %v\n", r.Version(), r.Timestamp(), r.deprecated, r.wip, r.Active())
				}

				t.Fatalf("failed.")
			}
		})
	}
}
