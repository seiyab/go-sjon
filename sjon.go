package sjon

type Serializer struct {
	structKeyNamer func(string) string
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
