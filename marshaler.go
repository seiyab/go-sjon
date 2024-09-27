package sjon

import (
	"encoding"
	"encoding/json"
	"io"
	"reflect"

	"github.com/pkg/errors"
)

var (
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

func marshalWithMarshalJSON(v reflect.Value, out io.Writer) (bool, error) {
	return marshalWithMethod(v, out, func(v json.Marshaler) ([]byte, error) {
		return v.MarshalJSON()
	})
}

func marshalWithMarshalText(v reflect.Value, out io.Writer) (bool, error) {
	return marshalWithMethod(v, out, func(v encoding.TextMarshaler) ([]byte, error) {
		b, err := v.MarshalText()
		if err != nil {
			return nil, err
		}
		q := make([]byte, 0, len(b)+2)
		q = append(q, '"')
		q = append(q, b...)
		q = append(q, '"')
		return q, nil
	})
}

func marshalWithMethod[T any](
	v reflect.Value,
	out io.Writer,
	method func(v T) ([]byte, error),
) (bool, error) {
	marshalerType := reflect.TypeOf((*T)(nil)).Elem()
	t := v.Type()
	if !(t.Implements(marshalerType) ||
		reflect.PointerTo(t).Implements(marshalerType)) {
		return false, nil
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			_, err := out.Write([]byte("null"))
			if err != nil {
				return false, err
			}
			return true, nil
		}
		v = v.Elem()
	}
	m, ok := v.Interface().(T)
	if !ok {
		return false, errors.Errorf("go-json: internal error: %v does not implement Marshaler", v.Type())
	}
	b, err := method(m)
	if err != nil {
		return false, err
	}
	_, err = out.Write(b)
	if err != nil {
		return false, err
	}
	return true, nil
}
