package types

import (
	"strconv"
	"strings"
)

type StrValue string

func (v StrValue) IsEmpty() bool {
	return v == ""
}

func (v StrValue) Default(value string) StrValue {
	if v == "" {
		return StrValue(value)
	}
	return v
}

func (v StrValue) String() string {
	return string(v)
}

func (v StrValue) Int() (int, error) {
	return strconv.Atoi(string(v))
}

func (v StrValue) Uint() (uint, error) {
	num, err := strconv.ParseUint(string(v), 10, 0)

	return uint(num), err
}

func (v StrValue) Int8() (int8, error) {
	num, err := strconv.ParseInt(string(v), 10, 8)

	return int8(num), err
}

func (v StrValue) Uint8() (uint8, error) {
	num, err := strconv.ParseUint(string(v), 10, 8)

	return uint8(num), err
}

func (v StrValue) Int16() (int16, error) {
	num, err := strconv.ParseInt(string(v), 10, 16)

	return int16(num), err
}

func (v StrValue) Uint16() (uint16, error) {
	num, err := strconv.ParseUint(string(v), 10, 16)

	return uint16(num), err
}

func (v StrValue) Int32() (int32, error) {
	num, err := strconv.ParseInt(string(v), 10, 32)

	return int32(num), err
}

func (v StrValue) Uint32() (uint32, error) {
	num, err := strconv.ParseUint(string(v), 10, 32)

	return uint32(num), err
}

func (v StrValue) Int64() (int64, error) {
	return strconv.ParseInt(string(v), 10, 64)
}

func (v StrValue) Uint64() (uint64, error) {
	return strconv.ParseUint(string(v), 10, 64)
}

func (v StrValue) Bool() (bool, error) {
	return strconv.ParseBool(string(v))
}

func (v StrValue) Float32() (float32, error) {
	num, err := strconv.ParseFloat(string(v), 32)

	return float32(num), err
}

func (v StrValue) Float64() (float64, error) {
	return strconv.ParseFloat(string(v), 64)
}

func (v StrValue) StringArray() []string {
	arr := strings.Split(string(v), ",")
	resultArr := make([]string, 0, len(arr))
	for _, item := range arr {
		if strings.TrimSpace(item) != "" {
			resultArr = append(resultArr, item)
		}
	}
	return resultArr
}
