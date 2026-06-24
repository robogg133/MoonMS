package plugin

import (
	"regexp"
)

func validIdentifier(s string) (valid bool) {
	valid, err := regexp.MatchString(`^[a-z0-9_]+$`, s)
	if err != nil {
		return false
	}
	return valid
}
