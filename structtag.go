package gonfig

import (
	"reflect"
	"strconv"
	"strings"
)

var (
	// DefaultTag denotes the value to use if not set by other means. Must be
	// resolve to the same type as the value receiver
	DefaultTag = "default"
	// DelimiterTag denotes the separator used for slice types and enum options
	DelimiterTag = "delimiter"
	// EnumTag denotes list of choices to constrain the value by. Every element
	// must be resolvable to the same type as the value reciver
	EnumTag = "enum"
	// EnvTag denotes the name of the environment variable
	EnvTag = "env"
	// FlagTag denotes the name of the cli flag
	FlagTag = "flag"
	// RequiredTag denotes whether or not the value must be set
	RequiredTag = "required"
	// UsageTag denotes help text for use with cli flags
	UsageTag = "usage"
)

func getDefaultValue(t reflect.StructField) *string {
	if defaultValue, ok := t.Tag.Lookup(DefaultTag); ok {
		return &defaultValue
	}

	return nil
}

func getDelimiter(t reflect.StructField) string {
	delim, ok := t.Tag.Lookup(DelimiterTag)
	if ok {
		return delim
	}

	return ","
}

func getEnum(t reflect.StructField) ([]string, error) {
	enumTag, ok := t.Tag.Lookup(EnumTag)
	if !ok {
		return nil, nil
	}

	delim := getDelimiter(t)
	options := strings.Split(enumTag, delim)
	return options, nil
}

func getEnvName(t reflect.StructField) *string {
	if envName, ok := t.Tag.Lookup(EnvTag); ok {
		return &envName
	}

	return nil
}

func getFlagName(t reflect.StructField) *string {
	if flagName, ok := t.Tag.Lookup(FlagTag); ok {
		return &flagName
	}

	return nil
}

func getRequired(t reflect.StructField) (bool, error) {
	requiredTag, ok := t.Tag.Lookup(RequiredTag)
	if !ok {
		return false, nil
	}

	required, err := strconv.ParseBool(requiredTag)
	if err != nil {
		return false, err
	}

	return required, nil
}
