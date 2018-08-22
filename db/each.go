package db

import (
	"reflect"
	"strings"
)

func each(v reflect.Value, keys map[string]bool, fn func(field Field) bool) bool {

	fd := Field{}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {

		fd.F = t.Field(i)
		fd.V = v.Field(i)

		if fd.F.Type.Kind() == reflect.Struct {
			if each(fd.V, keys, fn) {
				continue
			}
			return false
		}

		fd.Name = fd.F.Tag.Get("name")

		if fd.Name == "" {
			fd.Name = fd.F.Tag.Get("json")
			if fd.Name == "-" {
				continue
			}
			if fd.Name != "" {
				fd.Name = strings.Split(fd.Name, ",")[0]
			}
		} else if fd.Name == "-" {
			continue
		}

		if fd.Name == "" {
			fd.Name = strings.ToLower(fd.F.Name)
		}

		if keys[fd.Name] {
			continue
		}

		keys[fd.Name] = true

		switch fd.F.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint:
			fd.IsObject = false
		case reflect.Int64, reflect.Uint64:
			fd.IsObject = false
		case reflect.Float32, reflect.Float64:
			fd.IsObject = false
		case reflect.Bool:
			fd.IsObject = false
		case reflect.String:
			fd.IsObject = false
		default:
			fd.IsObject = true
		}

		if !fn(fd) {
			return false
		}
	}

	return true
}

func Each(object IObject, fn func(field Field) bool) bool {

	if object == nil {
		return true
	}

	keys := map[string]bool{}

	v := reflect.ValueOf(object)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	} else {
		return true
	}

	return each(v, keys, fn)
}
