package versionbundle

type Changelog struct {
	Component   string `json:"component" yaml:"component"`
	Description string `json:"description" yaml:"description"`
	Kind        string `json:"kind" yaml:"kind"`
}
