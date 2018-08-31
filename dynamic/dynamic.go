package dynamic

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type IString interface {
	String() string
}

type IGetter interface {
	GetValue(key string) interface{}
}

type ISetter interface {
	SetValue(key string, value interface{})
}

func Get(object interface{}, key string) interface{} {
	return GetWithAutoCreate(object, key, false)
}

func GetWithAutoCreate(object interface{}, key string, autocreate bool) interface{} {

	if object == nil {
		return nil
	}

	{
		v, ok := object.(IGetter)

		if ok {
			return v.GetValue(key)
		}

	}

	{
		v, ok := object.(map[string]interface{})

		if ok {
			vv, ok := v[key]
			if !ok && autocreate {
				vv = map[interface{}]interface{}{}
				v[key] = vv
			}
			return vv
		}
	}

	{
		v, ok := object.(map[interface{}]interface{})

		if ok {
			vv, ok := v[key]
			if !ok && autocreate {
				vv = map[interface{}]interface{}{}
				v[key] = vv
			}
			return vv
		}
	}

	{
		v, ok := object.([]interface{})

		if ok {
			if key == "@length" {
				return len(v)
			} else if key == "@first" {
				if len(v) > 0 {
					return v[0]
				}
			} else if key == "@last" {
				if len(v) > 0 {
					return v[len(v)-1]
				}
			} else {
				i, err := strconv.Atoi(key)
				if err == nil && i >= 0 && i < len(v) {
					return v[i]
				}
			}

			return nil
		}
	}

	{
		v := reflect.ValueOf(object)
		switch v.Kind() {
		case reflect.Map:
			switch v.Type().Key().Kind() {
			case reflect.String, reflect.Interface:

				vv := v.MapIndex(reflect.ValueOf(key))

				if autocreate && !vv.IsValid() {

					switch v.Type().Elem().Kind() {
					case reflect.Ptr:
						vv = reflect.New(v.Type().Elem().Elem())
						v.SetMapIndex(reflect.ValueOf(key), vv)
					case reflect.Interface:
						vv = reflect.ValueOf(map[interface{}]interface{}{})
						v.SetMapIndex(reflect.ValueOf(key), vv)
					}

				}

				if vv.IsValid() {
					if vv.CanInterface() {
						return vv.Interface()
					} else if vv.CanAddr() {
						return vv.Addr().Interface()
					}
				}
			}
		case reflect.Slice:

			var vv reflect.Value

			if key == "@length" {
				return v.Len()
			} else if key == "@first" {
				if v.Len() > 0 {
					vv = v.Index(0)
				}
			} else if key == "@last" {
				if v.Len() > 0 {
					vv = v.Index(v.Len() - 1)
				}
			} else {
				i, err := strconv.Atoi(key)
				if err == nil && i >= 0 && i < v.Len() {
					vv = v.Index(i)
				}
			}
			if vv.IsValid() && vv.CanInterface() {
				return vv.Interface()
			} else if vv.IsValid() && vv.CanAddr() {
				return v.Addr().Interface()
			}
		case reflect.Ptr:
			switch v.Type().Elem().Kind() {
			case reflect.Struct:
				vv := v.Elem().FieldByName(key)
				switch vv.Kind() {
				case reflect.Ptr:
					if vv.IsNil() {
						if autocreate {
							vv.Set(reflect.New(vv.Type().Elem()))
						} else {
							return nil
						}
					}
					return vv.Interface()
				case reflect.Interface:
					if vv.IsNil() {
						if autocreate {
							vv.Set(reflect.ValueOf(map[interface{}]interface{}{}))
						} else {
							return nil
						}
					}
					return vv.Interface()
				case reflect.Map:
					if vv.IsNil() {
						if autocreate {
							vv.Set(reflect.MakeMap(vv.Type()))
						} else {
							return nil
						}
					}
					return vv.Interface()
				default:
					if vv.IsValid() && vv.CanInterface() {
						return vv.Interface()
					} else if vv.CanAddr() {
						return vv.Addr().Interface()
					}
				}
			}
		}
	}

	return nil
}

func GetWithKeys(object interface{}, keys []string) interface{} {
	if keys == nil || len(keys) == 0 {
		return object
	} else if len(keys) == 1 {
		return Get(object, keys[0])
	} else {
		return GetWithKeys(Get(object, keys[0]), keys[1:])
	}
}

