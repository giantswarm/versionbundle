package versionbundle

import (
	"github.com/giantswarm/microerror"
)

var duplicatedCapabilityError = microerror.New("duplicated capability")

// IsDuplicatedCapability asserts duplicatedCapabilityError.
func IsDuplicatedCapability(err error) bool {
	return microerror.Cause(err) == duplicatedCapabilityError
}

var invalidCapabilityError = microerror.New("invalid capability")

// IsInvalidCapability asserts invalidCapabilityError.
func IsInvalidCapability(err error) bool {
	return microerror.Cause(err) == invalidCapabilityError
}

var invalidDependencyError = microerror.New("invalid dependency")

// IsInvalidDependency asserts invalidDependencyError.
func IsInvalidDependency(err error) bool {
	return microerror.Cause(err) == invalidDependencyError
}
