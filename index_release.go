package versionbundle

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
	"time"

	"github.com/giantswarm/microerror"
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
func CompileReleases(indexReleases []IndexRelease, bundles []Bundle) ([]Release, error) {
	releases, err := buildReleases(indexReleases, bundles)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// TODO(tuomas): Sort releases.
	// TODO(tuomas): Implement changelog cleanup etc.

	return releases, nil
}

func buildReleases(indexReleases []IndexRelease, bundles []Bundle) ([]Release, error) {
	bundleCache := make(map[string]Bundle)

	// Create cache of bundles for quick lookup
	for _, b := range bundles {
		bundleCache[b.ID()] = b
	}

	var releases []Release

	for _, ir := range indexReleases {
		release := Release{
			active:    ir.Active,
			timestamp: ir.Date.Format(releaseTimestampFormat),
			version:   ir.Version,
		}

		for _, a := range ir.Authorities {
			b, found := bundleCache[a.BundleID()]
			if !found {
				return nil, microerror.Maskf(invalidReleaseError, "IndexRelease v%s contains Authority with bundle ID %s that cannot be found from collected version bundles.")
			}
			release.bundles = append(release.bundles, b)
		}

		releases = append(releases, release)
	}

	return releases, nil
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