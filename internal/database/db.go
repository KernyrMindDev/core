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
	var err error
	var db *gorm.DB
	gormConfig := &gorm.Config{}
	// 调试模式, 开启数据库日志
	if appconfig.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}
	switch appconfig.Database.DBType {
	case model.Sqlite:
		// SQLite开启WAL模式防止并发锁表
		dsn := "data.db?_journal_mode=WAL"
		db, err = gorm.Open(sqlite.Open(dsn), gormConfig)
	case model.PostgreSql:
		// 构造postgresql连接信息
		var sslmode string
		if appconfig.Database.PgSSL {
			sslmode = "enable"
		} else {
			sslmode = "disable"
		}
		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", appconfig.Database.PgHost, appconfig.Database.PgUser, appconfig.Database.PgPassword, appconfig.Database.PgDbname, appconfig.Database.PgPort, sslmode)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	}

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

	log.Println("数据库初始化成功")
	return db
}
