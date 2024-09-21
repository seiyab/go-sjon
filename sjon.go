package sjon

type Serializer struct {
	KeyNamer func(string) string
}

func NewSerializer() Serializer {
	return Serializer{}
}
