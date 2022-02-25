package resource

import "encoding/json"

// MarshalJSON
//
// 平面：
// 		id,name
// 嵌套：
// 		id,name,user{id,name}
//     	id,name,user:address{id,city,street}
// 嵌套分页：
//      id,name,user.offset(1).limit(5).sort(-id){id,name}
func MarshalJSON(resource interface{}, fields []string) ([]byte, error) {
	m, err := toMapWithRelation(resource, fields)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}