func Each(object interface{}, fn func(key interface{}, value interface{}) bool) {

	if object == nil {
		return
	}

	{
		v, ok := object.(map[string]interface{})

		if ok {
			for key, value := range v {
				if !fn(key, value) {
					break
				}
			}
			return
		}
	}

	{
		v, ok := object.(map[interface{}]interface{})

		if ok {
			for key, value := range v {
				if !fn(key, value) {
					break
				}
			}
			return
		}
	}

	{
		v, ok := object.([]interface{})

		if ok {
			for key, value := range v {
				if !fn(key, value) {
					break
				}
			}
			return
		}
	}

	{
		v := reflect.ValueOf(object)
		switch v.Kind() {
		case reflect.Map:

			for _, key := range v.MapKeys() {
				if key.CanInterface() {
					vv := v.MapIndex(key)
					if vv.CanInterface() {
						if !fn(key.Interface(), vv.Interface()) {
							break
						}
					}
				}
			}
		case reflect.Slice:

			for i := 0; i < v.Len(); i++ {
				vv := v.Index(i)
				if vv.CanInterface() {
					if !fn(i, vv.Interface()) {
						break
					}
				}
			}
		case reflect.Ptr:
			switch v.Type().Elem().Kind() {
			case reflect.Struct:
				t := v.Type().Elem()
				for i := 0; i < t.NumField(); i++ {
					fd := t.Field(i)
					fdv := v.Elem().Field(i)
					if fdv.CanInterface() {
						if !fn(fd.Name, fdv.Interface()) {
							break
						}
					}

				}
			}
		}
	}

}

func IsNil(object interface{}) bool {

	if object == nil {
		return true
	}

	v := reflect.ValueOf(object)

	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return v.Bool() == false
	case reflect.Ptr:
		return v.IsNil()
	}

	return false
}

func IsEmpty(object interface{}) bool {

	if object == nil {
		return true
	}

	v := reflect.ValueOf(object)

	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return v.Bool() == false
	case reflect.Slice, reflect.Array:
		return v.Len() == 0
	case reflect.Map:
		return len(v.MapKeys()) == 0
	case reflect.Ptr:
		return v.IsNil()
	}

	return false
}

func Set(object interface{}, key string, value interface{}) {

	if object == nil {
		return
	}

	{
		v, ok := object.(ISetter)

		if ok {
			v.SetValue(key, value)
			return
		}

	}

	{
		v, ok := object.(map[string]interface{})

		if ok {
			if value == nil {
				delete(v, key)
			} else {
				v[key] = value
			}
			return
		}
	}

	{
		v, ok := object.(map[interface{}]interface{})

		if ok {
			if value == nil {
				delete(v, key)
			} else {
				v[key] = value
			}
			return
		}
	}

	{
		v := reflect.ValueOf(object)
		switch v.Kind() {
		case reflect.Map:
			switch v.Type().Key().Kind() {
			case reflect.String, reflect.Interface:
				if value == nil {
					switch v.Type().Elem().Kind() {
					case reflect.Map, reflect.Ptr, reflect.Interface:
						v.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
					}
				} else {
					vv := reflect.ValueOf(value).Convert(v.Type().Elem())
					if vv.IsValid() {
						v.SetMapIndex(reflect.ValueOf(key), vv)
					}
				}
			}
		case reflect.Ptr:
			switch v.Type().Elem().Kind() {
			case reflect.Interface:
				if v.IsNil() {
					vv := map[interface{}]interface{}{}
					Set(vv, key, value)
					v.Set(reflect.ValueOf(vv))
				} else {
					Set(v.Interface(), key, value)
				}
			case reflect.Struct:
				fd := v.Elem().FieldByName(key)
				if fd.IsValid() && fd.CanSet() {
					if value == nil {
						switch fd.Kind() {
						case reflect.Map, reflect.Ptr, reflect.Interface:
							fd.Set(reflect.ValueOf(value))
						}
					} else {
						switch fd.Kind() {
						case reflect.String:
							fd.SetString(StringValue(value, ""))
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							fd.SetInt(IntValue(value, 0))
						case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
							fd.SetUint(UintValue(value, 0))
						case reflect.Float32, reflect.Float64:
							fd.SetFloat(FloatValue(value, 0))
						case reflect.Bool:
							fd.SetBool(BooleanValue(value, false))
						case reflect.Interface:
							fd.Set(reflect.ValueOf(value))
						case reflect.Ptr:
							vv := reflect.ValueOf(value)
							if fd.Type() == vv.Type() {
								fd.Set(vv)
							} else if BooleanValue(value, false) {
								if fd.IsNil() {
									fd.Set(reflect.New(fd.Type().Elem()))
								}
							}
						}
					}
				}
			}
		}
	}

}

