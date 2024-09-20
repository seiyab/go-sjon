package sjon

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

const maxDepth = 1_000

func (s SJON) Marshal(v any) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	buf := bytes.NewBuffer(nil)
	err := s.reflectMarshal(reflect.ValueOf(v), buf, 0)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s SJON) reflectMarshal(v reflect.Value, out io.Writer, depth int) error {
	if depth > maxDepth {
		return errors.New("go-sjon: max depth exceeded")
	}
	m, ok := marshalers[v.Kind()]
	if !ok {
		return errors.New(
			fmt.Sprintf("go-sjon: unexpected kind %q", v.Kind()),
		)
	}
	var next marshalNext = func(v reflect.Value) error {
		return s.reflectMarshal(v, out, depth+1)
	}
	return m(v, out, next)
}

type marshalNext func(reflect.Value) error

var marshalers = map[reflect.Kind]func(reflect.Value, io.Writer, marshalNext) error{
	reflect.Array:      marshalArray,
	reflect.Slice:      marshalArray,
	reflect.Chan:       marshalNotSupported,
	reflect.Interface:  marshalNotSupported,
	reflect.Pointer:    marshalNotSupported,
	reflect.Struct:     marshalNotSupported,
	reflect.Map:        marshalNotSupported,
	reflect.Func:       marshalNotSupported,
	reflect.Int:        marshalInt,
	reflect.Int8:       marshalInt,
	reflect.Int16:      marshalInt,
	reflect.Int32:      marshalInt,
	reflect.Int64:      marshalInt,
	reflect.Uint:       marshalUint,
	reflect.Uint8:      marshalUint,
	reflect.Uint16:     marshalUint,
	reflect.Uint32:     marshalUint,
	reflect.Uint64:     marshalUint,
	reflect.Uintptr:    marshalNotSupported,
	reflect.String:     marshalNotSupported,
	reflect.Bool:       marshalBool,
	reflect.Float32:    marshalNotSupported,
	reflect.Float64:    marshalNotSupported,
	reflect.Complex64:  marshalNotSupported,
	reflect.Complex128: marshalNotSupported,
}

func marshalNotSupported(v reflect.Value, _ io.Writer, _ marshalNext) error {
	return errors.New(
		fmt.Sprintf("go-json: unsupported kind %q", v.Kind()),
	)
}

func marshalArray(v reflect.Value, out io.Writer, next marshalNext) error {
	_, err := out.Write([]byte("["))
	if err != nil {
		return err
	}
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			_, err := out.Write([]byte(","))
			if err != nil {
				return err
			}
		}
		err := next(v.Index(i))
		if err != nil {
			return err
		}
	}
	_, err = out.Write([]byte("]"))
	if err != nil {
		return err
	}
	return nil
}

func marshalInt(v reflect.Value, out io.Writer, _ marshalNext) error {
	_, err := out.Write([]byte(strconv.FormatInt(v.Int(), 10)))
	if err != nil {
		return err
	}
	return nil
}

func marshalUint(v reflect.Value, out io.Writer, _ marshalNext) error {
	_, err := out.Write([]byte(strconv.FormatUint(v.Uint(), 10)))
	if err != nil {
		return err
	}
	return nil
}

func marshalBool(v reflect.Value, out io.Writer, _ marshalNext) error {
	_, err := out.Write([]byte(fmt.Sprintf("%v", v.Bool())))
	if err != nil {
		return err
	}
	return nil
}
