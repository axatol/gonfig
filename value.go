package gonfig

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Value represents a settable receiver which is compatible with flag.Value
type Value interface {
	IsSet() bool
	Set(string) error
	Get() any
	String() string
}

// ValueParser receives a raw string and returns the resolved value or an error
type ValueParser func(string) (any, error)

// ValueImpl is the built-in concrete implementation of Value
type ValueImpl struct {
	set    bool
	value  reflect.Value
	parser ValueParser
}

func (v ValueImpl) IsSet() bool {
	return v.set
}

func (v *ValueImpl) Set(raw string) error {
	if !v.value.CanSet() {
		return fmt.Errorf("field cannot be set")
	}

	value, err := v.parser(raw)
	if err != nil {
		return err
	}

	v.value.Set(reflect.ValueOf(value))
	v.set = true
	return nil
}

func (v ValueImpl) Get() any {
	return v.value.Interface()
}

func (v ValueImpl) String() string {
	return v.value.String()
}

func NewValue(t reflect.StructField, v reflect.Value) Value {
	delim := getDelimiter(t)

	switch v.Type() {

	// bool

	case reflect.TypeFor[bool]():
		return &ValueImpl{value: v, parser: parseBoolValue}

	case reflect.TypeFor[[]bool]():
		return &ValueImpl{value: v, parser: parseSliceValue[bool](delim, parseBoolValue)}

	// strings

	case reflect.TypeFor[time.Duration]():
		return &ValueImpl{value: v, parser: ValueParser(parseDurationValue)}

	case reflect.TypeFor[[]time.Duration]():
		return &ValueImpl{value: v, parser: parseSliceValue[time.Duration](delim, parseDurationValue)}

	case reflect.TypeFor[string]():
		return &ValueImpl{value: v, parser: parseStringValue}

	case reflect.TypeFor[[]string]():
		return &ValueImpl{value: v, parser: parseSliceValue[string](delim, parseStringValue)}

	// int

	case reflect.TypeFor[int]():
		return &ValueImpl{value: v, parser: parseIntValue[int]}

	case reflect.TypeFor[[]int]():
		return &ValueImpl{value: v, parser: parseSliceValue[int](delim, parseIntValue[int])}

	case reflect.TypeFor[int8]():
		return &ValueImpl{value: v, parser: parseIntValue[int8]}

	case reflect.TypeFor[[]int8]():
		return &ValueImpl{value: v, parser: parseSliceValue[int8](delim, parseIntValue[int8])}

	case reflect.TypeFor[int16]():
		return &ValueImpl{value: v, parser: parseIntValue[int16]}

	case reflect.TypeFor[[]int16]():
		return &ValueImpl{value: v, parser: parseSliceValue[int16](delim, parseIntValue[int16])}

	case reflect.TypeFor[int32]():
		return &ValueImpl{value: v, parser: parseIntValue[int32]}

	case reflect.TypeFor[[]int32]():
		return &ValueImpl{value: v, parser: parseSliceValue[int32](delim, parseIntValue[int32])}

	case reflect.TypeFor[int64]():
		return &ValueImpl{value: v, parser: parseIntValue[int64]}

	case reflect.TypeFor[[]int64]():
		return &ValueImpl{value: v, parser: parseSliceValue[int64](delim, parseIntValue[int64])}

	// uint

	case reflect.TypeFor[uint]():
		return &ValueImpl{value: v, parser: parseUintValue[uint]}

	case reflect.TypeFor[[]uint]():
		return &ValueImpl{value: v, parser: parseSliceValue[uint](delim, parseUintValue[uint])}

	case reflect.TypeFor[uint8]():
		return &ValueImpl{value: v, parser: parseUintValue[uint8]}

	case reflect.TypeFor[[]uint8]():
		return &ValueImpl{value: v, parser: parseSliceValue[uint8](delim, parseUintValue[uint8])}

	case reflect.TypeFor[uint16]():
		return &ValueImpl{value: v, parser: parseUintValue[uint16]}

	case reflect.TypeFor[[]uint16]():
		return &ValueImpl{value: v, parser: parseSliceValue[uint16](delim, parseUintValue[uint16])}

	case reflect.TypeFor[uint32]():
		return &ValueImpl{value: v, parser: parseUintValue[uint32]}

	case reflect.TypeFor[[]uint32]():
		return &ValueImpl{value: v, parser: parseSliceValue[uint32](delim, parseUintValue[uint32])}

	case reflect.TypeFor[uint64]():
		return &ValueImpl{value: v, parser: parseUintValue[uint64]}

	case reflect.TypeFor[[]uint64]():
		return &ValueImpl{value: v, parser: parseSliceValue[uint64](delim, parseUintValue[uint64])}

	// float

	case reflect.TypeFor[float32]():
		return &ValueImpl{value: v, parser: parseFloatValue[float32]}

	case reflect.TypeFor[[]float32]():
		return &ValueImpl{value: v, parser: parseSliceValue[float32](delim, parseFloatValue[float32])}

	case reflect.TypeFor[float64]():
		return &ValueImpl{value: v, parser: parseFloatValue[float64]}

	case reflect.TypeFor[[]float64]():
		return &ValueImpl{value: v, parser: parseSliceValue[float64](delim, parseFloatValue[float64])}

	}

	return nil
}

func parseStringValue(s string) (any, error) {
	return s, nil
}

func parseDurationValue(s string) (any, error) {
	return time.ParseDuration(s)
}

func parseIntValue[T int64 | int32 | int16 | int8 | int](s string) (any, error) {
	value, err := strconv.ParseInt(s, 10, 64)
	return T(value), err
}

func parseUintValue[T uint64 | uint32 | uint16 | uint8 | uint](s string) (any, error) {
	value, err := strconv.ParseUint(s, 10, 64)
	return T(value), err
}

func parseFloatValue[T float64 | float32](s string) (any, error) {
	value, err := strconv.ParseFloat(s, 64)
	return T(value), err
}

func parseBoolValue(s string) (any, error) {
	return strconv.ParseBool(s)
}

func parseSliceValue[T any](delim string, parser ValueParser) ValueParser {
	return func(s string) (any, error) {
		inputs := strings.Split(s, delim)
		outputs := make([]T, len(inputs))

		for i, elem := range inputs {
			parsed, err := parser(elem)
			if err != nil {
				return nil, fmt.Errorf("could not parse slice element %d: %s", i, err)
			}

			value, ok := (parsed).(T)
			if !ok {
				return nil, fmt.Errorf("element %d parsed to unexpected type %T, expected %T", i, parsed, reflect.TypeFor[T]())
			}

			outputs[i] = value
		}

		return outputs, nil
	}
}
