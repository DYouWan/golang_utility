package database

import (
	"database/sql"
	"fmt"
	"github.com/dyouwan/utility/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Connect struct {
	db *gorm.DB
}

// GetConnect 获取数据库连接
func GetConnect(cfg *config.DatabaseConfig) (*Connect, error) {
	if cfg.Type == "mysql" {
		return initMySql(cfg)
	}
	return nil, fmt.Errorf("invalid database type")
}

func initMySql(cfg *config.DatabaseConfig) (*Connect, error) {
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
		return nil, err
	}

	conn := &Connect{db: db}
	err = conn.config(cfg.MaxOpenCons, cfg.MaxIdleCons)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// config 配置连接
func (c *Connect) config(maxOpenCons int, maxIdleCons int) error {
	// 如果连接已经关闭，则返回对应的错误信息
	if c.db == nil {
		return fmt.Errorf("database connection is already closed")
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(maxOpenCons)
	sqlDB.SetMaxIdleConns(maxIdleCons)

	return nil
}

// GetDB 获取gorm.DB
func (c *Connect) GetDB() *gorm.DB {
	// 如果连接已经关闭，则返回对应的错误信息
	if c.db == nil {
		panic("database connection is already closed")
	}
	return c.db
}

// GetSqlDB 获取标准库SqlDB
func (c *Connect) GetSqlDB() *sql.DB {
	// 如果连接已经关闭，则返回对应的错误信息
	if c.db == nil {
		panic("database connection is already closed")
	}
	sqlDB, err := c.db.DB()
	if err != nil {
		panic(err)
	}
	return sqlDB
}

// Close 关闭数据库连接
func (c *Connect) Close() error {
	// 如果连接已经关闭，则返回对应的错误信息
	if c.db == nil {
		return fmt.Errorf("database connection is already closed")
	}

	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}

	// 将连接标记为已关闭
	c.db = nil

	return nil
}
