package gonfig

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func Load(target any) error {
	config, err := NewConfig(target)
	if err != nil {
		return err
	}

	if err = config.ReadEnv(); err != nil {
		return err
	}

	if err = config.BindFlags(flag.CommandLine); err != nil {
		return err
	}

	flag.Parse()

	if err = config.Validate(); err != nil {
		return err
	}

	return nil
}

type Config struct {
	target any
	fields []*Field
}

func NewConfig(target any) (*Config, error) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return nil, fmt.Errorf("must specify non-nil pointer, got %s", v.Kind())
	}

	v = v.Elem()
	t := reflect.TypeOf(target).Elem()

	fields := []*Field{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		field, err := NewField(f, v.Field(i))
		if err != nil {
			return nil, fmt.Errorf("failed to configure field from struct tags '%s': %s", f.Name, err)
		}

		if !field.Required && field.EnvName == nil && field.FlagName == nil {
			continue
		}

		fields = append(fields, field)
	}

	config := Config{target, fields}
	return &config, nil
}

func (c Config) BindFlags(fs *flag.FlagSet) error {
	if fs.Parsed() {
		return fmt.Errorf("please do not call flagset.Parse() before config.BindFlags()")
	}

	for _, field := range c.fields {
		field.BindFlag(fs)
	}

	return nil
}

func (c Config) ReadEnv() error {
	for _, field := range c.fields {
		if field.EnvName == nil {
			continue
		}

		value, ok := os.LookupEnv(*field.EnvName)
		if !ok {
			continue
		}

		if err := field.Set(value); err != nil {
			return fmt.Errorf("failed to set field '%s': %s", field.Name, err)
		}
	}

	return nil
}

func (c Config) Validate() error {
	missing := []string{}
	for _, f := range c.fields {
		if f.Required && !f.Value.IsSet() {
			missing = append(missing, f.Name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required fields are unset: [%s]", strings.Join(missing, ", "))
	}

	return nil
}
