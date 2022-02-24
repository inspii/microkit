package resource

type GetOption struct {
	Fields []string
}

type ListOption struct {
	Offset    int
	Limit     int
	Sort      string
	WithCount bool
	Filters   []interface{}
	Fields    []string
}

type CRUD interface {
	List(dest interface{}, option ListOption) (total int, err error)
	GetByID(dest interface{}, id int, option GetOption) error
	Create(model interface{}) error
	Update(model interface{}) error
	DeleteByID(id int) error
}
