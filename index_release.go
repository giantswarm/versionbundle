package versionbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
)

const indexReleaseTimestampFormat = "2006-01-02T15:04:05.00Z"

type IndexRelease struct {
	Active      bool        `yaml:"active"`
	Authorities []Authority `yaml:"authorities"`
	Date        time.Time   `yaml:"date"`
	Version     string      `yaml:"version"`
}

// CompileReleases takes indexReleases and collected version bundles and
// compiles canonicalized Releases from them.
func CompileReleases(logger micrologger.Logger, indexReleases []IndexRelease, bundles []Bundle) ([]Release, error) {
	releases := buildReleases(logger, indexReleases, bundles)

	deduplicateReleaseChangelog(releases)

	return releases, nil
}

func buildReleases(logger micrologger.Logger, indexReleases []IndexRelease, bundles []Bundle) []Release {
	bundleCache := make(map[string]Bundle)

	// Create cache of bundles for quick lookup
	for _, b := range bundles {
		bundleCache[b.ID()] = b
	}

	var releases []Release

	for _, ir := range indexReleases {
		bundles, err := groupBundlesForIndexRelease(ir, bundleCache)
		if IsBundleNotFound(err) {
			logger.Log("level", "warning", "message", err.Error())
			continue
		}

		if err != nil {
			panic(err)
		}

		rc := ReleaseConfig{
			Bundles: bundles,
		}

		release, err := NewRelease(rc)
		if err != nil {
			logger.Log("level", "warning", "message", err.Error())
			continue
		}

		release.active = ir.Active
		release.deprecated = false
		release.timestamp = ir.Date.Format(indexReleaseTimestampFormat)
		release.version = ir.Version

		releases = append(releases, release)
	}

	return releases
}

func groupBundlesForIndexRelease(ir IndexRelease, bundles map[string]Bundle) ([]Bundle, error) {
	var groupedBundles []Bundle
	for _, a := range ir.Authorities {
		b, found := bundles[a.BundleID()]
		if !found {
			return nil, microerror.Maskf(bundleNotFoundError, "IndexRelease v%s contains Authority with bundle ID %s that cannot be found from collected version bundles.", ir.Version, a.BundleID())
		}
		groupedBundles = append(groupedBundles, b)
	}

	return groupedBundles, nil
}

// deduplicateReleaseChangelog removes duplicate changelog entries in
// consequtive release entries.
func deduplicateReleaseChangelog(releases []Release) []Release {
	if len(releases) < 2 {
		return releases
	}

	type LogState int
	const (
		New LogState = iota
		Removed
	)

	sort.Sort(SortReleasesByVersion(releases))

	prevChangelogs := make(map[string]LogState)
	for _, clog := range releases[0].Changelogs() {
		prevChangelogs[clog.String()] = New
	}

	// First one is always there.
	filteredReleases := append([]Release{}, releases[0])

	for _, r := range releases[1:] {
		curChangelogs := make(map[string]LogState)
		for _, clog := range r.Changelogs() {
			clogStr := clog.String()
			state, exists := prevChangelogs[clogStr]
			switch exists {
			case true:
				if state == New {
					curChangelogs[clogStr] = Removed
				}
				// r.Changelogs() returns a copy of changelogs so removal won't
				// mess iteration in this case.
				r.removeChangelogEntry(clog)
			case false:
				curChangelogs[clogStr] = New
			}
		}

		filteredReleases = append(filteredReleases, r)
		prevChangelogs = curChangelogs
	}

	return filteredReleases
}

// TODO define and implement validation rules
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

			if authority.Endpoint == nil {
				return microerror.Maskf(invalidReleaseError, "release %s authority %s doesn't have defined endpoint", release.Version, authority.Name)
			}

			if authority.Version == "" {
				return microerror.Maskf(invalidReleaseError, "release %s authority %s doesn't have defined version", release.Version, authority.Name)
			}
		}
	}
	return nil
}

func validateReleaseDates(indexReleases []IndexRelease) error {
	for _, release := range indexReleases {
		if release.Date.IsZero() {
			return microerror.Maskf(invalidReleaseError, "release %s has empty release date", release.Version)
		}
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
		authorities := make([]string, 0, len(release.Authorities))
		for _, a := range release.Authorities {
			authorities = append(authorities, a.BundleID())
		}

		sort.Strings(authorities)

		sha256Hash.Reset()
		sha256Hash.Write([]byte(strings.Join(authorities, ",")))

		hexHash := hex.EncodeToString(sha256Hash.Sum(nil))
		otherVer, exists = releaseChecksums[hexHash]
		if exists {
			return microerror.Maskf(invalidReleaseError, "duplicate release contents for versions %s and %s", otherVer, release.Version)
		}
		releaseChecksums[hexHash] = release.Version
	}

	return nil
}
