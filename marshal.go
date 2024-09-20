package sjon

import (
	"errors"
	"fmt"
	"reflect"
)

const maxDepth = 1_000

func (s SJON) Marshal(v any) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	return s.reflectMarshal(reflect.ValueOf(v), 0)
}

func (s SJON) reflectMarshal(v reflect.Value, depth int) ([]byte, error) {
	if depth > maxDepth {
		return nil, errors.New("go-sjon: max depth exceeded")
	}
	m, ok := marshalers[v.Kind()]
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("go-sjon: unexpected kind %q", v.Kind()),
		)
	}
	var next marshalNext = func(v reflect.Value) ([]byte, error) {
		return s.reflectMarshal(v, depth+1)
	}
	return m(v, next)
}

type marshalNext func(reflect.Value) ([]byte, error)

var marshalers = map[reflect.Kind]func(reflect.Value, marshalNext) ([]byte, error){
	reflect.Array:      marshalNotSupported,
	reflect.Slice:      marshalNotSupported,
	reflect.Chan:       marshalNotSupported,
	reflect.Interface:  marshalNotSupported,
	reflect.Pointer:    marshalNotSupported,
	reflect.Struct:     marshalNotSupported,
	reflect.Map:        marshalNotSupported,
	reflect.Func:       marshalNotSupported,
	reflect.Int:        marshalNotSupported,
	reflect.Int8:       marshalNotSupported,
	reflect.Int16:      marshalNotSupported,
	reflect.Int32:      marshalNotSupported,
	reflect.Int64:      marshalNotSupported,
	reflect.Uint:       marshalNotSupported,
	reflect.Uint8:      marshalNotSupported,
	reflect.Uint16:     marshalNotSupported,
	reflect.Uint32:     marshalNotSupported,
	reflect.Uint64:     marshalNotSupported,
	reflect.Uintptr:    marshalNotSupported,
	reflect.String:     marshalNotSupported,
	reflect.Bool:       marshalBool,
	reflect.Float32:    marshalNotSupported,
	reflect.Float64:    marshalNotSupported,
	reflect.Complex64:  marshalNotSupported,
	reflect.Complex128: marshalNotSupported,
}

func marshalNotSupported(v reflect.Value, _ marshalNext) ([]byte, error) {
	return nil, errors.New(
		fmt.Sprintf("go-json: unsupported kind %q", v.Kind()),
	)
}

func marshalBool(v reflect.Value, _ marshalNext) ([]byte, error) {
	return []byte(fmt.Sprintf("%v", v.Bool())), nil
}
