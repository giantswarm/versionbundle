package versionbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

type IndexRelease struct {
	Active      bool        `yaml:"active"`
	Apps        []App       `yaml:"apps"`
	Authorities []Authority `yaml:"authorities"`
	Date        time.Time   `yaml:"date"`
	Version     string      `yaml:"version"`
}

// CompileReleases takes indexReleases and collected version bundles and
// compiles canonicalized Releases from them.
func CompileReleases(logger micrologger.Logger, indexReleases []IndexRelease, bundles []Bundle) ([]Release, error) {
	releases, err := buildReleases(logger, indexReleases, bundles)
	if err != nil {
		return nil, err
	}

	releases = deduplicateReleaseChangelog(releases)

	return releases, nil
}

func buildReleases(logger micrologger.Logger, indexReleases []IndexRelease, bundles []Bundle) ([]Release, error) {
	bundleCache := make(map[string]Bundle)

	// Create cache of bundles for quick lookup
	for _, b := range bundles {
		bundleCache[b.ID()] = b
	}

	var releases []Release

	for _, ir := range indexReleases {
		bundles, err := groupBundlesForIndexRelease(ir, bundleCache)
		if IsBundleNotFound(err) {
			continue
		}

		if err != nil {
			return nil, err
		}

		rc := ReleaseConfig{
			Active:  ir.Active,
			Apps:    ir.Apps,
			Bundles: bundles,
			Date:    ir.Date,
			Version: ir.Version,
		}

		release, err := NewRelease(rc)
		if err != nil {
			logger.Log("level", "warning", "message", fmt.Sprintf("failed building new release from %s", ir.Version), "stack", fmt.Sprintf("%#v", err))
			continue
		}

		releases = append(releases, release)
	}

	return releases, nil
}

func groupBundlesForIndexRelease(ir IndexRelease, bundles map[string]Bundle) ([]Bundle, error) {
	var groupedBundles []Bundle
	for _, a := range ir.Authorities {
		b, found := bundles[a.BundleID()]
		if !found {
			return nil, microerror.Maskf(bundleNotFoundError, "IndexRelease %#q contains Authority with bundle ID %#q that cannot be found from collected version bundles.", ir.Version, a.Version)
		}
		groupedBundles = append(groupedBundles, b)
	}

	return groupedBundles, nil
}

// deduplicateReleaseChangelog removes duplicate changelog entries in
// consecutive release entries. Core concept of algorithm here is to first sort
// releases by their release version and then iterate them and compare current
// release to previous one that fulfills following requirements: smaller
// version number and earlier timestamp. Comparison of earlier timestamp is
// crucial here in order to calculate changelog correctly when newer patch
// releases have been introduced with lower version number
// (e.g. [1.0.0, 2.0.0] -> [1.0.0, 1.0.1, 2.0.0, 2.0.1]).
func deduplicateReleaseChangelog(releases []Release) []Release {
	if len(releases) < 2 {
		return releases
	}

	sort.Sort(SortReleasesByVersion(releases))

	return releases
}

// findPreviousRelease finds release that is older than argument r0. This
// function expects that releases is sorted by version as it is iterated
// backwards. If no previous release is found, empty Release is returned.
func findPreviousRelease(r0 Release, releases []Release) Release {
	for i := len(releases) - 1; i >= 0; i-- {
		if releases[i].timestamp.Before(r0.timestamp) {
			return releases[i]
		}
	}

	return Release{}
}

// ValidateIndexReleases ensures semantic rules for collection of indexReleases
// so that when used together, they form consistent and integral release index.
func ValidateIndexReleases(indexReleases []IndexRelease) error {
	if len(indexReleases) == 0 {
		return nil
	}

	var err error

	err = validateReleaseAuthorities(indexReleases)
	if err != nil {
		return microerror.Mask(err)
	}
	err = validateReleaseDates(indexReleases)
	if err != nil {
		return microerror.Mask(err)
	}
	err = validateUniqueReleases(indexReleases)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func validateReleaseAuthorities(indexReleases []IndexRelease) error {
	for _, release := range indexReleases {
		if len(release.Authorities) == 0 {
			return microerror.Maskf(invalidReleaseError, "release %s has no authorities", release.Version)
		}

		for _, authority := range release.Authorities {
			if authority.Name == "" {
				return microerror.Maskf(invalidReleaseError, "release %s contains authority without Name", release.Version)
			}

			if authority.Version == "" {
				return microerror.Maskf(invalidReleaseError, "release %s authority %s doesn't have defined version", release.Version, authority.Name)
			}
		}
	}
	return nil
}

func validateReleaseDates(indexReleases []IndexRelease) error {
	releaseDates := make(map[time.Time]string)
	for _, release := range indexReleases {
		if release.Date.IsZero() {
			return microerror.Maskf(invalidReleaseError, "release %s has empty release date", release.Version)
		}

		releaseDates[release.Date] = release.Version
	}

	return nil
}

func validateUniqueReleases(indexReleases []IndexRelease) error {
	releaseChecksums := make(map[string]string)
	releaseVersions := make(map[string]string)

	sha256Hash := sha256.New()

	for _, release := range indexReleases {
		// Verify release version number
		otherVer, exists := releaseVersions[release.Version]
		if exists {
			return microerror.Maskf(invalidReleaseError, "duplicate release versions %s and %s", otherVer, release.Version)
		}

		releaseVersions[release.Version] = release.Version

		// Verify release version contents
		appsAndAuthorities := make([]string, 0, len(release.Apps)+len(release.Authorities))
		for _, a := range release.Apps {
			appsAndAuthorities = append(appsAndAuthorities, a.AppID())
		}
		for _, a := range release.Authorities {
			appsAndAuthorities = append(appsAndAuthorities, a.BundleID())
		}

		sort.Strings(appsAndAuthorities)

		sha256Hash.Reset()
		_, err := sha256Hash.Write([]byte(strings.Join(appsAndAuthorities, ",")))
		if err != nil {
			return microerror.Mask(err)
		}

		hexHash := hex.EncodeToString(sha256Hash.Sum(nil))
		otherVer, exists = releaseChecksums[hexHash]
		if exists {
			return microerror.Maskf(invalidReleaseError, "duplicate release contents for versions %s and %s", otherVer, release.Version)
		}
		releaseChecksums[hexHash] = release.Version
	}

	return nil
}
