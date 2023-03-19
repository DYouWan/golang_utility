package migrations

import (
	"fmt"
	"gorm.io/gorm"
)

// Bootstrap 创建migrations表
func Bootstrap(db *gorm.DB) error {
	hasTable := db.Migrator().HasTable(&MigrationModel{})
	if hasTable {
		return nil
	}

	// 创建迁移表
	if err := db.AutoMigrate(&MigrationModel{}); err != nil {
		return fmt.Errorf("创建迁移表时出错: %s", err)
	}

	// 创建迁移表记录
	migration := &MigrationModel{Name: "bootstrap_migrations"}
	if err := db.Create(migration).Error; err != nil {
		return fmt.Errorf("将记录保存到迁移表时出错: %s", err)
	}
	return nil
}
