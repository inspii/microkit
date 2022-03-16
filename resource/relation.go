package resource

import "reflect"

type relationEntry struct {
	ID    int
	Index int
}

// Relation 处理多对多关系
//
// 例如：
//  有书籍（Book），作者（Author）及书籍作者关系（BookAuthor）对象，需要给书籍对象填入所有关联的作者对象：
//
//  type Author struct {
//  	ID int
//  	...
//  }
//
//  type Book struct {
//  	ID int
//  	...
//  	Authors []Author
//  }
//
//  type BookAuthor struct {
//  	BookID		int
//  	AuthorID	int
//  }
//
// 	Relation(
//		books, func(i int) int {
//			return books[i].ID
//		},
//		authors, func(i int) int {
//			return authors[i].ID
//		},
//		booksAuthors, func(i int) (int, int) {
//			return booksAuthors[i].BookID, booksAuthors[i].AuthorID
//		},
//		func(i int, j int) {
//			books[i].Authors = append(books[i].Authors, authors[j])
//		},
//	)
func Relation[T1 any, T2 any, T3 any](
	a []T1, getIDOfA func(i int) (aID int),
	b []T2, getIDOfB func(i int) (bID int),
	relation []T3, getRelation func(i int) (aID, bID int),
	do func(aIndex, bIndex int),
) {
	aIDMap := make(map[int]relationEntry)
	rva := reflect.ValueOf(a)
	for i := 0; i < rva.Len(); i++ {
		id := getIDOfA(i)
		aIDMap[id] = relationEntry{
			ID:    id,
			Index: i,
		}
	}

	bIDMap := make(map[int]relationEntry)
	rvb := reflect.ValueOf(b)
	for i := 0; i < rvb.Len(); i++ {
		id := getIDOfB(i)
		bIDMap[id] = relationEntry{
			ID:    id,
			Index: i,
		}
	}

	abMap := make(map[relationEntry][]relationEntry)
	rvr := reflect.ValueOf(relation)
	for i := 0; i < rvr.Len(); i++ {
		aID, bID := getRelation(i)
		aIndex, ok := aIDMap[aID]
		if !ok {
			continue
		}
		bIndex, ok := bIDMap[bID]
		if !ok {
			continue
		}
		abMap[aIndex] = append(abMap[aIndex], bIndex)
	}

	for aEntry, bEntries := range abMap {
		for _, bEntry := range bEntries {
			aIndex := aEntry.Index
			bIndex := bEntry.Index
			do(aIndex, bIndex)
		}
	}
}