func SetWithKeys(object interface{}, keys []string, value interface{}) {
	if keys == nil || len(keys) == 0 {
		return
	} else if len(keys) == 1 {
		Set(object, keys[0], value)
	} else {
		SetWithKeys(GetWithAutoCreate(object, keys[0], true), keys[1:], value)
	}
}

func StringValue(value interface{}, defaultValue string) string {

	if value == nil {
		return defaultValue
	}

	{
		v, ok := value.(string)
		if ok {
			return v
		}
	}

	{
		v, ok := value.(IString)
		if ok {
			return v.String()
		}
	}

	{
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return fmt.Sprintf("%d", v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return fmt.Sprintf("%d", v.Uint())
		case reflect.Float32, reflect.Float64:
			vv := v.Float()
			if float64(int64(vv)) == vv {
				return fmt.Sprintf("%d", int64(vv))
			}
			return fmt.Sprintf("%f", vv)
		case reflect.Bool:
			if v.Bool() {
				return "true"
			}
			return "false"
		case reflect.Ptr:
			if !v.IsNil() && v.Elem().CanInterface() {
				return StringValue(v.Elem().Interface(), defaultValue)
			}
		}
	}

	return defaultValue

}

func IntValue(value interface{}, defaultValue int64) int64 {

	if value == nil {
		return defaultValue
	}

	{
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int64(v.Uint())
		case reflect.Float32, reflect.Float64:
			return int64(v.Float())
		case reflect.Bool:
			if v.Bool() {
				return int64(1)
			}
			return int64(0)
		case reflect.String:
			s := v.String()
			if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
				vv, _ := strconv.ParseInt(s[2:], 16, 64)
				return vv
			} else if strings.HasPrefix(s, "0") {
				vv, _ := strconv.ParseInt(s[1:], 8, 64)
				return vv
			} else {
				vv, _ := strconv.ParseInt(s, 10, 64)
				return vv
			}
		case reflect.Ptr:
			if !v.IsNil() && v.Elem().CanInterface() {
				return IntValue(v.Elem().Interface(), defaultValue)
			}
		}
	}

	return defaultValue

}

func UintValue(value interface{}, defaultValue uint64) uint64 {

	if value == nil {
		return defaultValue
	}

	{
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return uint64(v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return v.Uint()
		case reflect.Float32, reflect.Float64:
			return uint64(v.Float())
		case reflect.Bool:
			if v.Bool() {
				return uint64(1)
			}
			return uint64(0)
		case reflect.String:
			s := v.String()
			if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
				vv, _ := strconv.ParseUint(s[2:], 16, 64)
				return vv
			} else if strings.HasPrefix(s, "0") {
				vv, _ := strconv.ParseUint(s[1:], 8, 64)
				return vv
			} else {
				vv, _ := strconv.ParseUint(s, 10, 64)
				return vv
			}
		case reflect.Ptr:
			if !v.IsNil() && v.Elem().CanInterface() {
				return UintValue(v.Elem().Interface(), defaultValue)
			}
		}
	}

	return defaultValue

}

func FloatValue(value interface{}, defaultValue float64) float64 {

	if value == nil {
		return defaultValue
	}

	{
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(v.Uint())
		case reflect.Float32, reflect.Float64:
			return v.Float()
		case reflect.Bool:
			if v.Bool() {
				return float64(1)
			}
			return float64(0)
		case reflect.String:
			vv, _ := strconv.ParseFloat(v.String(), 64)
			return vv
		case reflect.Ptr:
			if !v.IsNil() && v.Elem().CanInterface() {
				return FloatValue(v.Elem().Interface(), defaultValue)
			}
		}
	}

	return defaultValue

}

func BooleanValue(value interface{}, defaultValue bool) bool {

	if value == nil {
		return defaultValue
	}

	{
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() != 0 {
				return true
			}
			return false
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v.Uint() != 0 {
				return true
			}
			return false

		case reflect.Float32, reflect.Float64:
			if v.Float() != 0 {
				return true
			}
			return false
		case reflect.Bool:
			return v.Bool()
		case reflect.String:
			vv := v.String()
			if vv == "yes" || vv == "true" || vv == "1" {
				return true
			}
			return false
		case reflect.Ptr:
			if !v.IsNil() && v.Elem().CanInterface() {
				return BooleanValue(v.Elem().Interface(), defaultValue)
			}
		}
	}

	return defaultValue

}

