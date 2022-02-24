package types

import (
	"database/sql/driver"
	"errors"
)

// SQLJSON JSON格式数据
type SQLJSON string

func (j SQLJSON) Value() (driver.Value, error) {
	return string(j), nil
}

func (j *SQLJSON) Scan(args interface{}) error {
	switch args.(type) {
	case []byte:
		*j = SQLJSON(args.([]byte))
		return nil
	case string:
		*j = SQLJSON(args.(string))
		return nil
	default:
		return errors.New("not string")
	}
}
