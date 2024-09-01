package gonfig

import (
	"strings"
)

var normaliseDelimiterSet = []string{" ", "-", "_"}

func normaliseDelimiters(input, delim string) (output string) {
	for _, target := range normaliseDelimiterSet {
		if target != delim {
			output = strings.ReplaceAll(input, target, delim)
		}
	}

	return output
}

func toFlagName(input string) string {
	return strings.ToLower(normaliseDelimiters(input, "-"))
}

func toEnvName(input string) string {
	return strings.ToUpper(normaliseDelimiters(input, "_"))
}
