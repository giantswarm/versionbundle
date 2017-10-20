package versionbundle

import (
	"github.com/giantswarm/microerror"
)

var duplicatedCapabilityError = microerror.New("duplicated capability")

// IsDuplicatedCapability asserts duplicatedCapabilityError.
func IsDuplicatedCapability(err error) bool {
	return microerror.Cause(err) == duplicatedCapabilityError
}

var invalidDependencyError = microerror.New("invalid dependency")

// IsInvalidDependency asserts invalidDependencyError.
func IsInvalidDependency(err error) bool {
	return microerror.Cause(err) == invalidDependencyError
}
