package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"github.com/vancone/vancone-web-common-go/pkg/logger"
)

var Client *redis.Client
var localConfig Config

type Config struct {
	Url      string
	Password string
}

func Init(viper *viper.Viper) {
	err := viper.UnmarshalKey("redis", &localConfig)
	if err != nil {
		logger.Errorf("viper unmarshal err: %v", err)
		return
	}

	Client = redis.NewClient(&redis.Options{
		Addr:     localConfig.Url,
		Password: localConfig.Password,
		DB:       0,
	})

	_, err = Client.Ping(context.Background()).Result()
	if err != nil {
		logger.Errorf("Failed to connect to Redis: %v", err)
	}
}
