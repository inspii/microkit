package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

const (
	timeLayoutDateTime = "2006-01-02 15:04:05"
	timeLayoutDate     = "2006-01-02"
)

var local, _ = time.LoadLocation("Local")

// JSONTime "YYYY-MM-DD HH:MM:SS"格式时间
type JSONTime time.Time

func (t JSONTime) String() string {
	return time.Time(t).In(local).Format(timeLayoutDateTime)
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, t.String())
	return []byte(stamp), nil
}

func (t *JSONTime) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	if str == "" {
		return nil // 忽略空值，业务端调用接口容易迷惑
	}

	t0, err := time.ParseInLocation(timeLayoutDateTime, str, local)
	if err != nil {
		return err
	}

	*t = JSONTime(t0)
	return nil
}

// JSONDate "YYYY-MM-DD"格式时间
type JSONDate time.Time

func (d JSONDate) String() string {
	return time.Time(d).In(local).Format(timeLayoutDate)
}

func (d JSONDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *JSONDate) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}
	if str == "" {
		return nil // 忽略空值，业务端调用接口容易迷惑
	}

	t0, err := time.ParseInLocation(timeLayoutDate, str, local)
	if err != nil {
		return err
	}

	*d = JSONDate(t0)
	return nil
}

// JSONTimestamp 时间戳格式时间
type JSONTimestamp time.Time

func (d JSONTimestamp) String() string {
	return strconv.Itoa(int(time.Time(d).Unix()))
}

func (t JSONTimestamp) MarshalJSON() ([]byte, error) {
	val := strconv.Itoa(int(time.Time(t).Unix()))
	return []byte(val), nil
}

func (t *JSONTimestamp) UnmarshalJSON(bytes []byte) error {
	var num int64
	if err := json.Unmarshal(bytes, &num); err != nil {
		return err
	}
	if num == 0 {
		return nil // 忽略空值
	}

	t0 := time.Unix(num, 0).In(local)
	*t = JSONTimestamp(t0)
	return nil
}
