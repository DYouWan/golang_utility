package config

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
