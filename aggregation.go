package versionbundle

type Aggregation struct {
	BundledCapabilities [][]Capability `json:"bundledCapabilities" yaml:"bundledCapabilities"`
}
