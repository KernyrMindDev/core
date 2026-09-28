package database

import (
	"fmt"
	"log"

	"github.com/KernyrMindDev/core/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(appconfig *model.AppConfig) *gorm.DB {
	// 初始化数据库连接
	db, err := InitDatabaseConnect(&appconfig.Database, appconfig.Debug)

	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移/建表
	err = db.AutoMigrate(
		&model.User{},
		&model.Board{},
		&model.BoardNode{},
		&model.BoardConnection{},
		&model.AppDatabase{},
	)
	if err != nil {
		log.Fatalf("数据表自动迁移失败: %v", err)
	}

	// log.Println("数据库初始化成功")
	return db
}

// 初始化数据库连接
func InitDatabaseConnect(dbconfig *model.Database, debug bool) (*gorm.DB, error) {
	gormConfig := &gorm.Config{}
	// 调试模式, 开启数据库日志
	if debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	switch dbconfig.DBType {
	case model.Sqlite:
		// SQLite开启WAL模式防止并发锁表
		dsn := "data.db?_journal_mode=WAL"
		return gorm.Open(sqlite.Open(dsn), gormConfig)
	case model.PostgreSql:
		// 构造postgresql连接信息
		var sslmode string
		if dbconfig.PgSSL {
			sslmode = "enable"
		} else {
			sslmode = "disable"
		}
		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", dbconfig.PgHost, dbconfig.PgUser, dbconfig.PgPassword, dbconfig.PgDbname, dbconfig.PgPort, sslmode)
		return gorm.Open(postgres.Open(dsn), gormConfig)
	}
	return nil, fmt.Errorf("Unknown database type")
}
