package versionbundle

import (
	"time"
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

	return nil
}
