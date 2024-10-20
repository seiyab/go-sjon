package sjon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

const maxDepth = 1_000

func (s Serializer) Marshal(v any) ([]byte, error) {
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

func (s Serializer) reflectMarshal(v reflect.Value, out io.Writer, depth int) error {
	if depth > maxDepth {
		return errors.New("go-sjon: max depth exceeded")
	}
	if r, ok := s.replacers[v.Type()]; ok {
		return s.reflectMarshal(r(v), out, depth+1)
	}
	ok, err := marshalWithMarshalJSON(v, out)
	if err != nil || ok {
		return err
	}
	ok, err = marshalWithMarshalText(v, out)
	if err != nil || ok {
		return err
	}

	m, ok := marshalers[v.Kind()]
	if !ok {
		return fmt.Errorf("go-sjon: unexpected kind %q", v.Kind())
	}
	var next marshalNext = func(v reflect.Value, w io.Writer) error {
		return s.reflectMarshal(v, w, depth+1)
	}
	return m(&s, v, out, next)
}

type marshalNext func(reflect.Value, io.Writer) error

var marshalers = map[reflect.Kind]func(*Serializer, reflect.Value, io.Writer, marshalNext) error{
	reflect.Array:      marshalArray,
	reflect.Slice:      marshalArray,
	reflect.Chan:       marshalNotSupported,
	reflect.Interface:  marshalPointer,
	reflect.Pointer:    marshalPointer,
	reflect.Struct:     marshalStruct,
	reflect.Map:        marshalMap,
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
	reflect.String:     marshalString,
	reflect.Bool:       marshalBool,
	reflect.Float32:    marshalFloat[float32],
	reflect.Float64:    marshalFloat[float64],
	reflect.Complex64:  marshalNotSupported,
	reflect.Complex128: marshalNotSupported,
}

func marshalNotSupported(_ *Serializer, v reflect.Value, _ io.Writer, _ marshalNext) error {
	return fmt.Errorf("go-json: unsupported kind %q", v.Kind())
}

func marshalArray(_ *Serializer, v reflect.Value, out io.Writer, next marshalNext) error {
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
		err := next(v.Index(i), out)
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

func marshalPointer(_ *Serializer, v reflect.Value, out io.Writer, next marshalNext) error {
	if !v.IsNil() {
		return next(v.Elem(), out)
	}
	_, err := out.Write([]byte("null"))
	if err != nil {
		return err
	}
	return nil
}

func marshalInt(_ *Serializer, v reflect.Value, out io.Writer, _ marshalNext) error {
	_, err := out.Write([]byte(strconv.FormatInt(v.Int(), 10)))
	if err != nil {
		return err
	}
	return nil
}

func marshalUint(_ *Serializer, v reflect.Value, out io.Writer, _ marshalNext) error {
	_, err := out.Write([]byte(strconv.FormatUint(v.Uint(), 10)))
	if err != nil {
		return err
	}
	return nil
}

func marshalString(_ *Serializer, v reflect.Value, out io.Writer, _ marshalNext) error {
	b, err := json.Marshal(v.String())
	if err != nil {
		return err
	}
	_, err = out.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func marshalBool(_ *Serializer, v reflect.Value, out io.Writer, _ marshalNext) error {
	_, err := out.Write([]byte(fmt.Sprintf("%v", v.Bool())))
	if err != nil {
		return err
	}
	return nil
}

func marshalFloat[T ~float32 | float64](_ *Serializer, v reflect.Value, out io.Writer, _ marshalNext) error {
	b, err := json.Marshal(T(v.Float()))
	if err != nil {
		return err
	}
	_, err = out.Write(b)
	if err != nil {
		return err
	}
	return nil
}
