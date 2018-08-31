package json

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
)

func encodeObject(v reflect.Value, w *bytes.Buffer, idx int) (int, error) {

	var err error = nil
	ftype := v.Type()

	for i := 0; i < ftype.NumField(); i++ {
		fd := ftype.Field(i)
		fv := v.Field(i)

		var tags = strings.Split(fd.Tag.Get("json"), ",")

		if tags != nil && len(tags) > 0 && tags[0] != "" {

			if tags[0] == "-" {
				continue
			}

			if fv.CanInterface() {
				if len(tags) > 1 && tags[1] == "omitempty" && dynamic.IsNil(fv.Interface()) {
					continue
				}

				if idx != 0 {
					w.WriteString(",")
				}

				err = encode(tags[0], w)

				if err != nil {
					return idx, err
				}

				w.WriteString(":")

				err = encode(fv.Interface(), w)

				if err != nil {
					return idx, err
				}
				idx = idx + 1
			}

		} else if fv.Kind() == reflect.Struct {
			idx, err = encodeObject(fv, w, idx)
			if err != nil {
				return idx, err
			}
		}

	}

	return idx, nil
}

func encode(object interface{}, w *bytes.Buffer) error {

	if object == nil {
		w.WriteString("null")
		return nil
	}

	v := reflect.ValueOf(object)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		i := 0
		w.WriteString("{")
		for _, key := range v.MapKeys() {
			vv := v.MapIndex(key)
			if key.CanInterface() && vv.CanInterface() {
				if i != 0 {
					w.WriteString(",")
				}
				err := encode(dynamic.StringValue(key.Interface(), ""), w)
				if err != nil {
					return err
				}
				w.WriteString(":")
				err = encode(vv.Interface(), w)
				if err != nil {
					return err
				}
				i = i + 1
			}
		}
		w.WriteString("}")
	case reflect.Slice:
		w.WriteString("[")
		for i := 0; i < v.Len(); i++ {
			vv := v.Index(i)
			if vv.CanInterface() {
				if i != 0 {
					w.WriteString(",")
				}
				err := encode(vv.Interface(), w)
				if err != nil {
					return err
				}
			}
		}
		w.WriteString("]")
	case reflect.Struct:
		w.WriteString("{")
		_, err := encodeObject(v, w, 0)
		if err != nil {
			return err
		}
		w.WriteString("}")
	default:
		b, err := json.Marshal(object)
		if err != nil {
			return err
		}
		w.Write(b)
	}

	return nil
}

func Marshal(object interface{}) ([]byte, error) {

	w := bytes.NewBuffer(nil)

	err := encode(object, w)

	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func MarshalIndent(object interface{}, prefix, indent string) ([]byte, error) {

	w := bytes.NewBuffer(nil)

	enc := json.NewEncoder(w)

	enc.SetIndent(prefix, indent)

	err := enc.Encode(object)

	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
