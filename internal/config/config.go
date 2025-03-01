package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name    string `yaml:"name"`
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Prefork bool   `yaml:"prefork"`
		Debug   bool   `yaml:"debug"`
	} `yaml:"app"`
	Auth struct {
		AccessTokenExpMins  int    `yaml:"access_token_exp_mins"`
		RefreshTokenExpDays int    `yaml:"refresh_token_exp_days"`
		JWTSecret           string `yaml:"jwt_secrets"`
		JWTSalt             string `yaml:"jwt_salt"`
	} `yaml:"auth"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

func NewConfig() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("yaml")
	config.AddConfigPath("./../")
	config.AddConfigPath("./")

	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return config
}
