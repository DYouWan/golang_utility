package migrations

import (
	"errors"
	"gorm.io/gorm"
	"utility/log"
)

// MigrationStage ...
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
	migration := new(Migration)

	found := false
	result := db.Where("name = ?", migrationName).First(migration)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.INFO.Printf("开始运行%s迁移", migrationName)
	} else if result.Error != nil {
		panic(result.Error)
	} else {
		found = true
		log.INFO.Printf("跳过 %s 迁移", migrationName)
	}
	return found
}


// SaveMigration 将迁移记录保存到迁移表
func SaveMigration(db *gorm.DB, migrationName string) error {
	migration := new(Migration)
	migration.Name = migrationName

	if err := db.Create(migration).Error; err != nil {
		log.ERROR.Printf("将记录保存到迁移表时出错: %s", err)
		return err
	}
	return nil
}

// MigrateAll 运行引导，然后运行列出的所有迁移函数
func MigrateAll(db *gorm.DB, migrationFunctions []func(*gorm.DB) error) {
	if err := Bootstrap(db); err != nil {
		log.ERROR.Print(err)
	}

	for _, m := range migrationFunctions {
		if err := m(db); err != nil {
			log.ERROR.Print(err)
		}
	}
}