func EachReflect(v reflect.Value, fn func(name string, value reflect.Value) bool) {

	if v.Kind() == reflect.Struct {

		count := v.NumField()
		tp := v.Type()

		for i := 0; i < count; i++ {

			tf := tp.Field(i)
			fd := v.Field(i)

			name := tf.Tag.Get("json")

			if name == "-" {
				continue
			}

			if tf.Type.Kind() == reflect.Struct {
				EachReflect(fd, fn)
				continue
			}

			if name == "" {
				continue
			}

			name = strings.Split(name, ",")[0]

			if !fn(name, fd) {
				break
			}

		}
	}

}

func SetReflectValue(v reflect.Value, value interface{}) {

	switch v.Kind() {
	case reflect.String:
		v.SetString(StringValue(value, ""))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(IntValue(value, 0))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(UintValue(value, 0))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(FloatValue(value, 0))
	case reflect.Bool:
		v.SetBool(BooleanValue(value, false))
	case reflect.Ptr:
		switch v.Type().Elem().Kind() {
		case reflect.Struct:
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}

			EachReflect(v.Elem(), func(name string, v reflect.Value) bool {

				if v.IsValid() {
					SetReflectValue(v, Get(value, name))
				}

				return true
			})
		default:
			SetReflectValue(v.Elem(), value)
		}
	case reflect.Map:
		Each(value, func(key interface{}, value interface{}) bool {
			var kval reflect.Value
			var vval reflect.Value
			switch v.Type().Key().Kind() {
			case reflect.Interface:
				kval = reflect.ValueOf(key)
			case reflect.String:
				kval = reflect.ValueOf(StringValue(key, ""))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				kval = reflect.ValueOf(IntValue(key, 0))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				kval = reflect.ValueOf(UintValue(key, 0))
			case reflect.ValueOf(key).Kind():
				kval = reflect.ValueOf(key)
			}
			switch v.Type().Elem().Kind() {
			case reflect.Interface:
				vval = reflect.ValueOf(value)
			case reflect.String:
				vval = reflect.ValueOf(StringValue(value, ""))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				vval = reflect.ValueOf(IntValue(value, 0))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				vval = reflect.ValueOf(UintValue(value, 0))
			case reflect.Float32, reflect.Float64:
				vval = reflect.ValueOf(FloatValue(value, 0))
			case reflect.Bool:
				vval = reflect.ValueOf(BooleanValue(value, false))
			case reflect.ValueOf(value).Kind():
				vval = reflect.ValueOf(value)
			}
			if kval.IsValid() && kval.IsValid() {
				v.SetMapIndex(kval, vval)
			}
			return true
		})
	case reflect.Interface:
		if v.IsValid() && v.CanSet() {
			if value != nil {
				v.Set(reflect.ValueOf(value))
			}
		}
	}

}

func SetValue(object interface{}, value interface{}) {

	if object == nil {
		return
	}

	v := reflect.ValueOf(object)

	SetReflectValue(v, value)
}

func AddValue(object interface{}, value interface{}) {

	if object == nil {
		return
	}

	v := reflect.ValueOf(object)

	switch v.Kind() {
	case reflect.Ptr:
		switch v.Type().Elem().Kind() {
		case reflect.Slice:

			var vv reflect.Value
			var vval = reflect.ValueOf(value)

			if v.IsNil() {
				vv = reflect.MakeSlice(v.Type().Elem(), 0, 0)
			} else {
				vv = v.Elem()
			}

			switch v.Type().Elem().Elem().Kind() {
			case reflect.Interface:
				vv = reflect.ValueOf(value)
			case reflect.String:
				vv = reflect.AppendSlice(vv, reflect.ValueOf(StringValue(value, "")))
			case reflect.Ptr:
				switch v.Type().Elem().Elem().Elem().Kind() {
				case reflect.Struct:
					vval = reflect.New(v.Type().Elem().Elem().Elem())
					SetValue(vval.Interface(), value)
					vv = reflect.AppendSlice(vv, vval)
				}
			case vval.Kind():
				vv = reflect.AppendSlice(vv, vval)
			}

			v.Set(vv)
		}
	}
}
