package resource

import (
	"net/http"
	"net/url"
	"testing"
)

var ResultReservedWords = []string{"offset", "limit", "sort"}

func TestParseFilters(t *testing.T) {
	req := http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.test.com/search?fields=id,title&offset=0&limit=10&sort=-pop,-id&name[like]=liam")

	values := req.URL.Query()
	filters := ParseFilters(&values, ResultReservedWords)
	if len(filters) != 1 {
		t.Errorf("expected only 1 filter, but got: %d\n", len(filters))
	}
	if filters[0].Name != "name" || filters[0].Comparator != "[like]" || filters[0].Value != "liam" {
		t.Errorf("expected Name=name, Comparator=[like], Value=liam, but got: %v\n", filters[0])
	}
}

func TestParseFiltersEqual(t *testing.T) {
	req := http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.test.com/search?phone=123")

	values := req.URL.Query()
	filters := ParseFilters(&values, ResultReservedWords)
	if len(filters) != 1 {
		t.Errorf("expected only 1 filter, but got: %d\n", len(filters))
	}
	if filters[0].Name != "phone" || filters[0].Comparator != "[eq]" || filters[0].Value != "123" {
		t.Errorf("expected Name=phone, Comparator=[eq], Value=123, but got: %v\n", filters[1])
	}
}

func TestParseFiltersEmbed(t *testing.T) {
	req := http.Request{Method: "GET"}
	req.URL, _ = url.Parse("http://www.test.com/search?user.phone=123")

	values := req.URL.Query()
	filters := ParseFilters(&values, ResultReservedWords)
	if len(filters) != 1 {
		t.Errorf("expected only 1 filter, but got: %d\n", len(filters))
	}
	if filters[0].Name != "user.phone" || filters[0].Comparator != "[eq]" || filters[0].Value != "123" {
		t.Errorf("expected Name=user.phone, Comparator=[eq], Value=123, but got: %v\n", filters[1])
	}
}
