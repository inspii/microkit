package sql

// QueryOption 查询选型
type QueryOption func(config *queryConfig)

const notDeleted = "delete_time IS NULL"

func With(o Option) QueryOption {
	return func(config *queryConfig) {
		config.rawOptions = append(config.rawOptions, o)
	}
}

// WithNotDeleted 未删除数据
func WithNotDeleted() QueryOption {
	return func(config *queryConfig) {
		config.conditions = append(config.conditions, &condition{
			query: notDeleted,
		})
	}
}

// WithPreload 懒加载
func WithPreload(column string, conditions ...interface{}) QueryOption {
	return func(config *queryConfig) {
		config.preloads = append(config.preloads, &preloadConfig{
			column:     column,
			conditions: conditions,
		})
	}
}

// WithCondition 查询条件
func WithCondition(query string, args ...interface{}) QueryOption {
	return func(config *queryConfig) {
		config.conditions = append(config.conditions, &condition{
			query: query,
			args:  args,
		})
	}
}

// WithLimit 分页
func WithLimit(offset, limit int) QueryOption {
	return func(config *queryConfig) {
		config.offset = offset
		config.limit = limit
	}
}

// WithPage 分页
func WithPage(page, pageSize int) QueryOption {
	return func(config *queryConfig) {
		if page <= 0 {
			page = 1
		}
		config.offset = (page - 1) * pageSize
		config.limit = pageSize
	}
}

// WithOrder 排序
func WithOrder(order string) QueryOption {
	return func(config *queryConfig) {
		config.order = order
	}
}
