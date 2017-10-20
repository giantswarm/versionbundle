package versionbundle

import (
	"github.com/giantswarm/microerror"
)

var invalidAggregationError = microerror.New("invalid aggregation")

// IsInvalidAggregationError asserts invalidAggregationError.
func IsInvalidAggregationError(err error) bool {
	return microerror.Cause(err) == invalidAggregationError
}

var invalidBundleError = microerror.New("invalid bundle")

// IsInvalidBundleError asserts invalidBundleError.
func IsInvalidBundleError(err error) bool {
	return microerror.Cause(err) == invalidBundleError
}

var invalidCapabilityError = microerror.New("invalid capability")

// IsInvalidCapability asserts invalidCapabilityError.
func IsInvalidCapability(err error) bool {
	return microerror.Cause(err) == invalidCapabilityError
}

var invalidChangelogError = microerror.New("invalid changelog")

// IsInvalidChangelog asserts invalidChangelogError.
func IsInvalidChangelog(err error) bool {
	return microerror.Cause(err) == invalidChangelogError
}

var invalidComponentError = microerror.New("invalid component")

// IsInvalidComponent asserts invalidComponentError.
func IsInvalidComponent(err error) bool {
	return microerror.Cause(err) == invalidComponentError
}

var invalidDependencyError = microerror.New("invalid dependency")

// IsInvalidDependency asserts invalidDependencyError.
func IsInvalidDependency(err error) bool {
	return microerror.Cause(err) == invalidDependencyError
}
