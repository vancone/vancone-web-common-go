package database

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/vancone/vancone-web-common-go/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var DbConfig Config

type Config struct {
	Url      string
	Username string
	Password string
}

func Init(viper *viper.Viper) {
	err := viper.UnmarshalKey("database", &DbConfig)
	if err != nil {
		logger.Errorf("viper unmarshal err: %v", err)
		return
	}

	dsn := fmt.Sprintf("%s:%s@%s", DbConfig.Username, DbConfig.Password, DbConfig.Url)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("mysql connect error: %v", err)
	}
	if Db.Error != nil {
		logger.Fatalf("database error: %v", Db.Error)
	}

	// database/sql 默认永不回收连接（ConnMaxLifetime/ConnMaxIdleTime=0），
	// 空闲连接被 MySQL wait_timeout 或中间网络设备掐断后，请求从池里拿到
	// 陈旧连接会偶发 "invalid connection" 错误。定期轮换连接避免该问题。
	sqlDb, err := Db.DB()
	if err != nil {
		logger.Fatalf("get sql db error: %v", err)
	}
	sqlDb.SetMaxOpenConns(50)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(30 * time.Minute)
	sqlDb.SetConnMaxIdleTime(5 * time.Minute)
}
