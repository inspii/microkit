package resource

import (
	"fmt"
	"strings"
	"testing"
)

func TestExtractArg(t *testing.T) {
	if name, value := extractArg("offset(10)"); name != "offset" || value != "10" {
		t.Errorf("expected offset and 10, but got %s and %s", name, value)
	}
	if name, value := extractArg(".limit(10)"); name != "limit" || value != "10" {
		t.Errorf("expected limit and 10, but got %s and %s", name, value)
	}
	if name, value := extractArg("{id,fields}"); name != "fields" || value != "id,fields" {
		t.Errorf(`expected "fields" and "id,fields", but got %s and %s`, name, value)
	}
}

func TestChild(t *testing.T) {
	var childNode node

	childNode = child("comments{id,content}")
	if childNode.Name != "comments" ||
		childNode.Offset != 0 ||
		childNode.Limit != 0 ||
		!arrayStringCompare([]string{"id", "content"}, childNode.Fields) {

		t.Errorf(`expected comments, 0, 0, "id, content"; but got: %s, %d, %d, %v`,
			childNode.Name, childNode.Offset, childNode.Limit, childNode.Fields)
	}

	childNode = child("comments.offset(10).limit(10){id,content}")
	if childNode.Name != "comments" ||
		childNode.Offset != 10 ||
		childNode.Limit != 10 ||
		!arrayStringCompare([]string{"id", "content"}, childNode.Fields) {

		t.Errorf(`expected comments, 10, 10, "id, content"; but got: %s, %d, %d, %v`,
			childNode.Name, childNode.Offset, childNode.Limit, childNode.Fields)
	}

	childNode = child("comments:categories.offset(10).limit(10){id,title}")
	if childNode.Name != "comments" ||
		childNode.Offset != 0 ||
		childNode.Limit != 0 ||
		!arrayStringCompare([]string{"categories.offset(10).limit(10){id,title}"}, childNode.Fields) {

		t.Errorf(`expected comments, 0, 0, "categories.offset(10).limit(10){id,title}"; but got: %s, %d, %d, %v`,
			childNode.Name, childNode.Offset, childNode.Limit, childNode.Fields)
	}
}

func TestParseFields(t *testing.T) {
	fields := []string{"id", "title", "comments{id,content,*}", "comments:categories.offset(10).limit(10){id,title}"}
	aNode := parse("post", fields, 0, 10, []string{})
	printNode(aNode, 0)
	assertNode(t, aNode, []node{
		{
			Name:   "post",
			Fields: []string{"id", "title"},
			Offset: 0,
			Limit:  10,
			Sorts:  []string{},
		},
		{
			Name:   "comments",
			Fields: []string{"id", "content", "*"},
			Offset: 0,
			Limit:  0,
			Sorts:  []string{},
		},
		{
			Name:   "categories",
			Fields: []string{"id", "title"},
			Offset: 10,
			Limit:  10,
			Sorts:  []string{},
		},
	})
}

func arrayStringCompare(array1, array2 []string) bool {
	len1 := len(array1)
	len2 := len(array2)
	if len1 != len2 {
		return false
	}

	for i := 0; i < len1; i++ {
		if array1[i] != array2[i] {
			return false
		}
	}

	return true
}

func assertNode(t *testing.T, node *node, structures []node) []node {
	if structures == nil || len(structures) == 0 {
		return nil
	}

	structure := structures[0]
	if node == nil {
		t.Errorf("expected %v, but got nil", structure)
	}
	if node.Name != structure.Name || node.Offset != structure.Offset || node.Limit != structure.Limit ||
		!arrayStringCompare(node.Fields, structure.Fields) || !arrayStringCompare(node.Sorts, structure.Sorts) {
		t.Errorf("expected %v but got %v", structure, *node)
	}

	structures = structures[1:]
	for _, child := range node.Children {
		structures = assertNode(t, child, structures)
	}

	return structures
}

func printNode(node *node, paddingSpace int) {
	padding := strings.Repeat(" ", paddingSpace)
	fmt.Printf("%s{\n", padding)
	fmt.Printf("%sName:%s\n", padding, node.Name)
	fmt.Printf("%sFileds:%v\n", padding, node.Fields)
	fmt.Printf("%sOffset:%d\n", padding, node.Offset)
	fmt.Printf("%sLimit:%d\n", padding, node.Limit)
	fmt.Printf("%sSorts:%v\n", padding, node.Sorts)
	fmt.Printf("%sChildren:\n", padding)
	for _, child := range node.Children {
		printNode(child, paddingSpace+4)
	}
	fmt.Printf("%s}\n", padding)
}
