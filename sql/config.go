package sql

import "gorm.io/gorm"

type Option func(db *gorm.DB) *gorm.DB

// condition 查询条件项
type condition struct {
	query string
	args  []interface{}
}

// preloadConfig 懒加载配置
type preloadConfig struct {
	column     string
	conditions []interface{}
}

// queryConfig 查询条件
type queryConfig struct {
	preloads   []*preloadConfig
	conditions []*condition
	order      string
	offset     int
	limit      int
	rawOptions []Option
}

// parse 解析查询条件
func (c queryConfig) parse(tx *gorm.DB) *gorm.DB {
	for _, r := range c.rawOptions {
		tx = r(tx)
	}

	for _, c := range c.conditions {
		tx = tx.Where(c.query, c.args...)
	}

	for _, preload := range c.preloads {
		tx = tx.Preload(preload.column, preload.conditions...)
	}

	if c.order != "" {
		tx = tx.Order(c.order)
	}
	if c.limit > 0 {
		tx = tx.Offset(c.offset).Limit(c.limit)
	}
	return tx
}

// apply 应用查询选项
func (c *queryConfig) apply(options ...QueryOption) *queryConfig {
	for _, option := range options {
		option(c)
	}
	return c
}
