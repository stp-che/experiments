package test_helpers

import (
	"fmt"
	"reflect"
)

func Inspect(val interface{}) string {
	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Struct, reflect.Ptr:
		return inspectValue(v)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func inspectValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		s := "["
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				s += ", "
			}
			s += inspectValue(v.Index(i))
		}
		s += "]"
		return s
	case reflect.Struct:
		s := "{"
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				s += ", "
			}
			s += inspectValue(v.Field(i))
		}
		s += "}"
		return s
	case reflect.Ptr:
		return "&" + inspectValue(v.Elem())
	default:
		return fmt.Sprintf("%v", v)
	}
}
