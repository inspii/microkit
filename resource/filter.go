package resource

import (
	"net/url"
	"regexp"
)

const (
	Like             = "[like]"
	NotLike          = "[not_like]"
	In               = "[in]"
	NotIn            = "[not_in]"
	GreaterThan      = "[gt]"
	EqualGreaterThan = "[egt]"
	LessThan         = "[lt]"
	EqualLessThan    = "[elt]"
	Equal            = "[eq]"
	NotEqual         = "[not]"
	Between          = "[between]"
)

var comparators = []string{
	Like,
	NotLike,
	In,
	NotIn,
	GreaterThan,
	EqualGreaterThan,
	LessThan,
	EqualLessThan,
	Equal,
	NotEqual,
	Between,
}

var filterExp = regexp.MustCompile(`^([a-zA-Z1-9.]+)(\[[a-zA-Z1-9]+\])*$`)

func IsComparator(comparator string) bool {
	for _, one := range comparators {
		if one == comparator {
			return true
		}
	}
	return false
}

type Filter struct {
	Name       string
	Comparator string
	Value      string
}

func ParseFilters(values *url.Values, exceptField []string) []Filter {
	filters := make([]Filter, 0, len(*values))
	for param, valueList := range *values {
		if inArrayString(param, exceptField) {
			continue
		}
		matches := filterExp.FindStringSubmatch(param)
		if matches == nil {
			continue
		}

		var comparator string
		if matches[2] == "" {
			comparator = Equal
		} else {
			comparator = matches[2]
		}
		if !IsComparator(comparator) {
			continue
		}
		filter := Filter{
			Name:       matches[1],
			Value:      valueList[0],
			Comparator: comparator,
		}

		filters = append(filters, filter)
	}

	return filters
}

func inArrayString(needle string, array []string) bool {
	if len(array) == 0 {
		return true
	}
	for _, item := range array {
		if item == "*" {
			return true
		}
	}
	for _, item := range array {
		if needle == item {
			return true
		}
	}

	return false
}
