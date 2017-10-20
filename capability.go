package versionbundle

type Capability struct {
	Bundles []Bundle `json:"bundles" yaml:"bundles"`
	Name    string   `json:"name" yaml:"name"`
}
