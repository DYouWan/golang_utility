package migrations

import "gorm.io/gorm"

// MigrationModel 单个数据库迁移
type MigrationModel struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);not null;comment:名称"`
}

func (M MigrationModel) TableName() string {
	return "migration"
}
