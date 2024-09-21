package sjon

import (
	"io"
	"reflect"
	"strconv"
)

func marshalStruct(s *Serializer, v reflect.Value, out io.Writer, next marshalNext) error {
	out.Write([]byte("{"))
	for i := 0; i < v.NumField(); i++ {
		if i > 0 {
			out.Write([]byte(","))
		}
		name := v.Type().Field(i).Name
		if s.KeyNamer != nil {
			name = s.KeyNamer(name)
		}
		out.Write([]byte(strconv.Quote(name)))
		out.Write([]byte(":"))
		err := next(v.Field(i))
		if err != nil {
			return err
		}
	}
	out.Write([]byte("}"))
	return nil
}
