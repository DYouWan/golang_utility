package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"utility/config"
)

func NewDataBase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	if cfg.Type == "mysql" {
		return initMySql(cfg)
	}
	return nil, fmt.Errorf("无效的数据库类型")
}

func initMySql(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	args := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DatabaseName,
	)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return db, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenCons)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleCons)

	return db, nil
}
