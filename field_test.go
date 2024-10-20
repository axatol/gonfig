package gonfig

import (
	"flag"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestField(t *testing.T, target any) *Field {
	t.Helper()

	value := reflect.ValueOf(target).Elem()
	typeof := value.Type().Field(0)
	valueof := value.Field(0)

	field, err := NewField(typeof, valueof)
	assert.NoError(t, err)
	assert.NotNil(t, field)

	return field
}

func TestFieldDefault(t *testing.T) {
	field := newTestField(t, &struct {
		Prop string `default:"foo"`
	}{})

	// default value
	assert.Equal(t, "foo", field.Value.Get())

	// overridden
	err := field.Set("bar")
	assert.NoError(t, err)
	assert.Equal(t, "bar", field.Value.Get())
}

func TestFieldDelimiterDefault(t *testing.T) {
	field := newTestField(t, &struct {
		Prop []string
	}{})

	// default delimiter
	assert.Equal(t, ",", field.Delimiter)

	// matched default
	err := field.Set("foo,bar")
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, field.Value.Get())

	// unmatched
	err = field.Set("foo;bar")
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo;bar"}, field.Value.Get())

}

func TestFieldDelimiterCustom(t *testing.T) {
	field := newTestField(t, &struct {
		Prop []string `delim:";"`
	}{})

	// customised delimiter
	assert.Equal(t, ";", field.Delimiter)

	// unmatched default
	err := field.Set("foo,bar")
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo,bar"}, field.Value.Get())

	// matched
	err = field.Set("foo;bar")
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, field.Value.Get())
}

func TestFieldEnum(t *testing.T) {
	field := newTestField(t, &struct {
		Prop string `enum:"foo,bar"`
	}{})

	// customised enum
	assert.Equal(t, []string{"foo", "bar"}, field.Enum)

	// unmatched
	err := field.Set("baz")
	assert.ErrorContains(t, err, "value 'baz' was not a member of [foo, bar]")

	// matched
	err = field.Set("foo")
	assert.NoError(t, err)
	assert.Equal(t, "foo", field.Value.Get())
}

func TestFieldEnumCustomDelim(t *testing.T) {
	field := newTestField(t, &struct {
		Prop string `enum:"foo;bar" delim:";"`
	}{})

	// customised enum
	assert.Equal(t, []string{"foo", "bar"}, field.Enum)
}

func TestEnumSlice(t *testing.T) {
	field := newTestField(t, &struct {
		Prop []string `enum:"foo,bar"`
	}{})

	// customised enum
	assert.Equal(t, []string{"foo", "bar"}, field.Enum)

	// unmatched
	err := field.Set("foo,baz")
	assert.ErrorContains(t, err, "value 'baz' was not a member of [foo, bar]")

	// matched
	err = field.Set("foo,bar")
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, field.Value.Get())
}

func TestFieldEnvName(t *testing.T) {
	target := struct {
		Prop string `env:"PROP"`
	}{}
	config, err := NewConfig(&target)
	assert.NoError(t, err)

	// customised environment variable name
	assert.Equal(t, "PROP", *config.fields[0].EnvName)

	// field set
	t.Setenv("PROP", "foo")
	err = config.ReadEnv()
	assert.NoError(t, err)
	assert.Equal(t, "foo", target.Prop)
}

func TestFieldEnvNameIgnored(t *testing.T) {
	target := struct {
		Prop string
	}{}
	config, err := NewConfig(&target)
	assert.NoError(t, err)

	// no environment variable name to read
	assert.Nil(t, config.fields[0].EnvName)

	// field ignored
	t.Setenv("PROP", "foo")
	err = config.ReadEnv()
	assert.NoError(t, err)
	assert.Equal(t, "", target.Prop)
}

func TestFieldBindFlag(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	field := newTestField(t, &struct {
		Prop string `flag:"prop"`
	}{})

	// customised flag
	assert.Equal(t, "prop", *field.FlagName)

	// flag set
	field.BindFlag(fs)
	err := fs.Parse([]string{"-prop", "foo"})
	assert.NoError(t, err)
	assert.Equal(t, "foo", field.Value.Get())
}

func TestFieldBindFlagIgnored(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	field := newTestField(t, &struct {
		Prop string
	}{})

	// flag not set
	assert.Nil(t, field.FlagName)

	// flag not registered
	field.BindFlag(fs)
	fs.VisitAll(func(f *flag.Flag) {
		assert.NotEqual(t, "prop", f.Name)
	})
}

func TestFieldRequred(t *testing.T) {
	field := newTestField(t, &struct {
		Prop string `required:"true"`
	}{})

	// customised required
	assert.True(t, field.Required)
}

func TestFieldRequiredIgnored(t *testing.T) {
	field := newTestField(t, &struct {
		Prop string
	}{})

	// required not set
	assert.False(t, field.Required)
}
