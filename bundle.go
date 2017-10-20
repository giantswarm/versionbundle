package versionbundle

import "time"

type Bundle struct {
	Changelogs   []Changelog  `json:"changelogs" yaml:"changelogs"`
	Components   []Component  `json:"components" yaml:"components"`
	Dependencies []Dependency `json:"dependency" yaml:"dependency"`
	Time         time.Time    `json:"time" yaml:"time"`
	Version      string       `json:"version" yaml:"version"`
}
