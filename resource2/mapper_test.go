package resource

import (
	"reflect"
	"strconv"
	"testing"
)

func TestToMap(t *testing.T) {
	comment := struct {
		Id        int
		Content   string
		UserId    int
		CreatedAt string
	}{
		Id:        3,
		Content:   "hello",
		UserId:    5,
		CreatedAt: "2017-01-01",
	}

	fields := []string{"Id", "Content"}
	formatted, err := toMapWithRelation(comment, fields)
	if err != nil {
		t.Error(err)
	}

	commentMap, ok := formatted.(map[string]interface{})
	if !ok {
		t.Error("commentMap is not map[stirng]interface{}")
	}
	for key := range commentMap {
		if !inArrayString(key, fields) {
			t.Errorf("not expected field: %s", key)
		}
	}
}

func TestToMapOfSlice(t *testing.T) {
	type comment struct {
		Id        int
		Content   string
		UserId    int
		CreatedAt string
	}

	comments := []comment{
		{Id: 3,
			Content:   "hello",
			UserId:    5,
			CreatedAt: "2017-01-01",
		},
		{Id: 4,
			Content:   "world",
			UserId:    5,
			CreatedAt: "2017-01-02",
		},
	}

	fields := []string{"Id", "Content"}
	result, err := toMapWithRelation(comments, fields)
	if err != nil {
		t.Error(err)
	}

	r := reflect.ValueOf(result)
	if r.Kind() != reflect.Slice {
		t.Error("result is not a slice")
	}

	for i := 0; i < r.Len(); i++ {
		item := r.Index(i).Interface()
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			t.Error("itemMap is not a slice")
		}
		for key := range itemMap {
			if !inArrayString(key, fields) {
				t.Errorf("unexpected key: %s", key)
			}
		}
	}
}

type user struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type post struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserId  int    `json:"userId"`
}

func (p post) User() interface{} {
	return user{
		Id:   p.UserId,
		Name: "user:" + strconv.Itoa(p.UserId),
	}
}

type comment struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	PostId  int    `json:"postId"`
	UserId  int    `json:"userId"`
}

func (c comment) User() interface{} {
	return user{
		Id:   c.UserId,
		Name: "user:" + strconv.Itoa(c.UserId),
	}
}

func (c comment) Post() interface{} {
	return post{
		Id:      c.PostId,
		Title:   "post:" + strconv.Itoa(c.PostId),
		Content: "content",
		UserId:  999,
	}
}

func TestLoadRelation(t *testing.T) {
	comment := comment{
		Id:      3,
		Content: "Hehe",
		PostId:  10,
		UserId:  5,
	}
	node := &node{
		Name:   "comment",
		Fields: []string{"id", "content"},
	}

	result, err := loadRelation(comment, node)
	if err != nil {
		t.Error(err)
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Error("commentMap is not map[stirng]interface{}")
	}

	for key := range resultMap {
		if !inArrayString(key, []string{"id", "content"}) {
			t.Errorf("unexpected key: %s", key)
		}
	}
}

func TestLoadRelationArray(t *testing.T) {
	comments := []comment{
		{
			Id:      3,
			Content: "Hehe",
			PostId:  10,
			UserId:  5,
		},
		{
			Id:      4,
			Content: "Hi",
			PostId:  8,
			UserId:  6,
		},
	}
	node := &node{
		Name:   "comment",
		Fields: []string{"id", "content"},
	}

	result, err := loadRelation(comments, node)
	if err != nil {
		t.Error(err)
	}

	r := reflect.ValueOf(result)
	if r.Kind() != reflect.Slice {
		t.Error("result is not a slice")
	}

	for i := 0; i < r.Len(); i++ {
		item := r.Index(i).Interface()
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			t.Error("itemMap is not a slice")
		}
		for key := range itemMap {
			if !inArrayString(key, []string{"id", "content"}) {
				t.Errorf("unexpected key: %s", key)
			}
		}
	}
}

func TestLoadRelationEmbed(t *testing.T) {
	comment := comment{
		Id:      3,
		Content: "Hehe",
		PostId:  10,
		UserId:  5,
	}
	node := &node{
		Name:   "comment",
		Fields: []string{"id", "content"},
		Children: []*node{
			{
				Name:   "post",
				Fields: []string{"id", "title", "content"},
				Children: []*node{
					{
						Name:   "user",
						Fields: []string{"id", "name"},
					},
				},
			},
		},
	}

	result, err := loadRelation(comment, node)
	if err != nil {
		t.Error(err)
	}
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Error("commentMap is not map[stirng]interface{}")
	}
	assertField(t, result, []string{"id", "content", "post"})

	post, ok := resultMap["post"]
	if !ok {
		t.Error("expected key post exists, but got none")
	}
	assertField(t, post, []string{"id", "title", "content", "user"})

	ipostMap, ok := post.(map[string]interface{})
	if !ok {
		t.Errorf("expected type map[string]interface{}, but got: %s", reflect.TypeOf(post))
	}
	user, ok := ipostMap["user"]
	if !ok {
		t.Error("node user of post not found")
	}
	assertField(t, user, []string{"id", "name"})
}

func assertField(t *testing.T, ele interface{}, fields []string) {
	if eleMap, ok := ele.(map[string]interface{}); ok {
		for key := range eleMap {
			if !inArrayString(key, fields) {
				t.Errorf("unexpected key: %s", key)
			}
		}
	}
}
