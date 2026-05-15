package main

import (
	"strings"
)

func sanitize(s []rune) string {

	var builder strings.Builder

	nextUpper := true
	for i := range s {
		switch s[i] {
		case '.':
			builder.WriteRune('_')
			nextUpper = true
			continue
		case '_':
			nextUpper = true
			continue
		}
		if nextUpper {
			s[i] = rune(strings.ToUpper(string(s[i]))[0])
			nextUpper = false
		}

		builder.WriteRune(s[i])
	}
	return builder.String()
}
