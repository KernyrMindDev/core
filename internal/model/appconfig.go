package model

import "fmt"

// 应用配置结构体
type AppConfig struct {
	Port     int      `yaml:"port"`
	Database Database `yaml:"database"`
	Debug    bool     `yaml:"debug"`
}

type Database struct {
	DBType DatabaseType `yaml:"type"`
	// PostgreSQL独占配置
	PgHost     string `yaml:"host"`     // 主机名
	PgPort     int    `yaml:"port"`     //端口
	PgUser     string `yaml:"username"` // 用户名
	PgPassword string `yaml:"password"` // 密码
	PgDbname   string `yaml:"dbname"`   // 数据库名
	PgSSL      bool   `yaml:"ssl"`      // 是否开启加密
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

	// 强制转化数据库类型
	val := DatabaseType(raw)
	switch val {
	case Sqlite, PostgreSql:
		// 正确的数据库类型
		*dt = val
		return nil
	default:
		// 数据库类型不属于已知类型
		return fmt.Errorf("invalid database type '%s', expected '%s' or '%s'", raw, Sqlite, PostgreSql)
	}
}

// 应用数据库配置
type AppDatabase struct {
	Key   string `gorm:"primaryKey;type:varchar(100)"`
	Value string `gorm:"not null"`
}
