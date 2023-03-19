package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// Config 配置文件
type Config struct {
	//DatabaseConfig 数据库配置
	DatabaseConfig DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Type         string `yaml:"type"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"databaseName" `
	MaxIdleCons  int    `yaml:"maxIdleCons"`
	MaxOpenCons  int    `yaml:"maxOpenCons"`
}

func NewDatabaseConfig(configFile string) (*DatabaseConfig, error) {
	config := &Config{}

	// 读取YAML文件
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("初始化配置文件失败: %v", err)
	}
	return &config.DatabaseConfig, nil
}
