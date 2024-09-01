package gonfig

import (
	"flag"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// FlagSet represents the minimum interface to be able to bind using
// field.BindFlag
type FlagSet interface {
	Var(flag.Value, string, string)
}

// Field contains the metadata extracted from struct field tags and the value
// setter
type Field struct {
	Name      string
	Usage     string
	FlagName  *string
	EnvName   *string
	Delimiter string
	Required  bool
	Enum      []string
	Value     Value
}

// NewField extracts and parses the struct tags of a field and creates a value
// setter based on the struct field value receiver
func NewField(t reflect.StructField, v reflect.Value) (*Field, error) {
	field := Field{
		Name:      t.Name,
		Usage:     t.Tag.Get(UsageTag),
		FlagName:  getFlagName(t),
		EnvName:   getEnvName(t),
		Delimiter: getDelimiter(t),
		Value:     NewValue(t, v),
	}

	if required, err := getRequired(t); err != nil {
		return nil, err
	} else {
		field.Required = required
	}

	if enum, err := getEnum(t); err != nil {
		return nil, err
	} else {
		field.Enum = enum
		// TODO validate each element of enum
	}

	if field.Value == nil {
		return nil, fmt.Errorf("field type %s is unsupported", t.Type.Name())
	}

	if defaultValue := getDefaultValue(t); defaultValue != nil {
		if err := field.Set(*defaultValue); err != nil {
			return nil, fmt.Errorf("failed to set default value: %s", err)
		}
	}

	return &field, nil
}

// Set accepts a raw string value and validates it based on struct tags, then
// passes it through to the value setter for further parsing
func (f Field) Set(s string) error {
	if f.Enum != nil && !slices.Contains(f.Enum, s) {
		return fmt.Errorf("value %s was not a member of [%s]", s, strings.Join(f.Enum, ", "))
	}

	if err := f.Value.Set(s); err != nil {
		return err
	}

	return nil
}

// BindFlag registers the field with the given flagset
func (f Field) BindFlag(fs FlagSet) {
	if f.FlagName != nil {
		fs.Var(f.Value, *f.FlagName, f.Usage)
	}
}
