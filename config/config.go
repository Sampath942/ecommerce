package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port            string `mapstructure:"PORT"`
	DatabaseURL     string `mapstructure:"DATABASE_URL"`
	TestDatabaseURL string `mapstructure:"TEST_DATABASE_URL`
}

var AppConfig *Config

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	var cfg Config
	if err = viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to unmarshall config, %s", err)
	}

	AppConfig = &cfg
}
