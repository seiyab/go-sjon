package sjon

import "reflect"

func cloneMap[V any](m map[reflect.Type]V) map[reflect.Type]V {
	o := make(map[reflect.Type]V, len(m))
	for k, v := range m {
		o[k] = v
	}
	return o
}
