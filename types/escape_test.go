package types

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Author struct {
	Name   string
	Phones []string
}

type Book struct {
	Title   string
	Authors []*Author
}

func TestEscapeNilSlice(t *testing.T) {
	b := &Book{}
	err1 := EscapeNilSlice(b)
	bs, err2 := json.Marshal(b)
	fmt.Println(err1, err2, string(bs))
}

func TestEscapeNilSliceEmbed(t *testing.T) {
	liam := &Author{
		Name: "liam",
	}
	b := &Book{
		Title:   "Jane Eyre",
		Authors: []*Author{liam},
	}
	err1 := EscapeNilSlice(b)
	bs, err2 := json.Marshal(b)
	fmt.Println(err1, err2, string(bs))
}

func TestEscapeNilSliceMap(t *testing.T) {
	liam := &Author{
		Name: "liam",
	}
	b := &Book{
		Authors: []*Author{liam},
	}

	bookMap := map[string]*Book{
		"Jane Eyre": b,
	}

	err1 := EscapeNilSlice(bookMap)
	bs, err2 := json.Marshal(bookMap)
	fmt.Println(err1, err2, string(bs))
}
