package envelope

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/inspii/microkit/types"
)

const (
	errorCodeOK      = 0  // 正常
	errorCodeUnknown = -1 // 未知错误
)

// NamedData 支持自定义名称的响应 Data 或 Extra
type NamedData interface {
	DataName() string
}

// 响应封装
type responseEnvelope struct {
	Code    int
	Message string
	Data    interface{}
	Extras  []interface{}
}

// MarshalJSON 响应封装自定义序列化方式，如
//   {
//     "code": 0,
//     "message": "",
//     "total_count": 66, # 由 Extras 定义
//     "books": []        # 由 Data 定义
//   }
// Data 和 Extras 可通过实现 NamedData 接口，对键名进行自定义。
func (e responseEnvelope) MarshalJSON() ([]byte, error) {
	envelopeMap := make(map[string]interface{})
	envelopeMap["code"] = e.Code
	envelopeMap["message"] = e.Message

	if data, ok := e.Data.(NamedData); ok {
		dataName := data.DataName()
		envelopeMap[dataName] = e.Data
	} else {
		dataName := "data"
		if !types.IsNil(e.Data) {
			envelopeMap[dataName] = e.Data
		}
	}

	var unnamedExtras []interface{}
	for _, extra := range e.Extras {
		if extra == nil {
			continue
		}
		mapped, ok := toMap(extra)
		if ok {
			for k, v := range mapped {
				envelopeMap[k] = v
			}
		} else {
			if data, ok := e.Data.(NamedData); ok {
				extraName := data.DataName()
				envelopeMap[extraName] = extra
			} else {
				unnamedExtras = append(unnamedExtras, extra)
			}
		}
	}
	if len(unnamedExtras) > 0 {
		envelopeMap["extras"] = unnamedExtras
	}

	_ = types.EscapeNilSlice(&envelopeMap)
	return json.Marshal(envelopeMap)
}

// OK 响应成功
func OK(ctx *gin.Context, data interface{}, extras ...interface{}) {
	ctx.JSON(http.StatusOK, responseEnvelope{
		Code:   errorCodeOK,
		Data:   data,
		Extras: extras,
	})
}

// Error 响应失败
func Error(ctx *gin.Context, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, responseEnvelope{
			Code: errorCodeOK,
		})
		return
	}

	httpCode := http.StatusOK
	if err, ok := err.(types.ErrorWithHTTPCode); ok {
		httpCode = err.HTTPCode()
	}

	code := errorCodeUnknown
	if err, ok := err.(types.ErrorWithCode); ok {
		code = err.Code()
	}

	ctx.JSON(httpCode, responseEnvelope{
		Code:    code,
		Message: err.Error(),
	})
}

// toMap 将数据转换成字典（仅支持NamedData、结构体、部分字典）
func toMap(v interface{}) (map[string]interface{}, bool) {
	if data, ok := v.(NamedData); ok {
		return map[string]interface{}{
			data.DataName(): v,
		}, true
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map, reflect.Struct:
		// 使用JSON序列化，再反序列化成字典的形式，确保能正常解析结构体的所有 JSON Tag。
		// 仅使用反射方式较难完全支持结构体的所有JSON Tag，且无法支持 encoding/json 包的后期特性。
		mapped := make(map[string]interface{})
		bytes, err := json.Marshal(v)
		if err != nil {
			return nil, false
		}
		if err := json.Unmarshal(bytes, &mapped); err != nil {
			return nil, false
		}
		return mapped, true
	default:
		return nil, false
	}
}
