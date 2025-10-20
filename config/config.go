package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port                    string `mapstructure:"PORT"`
	DatabaseURL             string `mapstructure:"DATABASE_URL"`
	DatabaseURLTest         string `mapstructure:"DATABASE_URL_TEST`
	JWTSecret               string `mapstructure:"JWT_SECRET_KEY"`
	OtpSendingEmail         string `mapstructure:"OTP_SENDING_EMAIL"`
	OtpSendingEmailPassword string `mapstructure:"OTP_SENDING_EMAIL_PASSWORD"`
	OAuthClientID           string `mapstructure:"OAUTH_CLIENT_ID"`
	OAuthClientSecret       string `mapstructure:"OAUTH_CLIENT_SECRET"`
	OAuthRedirectURL        string `mapstructure:"OAUTH_REDIRECT_URL"`
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
