package versionbundle

type Capability struct {
	Bundles []Bundle `json:"bundles" yaml:"bundles"`
	Name    string   `json:"name" yaml:"name"`
}

type SortCapabilitiesByName []Capability

func (c SortCapabilitiesByName) Len() int           { return len(c) }
func (c SortCapabilitiesByName) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c SortCapabilitiesByName) Less(i, j int) bool { return c[i].Name < c[j].Name }
