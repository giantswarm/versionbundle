package versionbundle

import "github.com/giantswarm/microerror"

type kind string

const (
	KindAdded      kind = "added"
	KindChanged    kind = "changed"
	KindDeprecated kind = "deprecated"
	KindFixed      kind = "fixed"
	KindRemoved    kind = "removed"
	KindSecurity   kind = "security"
)

var (
	validKinds = []kind{
		KindAdded,
		KindChanged,
		KindDeprecated,
		KindFixed,
		KindRemoved,
		KindSecurity,
	}
)

type Changelog struct {
	Component   string `json:"component" yaml:"component"`
	Description string `json:"description" yaml:"description"`
	Kind        kind   `json:"kind" yaml:"kind"`
}

// TODO write tests
func (c Changelog) Validate() error {
	if c.Component == "" {
		return microerror.Maskf(invalidCapabilityError, "name must not be empty")
	}

	if c.Description == "" {
		return microerror.Maskf(invalidCapabilityError, "name must not be empty")
	}

	if c.Kind == "" {
		return microerror.Maskf(invalidDependencyError, "kind must not be empty")
	}
	var found bool
	for _, k := range validKinds {
		if c.Kind == k {
			found = true
		}
	}
	if !found {
		return microerror.Maskf(invalidDependencyError, "kind must be one of %#v", validKinds)
	}

	return nil
}
