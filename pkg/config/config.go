package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/vancone/vancone-web-common-go/pkg/encrypt"
)

var App AppConfig
var Mail MailConfig
var Server ServerConfig

type Config struct {
	App    AppConfig
	Mail   MailConfig
	Server ServerConfig
}

type AppConfig struct {
	Name string
}

type MailConfig struct {
	Protocol string
	Host     string
	Port     int
	Username string
	Password string
}

type ServerConfig struct {
	Port int64
}

func ReadConfig() *viper.Viper {
	viperConfig := viper.New()
	viperConfig.AddConfigPath("./conf")
	viperConfig.SetConfigType("yml")
	viperConfig.SetConfigName("config")
	if err := viperConfig.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Fatalln("Config file not found")
		}
	}
	err := processEncryptedKeys(viperConfig)
	if err != nil {
		log.Println("Failed to process encrypted keys", err)
	}

	initAppConfig(viperConfig)
	initMailConfig(viperConfig)
	initServerConfig(viperConfig)
	return viperConfig
}

func initAppConfig(viperConfig *viper.Viper) {
	err := viperConfig.UnmarshalKey("app", &App)
	if err != nil {
		log.Println("viper unmarshal err:", err)
	}
}

func initMailConfig(viperConfig *viper.Viper) {
	err := viperConfig.UnmarshalKey("mail", &Mail)
	if err != nil {
		log.Println("viper unmarshal err:", err)
	}
}

func initServerConfig(viperConfig *viper.Viper) {
	err := viperConfig.UnmarshalKey("server", &Server)
	if err != nil {
		log.Println("viper unmarshal err:", err)
	}
	if Server.Port == 0 {
		Server.Port = 8080
	}
}

func processEncryptedKeys(viperConfig *viper.Viper) error {
	for key, value := range viperConfig.AllSettings() {
		switch val := value.(type) {
		case string:
			if strings.HasPrefix(val, "ENC(") && strings.HasSuffix(val, ")") {
				encrypted := strings.TrimPrefix(strings.TrimSuffix(val, ")"), "ENC(")
				decrypted, err := encrypt.Decrypt(encrypted)
				if err != nil {
					return fmt.Errorf("failed to decrypt %s: %v %s", key, err, decrypted)
				}
				viperConfig.Set(key, decrypted)
			}
		case map[string]interface{}:
			// 处理嵌套的配置项
			subV := viper.New()
			subV.MergeConfigMap(val)
			if err := processEncryptedKeys(subV); err != nil {
				return err
			}
			// 更新处理后的嵌套配置
			viperConfig.Set(key, subV.AllSettings())
		case []interface{}:
			// 处理数组类型的值
			for i, item := range val {
				if strItem, ok := item.(string); ok {
					if strings.HasPrefix(strItem, "ENC(") && strings.HasSuffix(strItem, ")") {
						// 提取加密内容
						encContent := strings.TrimPrefix(strings.TrimSuffix(strItem, ")"), "ENC(")
						// 解密
						decrypted, err := encrypt.Decrypt(encContent)
						if err != nil {
							return fmt.Errorf("failed to decrypt %s[%d]: %v", key, i, err)
						}
						// 更新数组中的值
						val[i] = decrypted
					}
				}
			}
			// 更新处理后的数组
			viperConfig.Set(key, val)
		}
	}
	return nil
}
