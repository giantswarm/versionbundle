package versionbundle

import (
	"net/url"
)

type Authority struct {
	Endpoint *url.URL `yaml:"endpoint"`
	Name     string   `yaml:"Name"`
	Provider string   `yaml:"Provider"`
	Version  string   `yaml:"version"`
}
