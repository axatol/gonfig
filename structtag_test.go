package gonfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagDelimiter(t *testing.T) {
	t.Setenv("PROP", "1;2;3")
	target := struct {
		Prop []string `env:"PROP" delimiter:";"`
	}{}

	t.Cleanup(setCommandLineFlags(t))
	err := Load(&target)
	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "3"}, target.Prop)
}

func TestTagEnumUnmatched(t *testing.T) {
	t.Setenv("PROP", "foo")
	target := struct {
		Prop string `env:"PROP" enum:"debug,info,error"`
	}{}

	t.Cleanup(setCommandLineFlags(t))
	err := Load(&target)
	assert.ErrorContains(t, err, `failed to set field 'Prop': value foo was not a member of [debug, info, error]`)
	assert.Equal(t, "", target.Prop)
}

func TestTagEnumMatched(t *testing.T) {
	t.Setenv("PROP", "info")
	target := struct {
		Prop string `env:"PROP" enum:"debug,info,error"`
	}{}

	t.Cleanup(setCommandLineFlags(t))
	err := Load(&target)
	assert.NoError(t, err)
	assert.Equal(t, "info", target.Prop)
}

func TestTagRequiredInvalid(t *testing.T) {
	target := struct {
		Prop string `required:"prop"`
	}{}

	t.Cleanup(setCommandLineFlags(t))
	err := Load(&target)
	assert.ErrorContains(t, err, `failed to configure field from struct tags 'Prop': strconv.ParseBool: parsing "prop": invalid syntax`)
}

func TestTagRequiredValid(t *testing.T) {
	target := struct {
		Prop string `required:"true"`
	}{}

	t.Cleanup(setCommandLineFlags(t))
	err := Load(&target)
	assert.ErrorContains(t, err, "required fields are unset: [Prop]")
}
