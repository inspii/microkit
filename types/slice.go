package types

import (
	"reflect"
	"strconv"
	"strings"
)

// StrInts 以","分隔的整数数组
// 常用于HTTP请求中表示id列表，或数据库字段存放id列表，
type StrInts string

// NewStrInts 从切片新建 StrInts
func NewStrInts(slice []int) StrInts {
	if len(slice) == 0 {
		return ""
	}

	str := ""
	for _, v := range slice {
		str += strconv.Itoa(v) + ","
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
// 如果不是切片或数组，将奔溃，而不隐瞒该错误
func Each(slice interface{}, f func(i int)) {
	rv := reflect.ValueOf(slice)
	size := rv.Len()

	for i := 0; i < size; i++ {
		f(i)
	}
}

// Ints 从任意类型的切片或数组，生成Int数组
// 如果不是切片或数组，将奔溃，而不隐瞒该错误
func Ints(slice interface{}, f func(i int) int) []int {
	rv := reflect.ValueOf(slice)
	size := rv.Len()
	ints := make([]int, 0, size)
	for i := 0; i < size; i++ {
		ints = append(ints, f(i))
	}
	return ints
}

// IntsMap 从任意类型的切片或数组，生成映射表
// 如果不是切片或数组，将奔溃，而不隐瞒该错误
func IntsMap(slice interface{}, mapper func(i int) (key, value int)) map[int][]int {
	result := make(map[int][]int)

	rv := reflect.ValueOf(slice)
	size := rv.Len()
	for i := 0; i < size; i++ {
		k, v := mapper(i)
		result[k] = append(result[k], v)
	}
	return result
}
