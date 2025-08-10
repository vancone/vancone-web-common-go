package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Client *redis.Client
var localConfig Config

type Config struct {
	Addr     string
	Password string
}

func Init(viper *viper.Viper) {
	err := viper.UnmarshalKey("redis", &localConfig)
	if err != nil {
		log.Println("viper unmarshal err:", err)
		return
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     localConfig.Addr,
		Password: localConfig.Password,
		DB:       0,
	})

	_, err = Client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis", err)
	}
}
