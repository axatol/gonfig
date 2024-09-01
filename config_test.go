package gonfig

import (
	"flag"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPointerInput(t *testing.T) {
	target := struct{}{}
	_, err := NewConfig(target)
	assert.ErrorContains(t, err, "must specify non-nil pointer")
}

func TestAllTypes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{name: "bool", input: "true", expected: true},
		{name: "bool_slice", input: "false,true", expected: []bool{false, true}},

		{name: "duration", input: "1h5s", expected: time.Hour + time.Second*5},
		{name: "duration_slice", input: "2h,4s", expected: []time.Duration{time.Hour * 2, time.Second * 4}},
		{name: "string", input: "prop", expected: "prop"},
		{name: "string_slice", input: "lorem,ipsum", expected: []string{"lorem", "ipsum"}},

		{name: "float32", input: "8", expected: float32(8)},
		{name: "float32_slice", input: "7,6", expected: []float32{7, 6}},
		{name: "float64", input: "8", expected: float64(8)},
		{name: "float64_slice", input: "7,6", expected: []float64{7, 6}},

		{name: "int", input: "1", expected: int(1)},
		{name: "int_slice", input: "1,2,3", expected: []int{1, 2, 3}},
		{name: "int8", input: "6", expected: int8(6)},
		{name: "int8_slice", input: "5,4", expected: []int8{5, 4}},
		{name: "int16", input: "6", expected: int16(6)},
		{name: "int16_slice", input: "5,4", expected: []int16{5, 4}},
		{name: "int32", input: "6", expected: int32(6)},
		{name: "int32_slice", input: "5,4", expected: []int32{5, 4}},
		{name: "int64", input: "6", expected: int64(6)},
		{name: "int64_slice", input: "5,4", expected: []int64{5, 4}},

		{name: "uint", input: "1", expected: uint(1)},
		{name: "uint_slice", input: "1,2,3", expected: []uint{1, 2, 3}},
		{name: "uint8", input: "5", expected: uint8(5)},
		{name: "uint8_slice", input: "4,3", expected: []uint8{4, 3}},
		{name: "uint16", input: "5", expected: uint16(5)},
		{name: "uint16_slice", input: "4,3", expected: []uint16{4, 3}},
		{name: "uint32", input: "5", expected: uint32(5)},
		{name: "uint32_slice", input: "4,3", expected: []uint32{4, 3}},
		{name: "uint64", input: "5", expected: uint64(5)},
		{name: "uint64_slice", input: "4,3", expected: []uint64{4, 3}},
	}

	target := struct {
		Bool      bool   `env:"bool"`
		BoolSlice []bool `env:"bool_slice"`

		Duration      time.Duration   `env:"duration"`
		DurationSlice []time.Duration `env:"duration_slice"`
		String        string          `env:"string"`
		StringSlice   []string        `env:"string_slice"`

		Float32      float32   `env:"float32"`
		Float32Slice []float32 `env:"float32_slice"`
		Float64      float64   `env:"float64"`
		Float64Slice []float64 `env:"float64_slice"`

		Int        int     `env:"int"`
		IntSlice   []int   `env:"int_slice"`
		Int8       int8    `env:"int8"`
		Int8Slice  []int8  `env:"int8_slice"`
		Int16      int16   `env:"int16"`
		Int16Slice []int16 `env:"int16_slice"`
		Int32      int32   `env:"int32"`
		Int32Slice []int32 `env:"int32_slice"`
		Int64      int64   `env:"int64"`
		Int64Slice []int64 `env:"int64_slice"`

		Uint        uint     `env:"uint"`
		UintSlice   []uint   `env:"uint_slice"`
		Uint8       uint8    `env:"uint8"`
		Uint8Slice  []uint8  `env:"uint8_slice"`
		Uint16      uint16   `env:"uint16"`
		Uint16Slice []uint16 `env:"uint16_slice"`
		Uint32      uint32   `env:"uint32"`
		Uint32Slice []uint32 `env:"uint32_slice"`
		Uint64      uint64   `env:"uint64"`
		Uint64Slice []uint64 `env:"uint64_slice"`
	}{}

	for _, test := range tests {
		t.Setenv(test.name, test.input)
	}

	// err = config.ReadEnv()
	// assert.NoError(t, err)

	var oldFlags = flag.CommandLine
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	t.Cleanup(func() { flag.CommandLine = oldFlags })

	err := Load(&target)
	assert.NoError(t, err)

	v := reflect.ValueOf(target)
	for i, test := range tests {
		actual := v.Field(i).Interface()
		assert.Equal(t, test.expected, actual)
	}
}
