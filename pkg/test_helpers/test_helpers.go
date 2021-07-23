package test_helpers

import (
	"fmt"
	"reflect"
)

func Inspect(val interface{}) string {
	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Array, reflect.Slice:
		s := "["
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				s += ", "
			}
			s += fmt.Sprintf("%v", v.Index(i))
		}
		s += "]"
		return s
	default:
		return fmt.Sprintf("%v", val)
	}
}
