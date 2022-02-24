package types

import (
	"encoding/json"
)

// StrJSON 将字符串进行JSON序列化或反序列化
// 空字符串为`null`值，反之亦然
type StrJSON string

func (j StrJSON) MarshalJSON() ([]byte, error) {
	if j == "" {
		return []byte("null"), nil
	}
	if err := j.checkBytes([]byte(j)); err != nil {
		return nil, err
	}

	return []byte(j), nil
}

func (j *StrJSON) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*j = ""
		return nil
	}
	if err := j.checkBytes(data); err != nil {
		return err
	}
	*j = StrJSON(data)
	return nil
}

func (j StrJSON) checkBytes(data []byte) error {
	var tmp interface{}
	return json.Unmarshal(data, &tmp)
}
