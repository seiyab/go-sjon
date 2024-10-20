package sjon

import "reflect"

// Serializer provides JSON serialization with non-invasive customization
type Serializer struct {
	structKeyNamer func(string) string
	replacers      map[reflect.Type]func(reflect.Value) reflect.Value
}

// NewSerializer returns a new Serializer
func NewSerializer() Serializer {
	return Serializer{}
}

// Config is used to configure instances
// Usual user don't need to know about it
type Config[T any] func(v T) T

// With returns a new Serializer with passed configuration
func (s Serializer) With(config Config[Serializer]) Serializer {
	return config(s)
}

// StructKeyNamer returns a Config that changes how Serializer serializes struct key
func StructKeyNamer(keyNamer func(string) string) Config[Serializer] {
	return func(s Serializer) Serializer {
		s.structKeyNamer = keyNamer
		return s
	}
}

// Replacer takes a replace function and returns a config
// Serializer with Replacer applies the replace function into values with replace function's argument type
func Replacer[Before, After any](replace func(Before) After) Config[Serializer] {
	ty := reflect.TypeOf(replace)
	replaceValue := reflect.ValueOf(replace)
	reflectReplace := func(v reflect.Value) reflect.Value {
		rtn := replaceValue.Call([]reflect.Value{v})
		return rtn[0]
	}
	return func(s Serializer) Serializer {
		s.replacers = cloneMap(s.replacers) // TODO: need to improve time complexity?
		s.replacers[ty.In(0)] = reflectReplace
		return s
	}
}
