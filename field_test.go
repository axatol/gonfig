package gonfig

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewField(t *testing.T) {
	target := struct {
		Prop string `default:"foo" enum:"foo,bar" env:"PROP" flag:"prop" required:"true" usage:"usage"`
	}{}
	value := reflect.ValueOf(&target).Elem()
	typeof := value.Type().Field(0)
	valueof := value.Field(0)
	field, err := NewField(typeof, valueof)
	assert.NoError(t, err)
	assert.NotNil(t, field)
	t.Logf("v: %v\n", field)
	assert.Equal(t, "foo", field.Value.Get()) // default
	assert.Equal(t, []string{"foo", "bar"}, field.Enum)
	assert.NotNil(t, field.EnvName)
	assert.Equal(t, "PROP", *field.EnvName)
	assert.NotNil(t, field.FlagName)
	assert.Equal(t, "prop", *field.FlagName)
	assert.Equal(t, true, field.Required)
	assert.Equal(t, "usage", field.Usage)
}
