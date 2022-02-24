package sql

import (
	"fmt"
	"gorm.io/gorm/logger"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDB 创建数据库连接实例
func NewDB(dbHost string, dbPort int, dbName, dbUser, dbPass string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.WithError(err).Fatalf("db connect error (%s)", dsn)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.WithError(err).Fatalf("db connect error (%s)", dsn)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)
	return db
}
