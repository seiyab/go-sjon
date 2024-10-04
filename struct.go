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
		tagItems := strings.Split(tag, ",")
		if len(tagItems) > 1 && tagItems[1] == "omitempty" && isEmpty(v.Field(i)) {
			continue
		}

		if !isFirst {
			out.Write([]byte(","))
		} else {
			isFirst = false
		}
		name := v.Type().Field(i).Name
		if tagItems[0] != "" {
			name = tagItems[0]
		} else if s.structKeyNamer != nil {
			name = s.structKeyNamer(name)
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

func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return v.Bool() == false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.IsZero()
	case reflect.Pointer, reflect.Interface:
		return v.IsNil()
	}
	return false
}
