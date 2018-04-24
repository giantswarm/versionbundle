package versionbundle

import (
	"fmt"
	"reflect"
	"testing"
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
					timestamp: "2018-02-02T00:00:00.000000Z",
					version:   "1.0.1",
				},
			},
			expectedOutput: []Release{
				{
					timestamp: "2018-02-02T00:00:00.000000Z",
					version:   "1.0.1",
				},
			},
		},
		{
			name: "case 1: test with two Releases",
			input: []Release{
				{
					timestamp: "2018-02-08T00:00:00.000000Z",
					version:   "1.1.1",
				},
				{
					timestamp:  "2018-02-02T00:00:00.000000Z",
					version:    "1.0.1",
					deprecated: true,
				},
			},
			expectedOutput: []Release{
				{
					timestamp:  "2018-02-02T00:00:00.000000Z",
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: "2018-02-08T00:00:00.000000Z",
					version:   "1.1.1",
				},
			},
		},
		{
			name: "case 2: test with two Releases with same version",
			input: []Release{
				{
					timestamp: "2018-02-08T00:00:00.000000Z",
					version:   "1.1.1",
				},
				{
					timestamp:  "2018-02-02T00:00:00.000000Z",
					version:    "1.1.1",
					deprecated: true,
				},
			},
			expectedOutput: []Release{
				{
					timestamp: "2018-02-08T00:00:00.000000Z",
					version:   "1.1.1",
				},
			},
		},
		{
			name: "case 3: test with three same versions in multiple Releases",
			input: []Release{
				{
					timestamp:  "2018-02-08T00:00:00.000000Z",
					version:    "1.1.1",
					deprecated: true,
				},
				{
					timestamp: "2018-03-22T00:00:00.000000Z",
					version:   "1.2.0",
					wip:       true,
				},
				{
					timestamp:  "2018-02-02T00:00:00.000000Z",
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: "2018-02-18T00:00:00.000000Z",
					version:   "1.1.1",
				},
				{
					timestamp:  "2018-02-12T00:00:00.000000Z",
					version:    "1.1.1",
					deprecated: true,
				},
			},
			expectedOutput: []Release{
				{
					timestamp:  "2018-02-02T00:00:00.000000Z",
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: "2018-02-18T00:00:00.000000Z",
					version:   "1.1.1",
				},
				{
					timestamp: "2018-03-22T00:00:00.000000Z",
					version:   "1.2.0",
					wip:       true,
				},
			},
		},
		{
			name: "case 4: test with two Releases with same version where older is active",
			input: []Release{
				{
					timestamp: "2018-02-08T00:00:00.000000Z",
					version:   "1.1.1",
					wip:       true,
				},
				{
					timestamp: "2018-02-02T00:00:00.000000Z",
					version:   "1.1.1",
				},
			},
			expectedOutput: []Release{
				{
					timestamp: "2018-02-02T00:00:00.000000Z",
					version:   "1.1.1",
				},
			},
		},
		{
			name: "case 5: test with three same versions in multiple Releases where middle one of duplicates is active",
			input: []Release{
				{
					timestamp:  "2018-02-08T00:00:00.000000Z",
					version:    "1.1.1",
					deprecated: true,
				},
				{
					timestamp: "2018-03-22T00:00:00.000000Z",
					version:   "1.2.0",
					wip:       true,
				},
				{
					timestamp:  "2018-02-02T00:00:00.000000Z",
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: "2018-02-18T00:00:00.000000Z",
					version:   "1.1.1",
					wip:       true,
				},
				{
					timestamp: "2018-02-12T00:00:00.000000Z",
					version:   "1.1.1",
				},
			},
			expectedOutput: []Release{
				{
					timestamp:  "2018-02-02T00:00:00.000000Z",
					version:    "1.0.1",
					deprecated: true,
				},
				{
					timestamp: "2018-02-12T00:00:00.000000Z",
					version:   "1.1.1",
				},
				{
					timestamp: "2018-03-22T00:00:00.000000Z",
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
