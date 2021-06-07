package config

import (
	"log"
	"time"

	"github.com/chitchat-awsome/pkg/psqlconnector"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address      string        `mapstructure:"address"`
	ReadTimeout  time.Duration `mapstructure:"readtimeout"`
	WriteTimeout time.Duration `mapstructure:"writetimeout"`
	Static       string        `mapstructure:"static"`
}

type AppConfig struct {
	Server ServerConfig                     `mapstructure:"application"`
	Psql   psqlconnector.PsqlConfigurations `mapstructure:"psql"`
}

var values AppConfig

func init() {
	config := viper.New()
	config.SetConfigName("config") // config file name
	config.AddConfigPath("./config/")
	// config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("error read config / %s", err)
	}

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("Error while reading config: %s", err)
	}

	if err := config.Unmarshal(&values); err != nil {
		log.Fatalf("Error while parsing config: %s", err)
	}
}

func GetConfig() *AppConfig {
	return &values
}
