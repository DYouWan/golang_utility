package migrations

import "gorm.io/gorm"

// Migration 单个数据库迁移
type Migration struct {
	gorm.Model
	Name string `sql:"size:255"`
}

