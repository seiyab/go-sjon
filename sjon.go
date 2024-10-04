package sjon

type Serializer struct {
	structKeyNamer func(string) string
}

func NewSerializer() Serializer {
	return Serializer{}
}

func (s Serializer) WithStructKeyNamer(kn func(string) string) Serializer {
	s.structKeyNamer = kn
	return s
}
