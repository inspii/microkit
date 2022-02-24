package envelope

import (
	"github.com/inspii/microkit/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	pageName     = "page"
	pageSizeName = "page_size"
	offsetName   = "offset"
	limitName    = "limit"
)

var (
	defaultPageSize = 20
)

// GetParam 从URL中获取参数
func GetParam(c *gin.Context, key string) types.StrValue {
	return types.StrValue(c.Param(key))
}

// GetQuery 从URL中获取参数
func GetQuery(c *gin.Context, key string) types.StrValue {
	return types.StrValue(c.Query(key))
}

type Pagination struct {
	Offset   int
	Limit    int
	Page     int
	PageSize int
}

// SetDefaultPageSize 设置默认分页大小
func SetDefaultPageSize(pageSize int) {
	defaultPageSize = pageSize
}

// GetPagination 获取分页信息
func GetPagination(c *gin.Context) Pagination {
	if c.Query(offsetName) != "" || c.Query(limitName) != "" {
		offset, limit := GetOffsetLimit(c)
		page, pageSize := ToPagination(offset, limit)
		return Pagination{
			Offset:   offset,
			Limit:    limit,
			Page:     page,
			PageSize: pageSize,
		}
	} else {
		page, pageSize := GetPaging(c)
		offset, limit := ToOffsetLimit(page, pageSize)
		return Pagination{
			Offset:   offset,
			Limit:    limit,
			Page:     page,
			PageSize: pageSize,
		}
	}
}

// GetPaging 获取分页信息
func GetPaging(c *gin.Context) (page, pageSize int) {
	page = toInt(c.Query(pageName), 1)
	if page < 1 {
		page = 1
	}

	pageSize = toInt(c.Query(pageSizeName), defaultPageSize)
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return
}

// GetOffsetLimit 获取分页信息
func GetOffsetLimit(c *gin.Context) (offset, limit int) {
	offset = toInt(c.Query(offsetName), 0)
	if offset < 0 {
		offset = 0
	}

	limit = toInt(c.Query(limitName), defaultPageSize)
	if limit <= 0 {
		limit = defaultPageSize
	}

	return
}

// ToOffsetLimit 分页信息转换
func ToOffsetLimit(page, pageSize int) (offset, limit int) {
	offset = (page - 1) * pageSize
	limit = pageSize
	return
}

// ToPagination 分页信息转换
func ToPagination(offset, limit int) (page, pageSize int) {
	if limit <= 0 {
		limit = defaultPageSize
	}

	page = offset/limit + 1
	pageSize = limit
	return
}

func toInt(val string, defaultVal int) int {
	if val == "" {
		return defaultVal
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return num
}
