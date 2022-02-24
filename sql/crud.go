package sql

import (
	"gorm.io/gorm/clause"
	"time"

	"gorm.io/gorm"
)

// GetByID 按ID查询
func GetByID(tx *gorm.DB, model interface{}, id int, options ...QueryOption) error {
	options = append(options, WithCondition("id=?", id))
	return GetOne(tx, model, options...)
}

// GetOne 查询一条数据
func GetOne(tx *gorm.DB, model interface{}, options ...QueryOption) error {
	config := &queryConfig{}
	t := config.apply(options...).parse(tx)
	return t.First(model).Error
}

// List 查询列表
func List(tx *gorm.DB, modelArray interface{}, options ...QueryOption) error {
	config := &queryConfig{}
	tx = config.apply(options...).parse(tx)
	return tx.Find(modelArray).Error
}

// Count 查询数据量
func Count(tx *gorm.DB, model interface{}, options ...QueryOption) (int, error) {
	config := &queryConfig{}
	config.apply(options...)
	config.offset = 0
	config.limit = 0

	tx = config.parse(tx.Model(model))

	var count int64
	err := tx.Count(&count).Error
	return int(count), err
}

// ListWithCount 查询列表及数据量
func ListWithCount(tx *gorm.DB, modelArray interface{}, options ...QueryOption) (int, error) {
	count, err := Count(tx, modelArray, options...)
	if err != nil {
		return 0, err
	}
	if count <= 0 {
		return 0, nil
	}
	err = List(tx, modelArray, options...)
	return count, err
}

// Create 创建记录，不支持级联创建
func Create(tx *gorm.DB, model interface{}) error {
	return tx.Omit(clause.Associations).Create(model).Error
}

// UpdateByID 按ID更新数据，不支持级联更新
func UpdateByID(tx *gorm.DB, model interface{}, id int, options ...QueryOption) error {
	options = append(options, WithCondition("id=?", id))
	return UpdateOne(tx, model, options...)
}

// UpdateOne 更新一条数据，不支持级联更新
func UpdateOne(tx *gorm.DB, model interface{}, options ...QueryOption) error {
	config := &queryConfig{}
	t := config.apply(options...).parse(tx)
	return t.Omit(clause.Associations).Limit(1).Updates(model).Error
}

// SoftDeleteByID 按ID软删除数据
func SoftDeleteByID(tx *gorm.DB, model interface{}, id int, options ...QueryOption) error {
	options = append(options, WithCondition("id=?", id))
	return SoftDeleteOne(tx, model, options...)
}

// SoftDeleteOne 软删除一条数据
func SoftDeleteOne(tx *gorm.DB, model interface{}, options ...QueryOption) error {
	now := time.Now()
	config := &queryConfig{}
	t := config.apply(options...).parse(tx)
	return t.Omit(clause.Associations).Model(model).Limit(1).UpdateColumn("delete_time", now).Error
}

// DeleteByID 按ID删除数据
func DeleteByID(tx *gorm.DB, model interface{}, id int, options ...QueryOption) error {
	options = append(options, WithCondition("id=?", id))
	return DeleteOne(tx, model, options...)
}

// DeleteOne 删除一条数据
func DeleteOne(tx *gorm.DB, model interface{}, options ...QueryOption) error {
	config := &queryConfig{}
	t := config.apply(options...).parse(tx)
	return t.Omit(clause.Associations).Limit(1).Delete(model).Error
}

// Delete 删除数据
// limit 限制删除条数，防止误删，取 -1 则不限制
func Delete(tx *gorm.DB, model interface{}, limit int, options ...QueryOption) error {
	config := &queryConfig{}
	t := config.apply(options...).parse(tx)
	if limit > 0 {
		t = t.Limit(limit)
	}
	return t.Omit(clause.Associations).Delete(model).Error
}
