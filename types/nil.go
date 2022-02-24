package types

import "reflect"

// IsNil 是否为Nil
func IsNil(v interface{}) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		return rv.IsNil()
	}

	return false
}
