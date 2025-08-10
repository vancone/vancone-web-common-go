package database

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var DbConfig Config

type Config struct {
	Url      string
	User     string
	Password string
}

func Init(viper *viper.Viper) {
	err := viper.UnmarshalKey("database", &DbConfig)
	if err != nil {
		log.Println("viper unmarshal err:", err)
		return
	}

	dsn := fmt.Sprintf("%s:%s@%s", DbConfig.User, DbConfig.Password, DbConfig.Url)
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("mysql connect error", err)
	}
	if Db.Error != nil {
		log.Fatalln("database error", Db.Error)
	}
}
