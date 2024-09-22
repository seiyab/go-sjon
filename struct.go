package sjon

import (
	"io"
	"reflect"
	"strconv"
	"strings"
)

func marshalStruct(s *Serializer, v reflect.Value, out io.Writer, next marshalNext) error {
	out.Write([]byte("{"))
	isFirst := true
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		tag := f.Tag.Get("json")
		if tag == "-" {
			continue
		}
		if !isFirst {
			out.Write([]byte(","))
		} else {
			isFirst = false
		}
		name := v.Type().Field(i).Name
		if tag != "" {
			tagItems := strings.Split(tag, ",")
			if tagItems[0] != "" {
				name = tagItems[0]
			}
		} else if s.KeyNamer != nil {
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
