package sjon

import "reflect"

type Serializer struct {
	structKeyNamer func(string) string
	replacers      map[reflect.Type]func(reflect.Value) reflect.Value
}

func NewSerializer() Serializer {
	return Serializer{}
}

type Config[T any] func(v T) T

func (s Serializer) With(config Config[Serializer]) Serializer {
	return config(s)
}

func StructKeyNamer(keyNamer func(string) string) Config[Serializer] {
	return func(s Serializer) Serializer {
		s.structKeyNamer = keyNamer
		return s
	}
}

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
