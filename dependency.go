package versionbundle

import "strings"

const (
	OperatorEqual          = "=="
	OperatorGreaterOrEqual = ">="
	OperatorLessOrEqual    = "<="
	OperatorNotEqual       = "!="
)

type Dependency struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

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

// TODO
func (d Dependency) Validate() error {
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
