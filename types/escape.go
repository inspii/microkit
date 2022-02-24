package types

import (
	"errors"
	"reflect"
)

// EscapeNilSlice 将数据中所有的 nil 切片或数组，设置为长度为 0 的空数组，避免前端处理失败。
// 注意：必须使用指针传递，否则无法替换掉nil
func EscapeNilSlice(v interface{}) error {
	return escapeNilSlice(reflect.ValueOf(v))
}

func escapeNilSlice(rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Ptr:
		ele := rv.Elem()
		return escapeNilSlice(ele)
	case reflect.Map:
		for _, k := range rv.MapKeys() {
			item := rv.MapIndex(k)
			if err := escapeNilSlice(item); err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			item := rv.Field(i)
			if err := escapeNilSlice(item); err != nil {
				return err
			}
		}
		return nil
	case reflect.Slice, reflect.Array:
		if !rv.CanSet() {
			return errors.New("value not settable")
		}

		if rv.IsNil() {
			zero := reflect.MakeSlice(rv.Type(), 0, 0)
			rv.Set(zero)
		} else {
			for i := 0; i < rv.Len(); i++ {
				item := rv.Index(i)
				if err := escapeNilSlice(item); err != nil {
					return err
				}
			}
		}
		return nil
	default:
		return nil
	}
}
