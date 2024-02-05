package config

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/spf13/viper"
)

var Env *envConfigs

func InitEnvConfigs() {
	Env = loadEnvVariables()
}

type envConfigs struct {
	ServerPort   string `mapstructure:"SERVER_PORT"`
	AuthUserName string `mapstructure:"AUTH_USERNAME"`
	AuthPassword string `mapstructure:"AUTH_PASSWORD"`
}

func loadEnvVariables() (config *envConfigs) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}
