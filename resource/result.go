package resource

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var ResultReservedWord = []string{"fields", "offset", "limit", "sort"}

var fieldExp = regexp.MustCompile(`([\w]+[:\w]*|\*)(\.offset\([0-9]+\))?(\.limit\([0-9]+\))?(\.order\([\w-]+\))?({(\w,*|\*,*)*})?`)

type node struct {
	Name        string
	Fields      []string
	Offset      int64
	Limit       int64
	Sorts       []string
	HasChildren bool
}

type Node struct {
	Name     string
	Fields   []string
	Offset   int64
	Limit    int64
	Sorts    []string
	Children []*Node
}

func ParseResult(request *http.Request) *Node {
	sorts := strings.Split(request.URL.Query().Get("sort"), ",")
	offset, _ := strconv.ParseInt(request.URL.Query().Get("offset"), 10, 64)
	limit, _ := strconv.ParseInt(request.URL.Query().Get("limit"), 10, 64)
	if offset <= 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 10
	}

	fields := fieldExp.FindAllString(request.URL.Query().Get("fields"), -1)
	node := ParseFields("", fields, offset, limit, sorts)

	return node
}

func ParseFields(name string, fields []string, offset int64, limit int64, sorts []string) *Node {
	currentNode := Node{
		Name:   name,
		Offset: offset,
		Limit:  limit,
		Sorts:  sorts,
	}

	children := make(map[string][]node)

	for _, field := range fields {
		matches := fieldExp.FindStringSubmatch(field)
		if matches == nil {
			continue
		}
		if matches[3] == "" && matches[4] == "" && matches[5] == "" {
			currentNode.Fields = append(currentNode.Fields, field)
		} else {
			names := strings.Split(matches[1], ":")
			nodeName := names[0]
			child := child(field)

			if _, ok := children[nodeName]; ok {
				children[nodeName] = append(children[nodeName], child)
			} else {
				children[nodeName] = []node{child}
			}
		}
	}

	for name, nodes := range children {
		var childFields []string
		var childOffset, childLimit int64
		var childSorts []string
		for _, node := range nodes {
			if !node.HasChildren {
				childOffset = node.Offset
				childLimit = node.Limit
				childSorts = node.Sorts
			}
			childFields = append(childFields, node.Fields...)
		}
		child := ParseFields(name, childFields, childOffset, childLimit, childSorts)
		currentNode.Children = append(currentNode.Children, child)
	}

	return &currentNode
}

func child(field string) (child node) {
	matches := fieldExp.FindStringSubmatch(field)
	if matches == nil {
		return node{
			Name:   "",
			Fields: []string{},
		}
	}
	names := strings.Split(matches[1], ":")
	if len(names) == 1 {
		args := make(map[string]string)
		for _, arg := range matches[2:] {
			name, value := extractArg(arg)
			args[name] = value
		}
		child := node{
			Name: names[0],
		}
		if offset, ok := args["offset"]; ok {
			child.Offset, _ = strconv.ParseInt(offset, 10, 64)
		}
		if limit, ok := args["limit"]; ok {
			child.Limit, _ = strconv.ParseInt(limit, 10, 64)
		}
		if fields, ok := args["fields"]; ok {
			child.Fields = strings.Split(fields, ",")
		}
		if sorts, ok := args["sort"]; ok {
			child.Sorts = strings.Split(sorts, ",")
		}
		return child
	} else {
		return node{
			Name:        names[0],
			Fields:      []string{strings.Replace(field, names[0]+":", "", 1)},
			HasChildren: true,
		}
	}
}

func extractArg(arg string) (name, value string) {
	exp := regexp.MustCompile(`(\w+)\(([\w]*)\)|{([\w,\-+*]*)}`)
	matches := exp.FindStringSubmatch(arg)

	if len(matches) < 4 {
		return
	}

	name = matches[1]
	value = matches[2]
	if name == "" && value == "" {
		name = "fields"
		value = matches[3]
	}
	return
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
