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

// TODO define and implement validation rules
func ValidateIndexReleases(indexReleases []IndexRelease) error {
	if len(indexReleases) == 0 {
		return nil
	}

	err := validateUniqueReleases(indexReleases)
	if err != nil {
		return microerror.Mask(err)
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
			n := strings.TrimSpace(a.Name)
			p := strings.TrimSpace(a.Provider)
			v := strings.TrimSpace(a.Version)
			authorities = append(authorities, n+":"+p+":"+v)
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
