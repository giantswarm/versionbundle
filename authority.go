package versionbundle

import (
	"net/url"
)

type Authority struct {
	Endpoint *url.URL `yaml:"endpoint"`
	Name     string   `yaml:"Name"`
	Version  string   `yaml:"version"`
}
