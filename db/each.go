package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hailongz/kk-lib/dynamic"
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

		fd.IsObject = false
		fd.DBValue = ""
		fd.DBIndex = fd.F.Tag.Get("index")
		length := dynamic.IntValue(fd.F.Tag.Get("length"), 0)

		keys[fd.Name] = true

		switch fd.F.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint:
			if length != 0 {
				fd.DBType = fmt.Sprintf("INT(%d)", length)
			} else {
				fd.DBType = "INT"
			}
			fd.DBValue = dynamic.StringValue(dynamic.IntValue(fd.F.Tag.Get("default"), 0), "0")
		case reflect.Int64, reflect.Uint64:
			if length != 0 {
				fd.DBType = fmt.Sprintf("BIGINT(%d)", length)
			} else {
				fd.DBType = "BIGINT"
			}
			fd.DBValue = dynamic.StringValue(dynamic.IntValue(fd.F.Tag.Get("default"), 0), "0")
		case reflect.Float32, reflect.Float64:
			if length != 0 {
				fd.DBType = fmt.Sprintf("DOUBLE(%d)", length)
			} else {
				fd.DBType = "DOUBLE"
			}
			fd.DBValue = dynamic.StringValue(dynamic.FloatValue(fd.F.Tag.Get("default"), 0), "0")
		case reflect.Bool:
			fd.DBType = "INT(1)"
			if dynamic.BooleanValue(fd.F.Tag.Get("default"), false) {
				fd.DBValue = "1"
			} else {
				fd.DBValue = "0"
			}
		case reflect.String:
			if length == 0 {
				fd.DBType = "VARCHAR(64)"
				fd.DBValue = fmt.Sprintf("'%s'", fd.F.Tag.Get("default"))
			} else if length > 65535 {
				fd.DBType = fmt.Sprintf("LONGTEXT(%d)", length)
				fd.DBValue = ""
			} else if length > 4096 {
				fd.DBType = fmt.Sprintf("TEXT(%d)", length)
				fd.DBValue = ""
			} else if length == -1 {
				fd.DBType = "TEXT"
				fd.DBValue = ""
			} else if length == -2 {
				fd.DBType = "LONGTEXT"
				fd.DBValue = ""
			} else {
				fd.DBType = fmt.Sprintf("VARCHAR(%d)", length)
				fd.DBValue = fmt.Sprintf("'%s'", fd.F.Tag.Get("default"))
			}
		default:
			fd.DBValue = ""
			if length == 0 {
				fd.DBType = "TEXT"
			} else {
				fd.DBType = fmt.Sprintf("TEXT(%d)", length)
			}
			fd.IsObject = true
		}

		if fd.DBValue != "" {
			fd.DBValue = "DEFAULT " + fd.DBValue
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
