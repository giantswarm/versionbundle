package versiongraph

import (
	"github.com/giantswarm/microerror"
)

var duplicatedDependencyError = microerror.New("duplicated dependency")

// IsDuplicatedDependency asserts duplicatedDependencyError.
func IsDuplicatedDependency(err error) bool {
	return microerror.Cause(err) == duplicatedDependencyError
}
