package model

import "fmt"

// 应用配置结构体
type AppConfig struct {
	Port     int `yaml:"port"`
	Database struct {
		DBType DatabaseType `yaml:"type"`
		// PostgreSQL独占配置
		PgHost     string `yaml:"host"`
		PgPort     int    `yaml:"port"`
		PgUser     string `yaml:"username"`
		PgPassword string `yaml:"password"`
		PgDbname   string `yaml:"dbname"`
		PgSSL      bool   `yaml:"ssl"`
	} `yaml:"database"`
	Debug bool `yaml:"debug"`
}

// 数据库类型枚举
type DatabaseType string

const (
	Sqlite     DatabaseType = "sqlite"
	PostgreSql DatabaseType = "postgresql"
)

// 实现 yaml.Unmarshaler 接口
func (dt *DatabaseType) UnmarshalYAML(unmarshal func(any) error) error {
	var raw string
	if err := unmarshal(&raw); err != nil {
		return err
	}

	val := DatabaseType(raw)
	switch val {
	case Sqlite, PostgreSql:
		*dt = val
		return nil
	default:
		return fmt.Errorf("invalid database type '%s', expected '%s' or '%s'", raw, Sqlite, PostgreSql)
	}
}

// 应用数据库配置
type AppDatabase struct {
	Key   string `gorm:"primaryKey;type:varchar(100)"`
	Value string `gorm:"not null"`
}
