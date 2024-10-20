package sjon

import (
	"bytes"
	"io"
	"reflect"
	"sort"

	"github.com/pkg/errors"
)

type printedMapKey struct {
	key     reflect.Value
	printed []byte
}

func marshalMap(s *Serializer, v reflect.Value, out io.Writer, next marshalNext) error {
	_, err := out.Write([]byte("{"))
	if err != nil {
		return err
	}
	pmks, err := printMapKeys(s, v, next)
	if err != nil {
		return err
	}
	for i, key := range pmks {
		if i > 0 {
			_, err := out.Write([]byte(","))
			if err != nil {
				return err
			}
		}
		_, err := out.Write(key.printed)
		if err != nil {
			return err
		}
		_, err = out.Write([]byte(":"))
		if err != nil {
			return err
		}
		err = next(v.MapIndex(key.key), out)
		if err != nil {
			return err
		}
	}
	_, err = out.Write([]byte("}"))
	if err != nil {
		return err
	}
	return nil
}

func printMapKeys(s *Serializer, v reflect.Value, next marshalNext) ([]printedMapKey, error) {
	keys := v.MapKeys()
	printedToPK := make(map[string]printedMapKey)
	for _, originalKey := range keys {
		key := originalKey
		replacer, ok := s.replacers[key.Type()]
		if ok {
			key = replacer(key)
		}
		buf := bytes.NewBuffer(nil)
		switch key.Kind() {
		case reflect.String:
			err := next(key, buf)
			if err != nil {
				return nil, err
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			_, err := buf.Write([]byte{'"'})
			if err != nil {
				return nil, err
			}
			err = next(key, buf)
			if err != nil {
				return nil, err
			}
			_, err = buf.Write([]byte{'"'})
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.Errorf("go-sjon: unsupported map key type %q", key.Kind())
		}
		bs := buf.Bytes()
		if _, ok := printedToPK[string(bs)]; ok {
			// TODO: consider priority and handle correctly
			return nil, errors.New("go-sjon: duplicate key")
		}
		printedToPK[string(bs)] = printedMapKey{originalKey, bs}
	}
	pks := make([]printedMapKey, 0, len(keys))
	for _, v := range printedToPK {
		pks = append(pks, v)
	}
	sort.Slice(pks, func(i, j int) bool {
		return string(pks[i].printed) < string(pks[j].printed)
	})
	return pks, nil
}
