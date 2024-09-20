package sjon

import (
	"io"
	"reflect"
	"strconv"
)

func marshalStruct(v reflect.Value, out io.Writer, next marshalNext) error {
	out.Write([]byte("{"))
	for i := 0; i < v.NumField(); i++ {
		if i > 0 {
			out.Write([]byte(","))
		}
		out.Write([]byte(strconv.Quote(v.Type().Field(i).Name)))
		out.Write([]byte(":"))
		err := next(v.Field(i))
		if err != nil {
			return err
		}
	}
	out.Write([]byte("}"))
	return nil
}