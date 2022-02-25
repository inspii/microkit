package resource

import (
	"regexp"
	"strconv"
	"strings"
)

// 平面：
// 		id,name
// 嵌套：
// 		id,name,user{id,name}
//     	id,name,user:address{id,city,street}
// 嵌套分页：
//      id,name,user.offset(1).limit(5).sort(-id){id,name}
var fieldExp = regexp.MustCompile(`([\w]+[:\w]*|\*)(\.offset\([0-9]+\))?(\.limit\([0-9]+\))?(\.sort\([\w-]+\))?({(\w,*|\*,*)*})?`)

type node struct {
	Name        string
	Fields      []string
	Offset      int
	Limit       int
	Sorts       []string
	HasChildren bool
	Children    []*node
}

func parseFields(fields []string, offset int, limit int, sorts []string) *node {
	return parse("", fields, offset, limit, sorts)
}

func parse(nodeName string, fields []string, offset int, limit int, sorts []string) *node {
	currentNode := node{
		Name:   nodeName,
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
			children[nodeName] = append(children[nodeName], child)
		}
	}

	for name, nodes := range children {
		var childFields []string
		var childOffset, childLimit int
		var childSorts []string
		for _, node := range nodes {
			if !node.HasChildren {
				childOffset = node.Offset
				childLimit = node.Limit
				childSorts = node.Sorts
			}
			childFields = append(childFields, node.Fields...)
		}
		child := parse(name, childFields, childOffset, childLimit, childSorts)
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
			child.Offset, _ = strconv.Atoi(offset)
		}
		if limit, ok := args["limit"]; ok {
			child.Limit, _ = strconv.Atoi(limit)
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
