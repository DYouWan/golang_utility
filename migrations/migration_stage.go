package migrations

import (
	"errors"
	"gorm.io/gorm"
)

// MigrationStage 对外提供迁移的口子
type MigrationStage struct {
	Name     string
	Function func(db *gorm.DB, name string) error
}

// Migrate 迁移
func Migrate(db *gorm.DB, migrations []MigrationStage) error {
	for _, m := range migrations {
		if MigrationExists(db, m.Name) {
			continue
		}
		if err := m.Function(db, m.Name); err != nil {
			return err
		}
		if err := SaveMigration(db, m.Name); err != nil {
			return err
		}
	}
	return nil
}

// MigrationExists 检查迁移是否已经运行过
func MigrationExists(db *gorm.DB, migrationName string) bool {
	found := false
	result := db.Where("name = ?", migrationName).First(&MigrationModel{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Info("The %s migration starts", migrationName)
	} else if result.Error != nil {
		panic(result.Error)
	} else {
		found = true
		logger.Info("%s There is a migration record, skip this time", migrationName)
	}
	return found
}

// SaveMigration 记录执行外部的方法
func SaveMigration(db *gorm.DB, migrationName string) error {
	migration := &MigrationModel{Name: migrationName}

	if err := db.Create(migration).Error; err != nil {
		return err
	}
	return nil
}

// MigrateAll 运行引导，然后运行列出的所有迁移函数
func MigrateAll(db *gorm.DB, migrationFunctions []func(*gorm.DB) error) {
	if err := Bootstrap(db); err != nil {
		logger.Error(err)
	}

	for _, m := range migrationFunctions {
		if err := m(db); err != nil {
			logger.Error(err)
		}
	}
}
