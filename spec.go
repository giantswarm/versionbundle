package versiongraph

type Dependency struct {
	Name    string
	Version string

	Conflicts []string
	Requires  []string
}

type Graph struct {
	Dependencies []Dependency
}

type Interface interface {
	Bundle() ([]Graph, error)
}
