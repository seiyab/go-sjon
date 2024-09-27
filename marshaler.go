package sjon

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

var jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()

func shouldMarshalWithMethod(v reflect.Value) bool {
	t := v.Type()
	return t.Implements(jsonMarshalerType) ||
		reflect.PointerTo(t).Implements(jsonMarshalerType)
}

func marshalWithMethod(v reflect.Value) ([]byte, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return []byte("null"), nil
		}
		v = v.Elem()
	}
	m, ok := v.Interface().(json.Marshaler)
	if !ok {
		return nil, errors.Errorf("go-json: internal error: %v does not implement json.Marshaler", v.Type())
	}
	return m.MarshalJSON()
}
