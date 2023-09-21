package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// NewMySQL 创建 MySQL 连接实例
func NewMySQL(dsn string) *gorm.DB {
	m, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if db, err := m.DB(); err != nil {
		panic(err)
	} else {
		db.SetConnMaxLifetime(5 * time.Minute)
	}
	return m
}
