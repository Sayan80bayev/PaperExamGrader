package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL   string `mapstructure:"DATABASE_URL"`
	Port          string `mapstructure:"PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	MinioBucket   string `mapstructure:"MINIO_BUCKET"`
	MinioEndpoint string `mapstructure:"MINIO_ENDPOINT"`
	AccessKey     string `mapstructure:"ACCESS_KEY"`
	SecretKey     string `mapstructure:"SECRET_KEY"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile("config/config.yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Couldn't load config.yaml: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
