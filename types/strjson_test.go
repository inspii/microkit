package types

import (
	"encoding/json"
	"testing"
)

func TestJSON(t *testing.T) {
	table := [][]interface{}{
		{"", []byte("null")},
		{"1", []byte("1")},
	}

	for _, c := range table {
		a := StrJSON(c[0].(string))
		bytes, err := json.Marshal(a)
		if err != nil {
			t.Error(err)
		}
		if string(bytes) != string(c[1].([]byte)) {
			t.Error(err)
		}
	}

	table2 := [][]interface{}{
		{[]byte("null"), ""},
		{[]byte("1"), "1"},
	}

	for _, c := range table2 {
		var j StrJSON
		err := json.Unmarshal(c[0].([]byte), &j)
		if err != nil {
			t.Error(err)
		}
		if string(j) != c[1] {
			t.Error(err)
		}
	}
}
