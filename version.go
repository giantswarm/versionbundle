package versionbundle

import "strconv"

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	if err != nil {
		return false
	}

	return true
}
