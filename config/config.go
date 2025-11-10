package config

import (
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	DB  `yaml:"db"`
	Log `yaml:"logger"`
}

type DB struct {
	PoolMax  int64  `yaml:"pool_max" env:"DB_POOl_MAX"`
	Host     string `yaml:"host" env:"DB_HOST"`
	User     string `yaml:"user" env:"DB_USER"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME"`
	Port     string `yaml:"port" env:"DB_PORT"`
}

type Log struct {
	Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
}

func NewConfig() *Configuration {
	var config Configuration
	viper.SetConfigFile("config/config.yml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return &config
}
