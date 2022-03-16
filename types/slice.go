package types

import (
	"reflect"
	"strconv"
	"strings"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// StrInts 以","分隔的整数数组
// 常用于HTTP请求中表示id列表，或数据库字段存放id列表，
type StrInts string

// NewStrInts 从切片新建 StrInts
func NewStrInts[T Integer](slice []T) StrInts {
	if len(slice) == 0 {
		return ""
	}

	str := ""
	for _, v := range slice {
		str += strconv.Itoa(int(v)) + ","
	}
	return StrInts(str[0 : len(str)-1])
}

func (s StrInts) String() string {
	return string(s)
}

// Slice 将 StrInts 转化成切片
func (s StrInts) Slice() []int {
	splits := strings.Split(string(s), ",")
	slice := make([]int, 0, len(splits))
	for _, v := range splits {
		if v == "" {
			continue
		}
		if num, err := strconv.Atoi(v); err == nil {
			slice = append(slice, num)
		}
	}
	return slice
}

// Each 对任意类型的切片或数组，进行遍历执行
func Each[T any](slice []T, f func(i int)) {
	rv := reflect.ValueOf(slice)
	size := rv.Len()

	for i := 0; i < size; i++ {
		f(i)
	}
}

// Ints 从任意类型的切片或数组，生成Int数组
func Ints[T any](slice []T, f func(i int) int) []int {
	rv := reflect.ValueOf(slice)
	size := rv.Len()
	ints := make([]int, 0, size)
	for i := 0; i < size; i++ {
		ints = append(ints, f(i))
	}
	return ints
}

// IntsMap 从任意类型的切片或数组，生成映射表
func IntsMap[T any](slice []T, mapper func(i int) (key, value int)) map[int][]int {
	result := make(map[int][]int)

	rv := reflect.ValueOf(slice)
	size := rv.Len()
	for i := 0; i < size; i++ {
		k, v := mapper(i)
		result[k] = append(result[k], v)
	}
	return result
}
