package sjon

type SJON struct {
	KeyNamer func(string) string
}

func New() SJON {
	return SJON{}
}
