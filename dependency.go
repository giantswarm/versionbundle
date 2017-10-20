package versionbundle

import (
	"strings"

	"github.com/giantswarm/microerror"
)

const (
	OperatorEqual          = "=="
	OperatorGreaterOrEqual = ">="
	OperatorLessOrEqual    = "<="
	OperatorNotEqual       = "!="
)

var (
	validOperators = []string{
		OperatorEqual,
		OperatorGreaterOrEqual,
		OperatorLessOrEqual,
		OperatorNotEqual,
	}
)

type Dependency struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

// TODO write tests
func (d Dependency) Matches(c Component) bool {
	if d.Name != c.Name {
		return false
	}

	dependencyVersion, componentVersion, operator := parseParameters(d.Version, c.Version)

	if operator == OperatorEqual {
		if componentVersion == dependencyVersion {
			return true
		}
	}

	if operator == OperatorGreaterOrEqual {
		if componentVersion >= dependencyVersion {
			return true
		}
	}

	if operator == OperatorLessOrEqual {
		if componentVersion <= dependencyVersion {
			return true
		}
	}

	if operator == OperatorNotEqual {
		if componentVersion != dependencyVersion {
			return true
		}
	}

	return false
}

// TODO write tests
func (d Dependency) Validate() error {
	inputSplit := strings.Split(d.Version, " ")
	if len(inputSplit) != 2 {
		return microerror.Maskf(invalidDependencyError, "input format must be '<operator> <semver version>'")
	}

	operator := inputSplit[0]
	if operator == "" {
		return microerror.Maskf(invalidDependencyError, "operator must not be empty")
	}
	var found bool
	for _, o := range validOperators {
		if operator == o {
			found = true
		}
	}
	if !found {
		return microerror.Maskf(invalidDependencyError, "operator format must be one of %#v", validOperators)
	}

	versionSplit := strings.Split(inputSplit[1], ".")
	if len(versionSplit) != 3 {
		return microerror.Maskf(invalidDependencyError, "version format must be '<major>.<minor>.<patch>'")
	}

	if !isNumber(versionSplit[0]) {
		return microerror.Maskf(invalidDependencyError, "major version must be int")
	}

	minor := versionSplit[1]
	if !isNumber(minor) && minor != "x" {
		return microerror.Maskf(invalidDependencyError, "minor version must be int or wildcard ('x')")
	}

	patch := versionSplit[2]
	if !isNumber(patch) && patch != "x" {
		return microerror.Maskf(invalidDependencyError, "patch version must be int or wildcard ('x')")
	}

	return nil
}

func parseParameters(dependencyVersion, componentVersion string) (string, string, string) {
	split := strings.Split(dependencyVersion, " ")
	dependencyVersion = split[1]

	i := strings.Index(dependencyVersion, "x")
	if i == -1 {
		return dependencyVersion, componentVersion, split[0]
	}

	return dependencyVersion[:i-1], componentVersion[:i-1], split[0]
}
