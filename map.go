package sjon

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"sort"
)

type printedMapKey struct {
	key     reflect.Value
	printed []byte
}

func marshalMap(_ *Serializer, v reflect.Value, out io.Writer, next marshalNext) error {
	_, err := out.Write([]byte("{"))
	if err != nil {
		return err
	}
	pmks, err := printMapKeys(v, next)
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

func printMapKeys(v reflect.Value, next marshalNext) ([]printedMapKey, error) {
	keys := v.MapKeys()
	printedToPK := make(map[string]printedMapKey)
	for _, key := range keys {
		buf := bytes.NewBuffer(nil)
		err := next(key, buf)
		if err != nil {
			return nil, err
		}
		bs := buf.Bytes()
		if _, ok := printedToPK[string(bs)]; ok {
			// TODO: consider priority and handle correctly
			return nil, errors.New("go-sjon: duplicate key")
		}
		printedToPK[string(bs)] = printedMapKey{key, bs}
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
