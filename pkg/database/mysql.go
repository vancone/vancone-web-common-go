package database

import (
	"fmt"

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
}
