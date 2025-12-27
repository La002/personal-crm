package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Configuration struct {
	DB    DB    `yaml:"db"`
	Log   Log   `yaml:"log"`
	OAuth OAuth `yaml:"oauth"`
	JWT   JWT   `yaml:"jwt"`
}

type DB struct {
	URL      string `yaml:"url" env:"DB_URL"` // Cloud database connection string
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

type OAuth struct {
	GoogleClientID     string `yaml:"google_client_id" mapstructure:"google_client_id" env:"OAUTH_GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `yaml:"google_client_secret" mapstructure:"google_client_secret" env:"OAUTH_GOOGLE_CLIENT_SECRET"`
	RedirectURL        string `yaml:"redirect_url" mapstructure:"redirect_url" env:"OAUTH_REDIRECT_URL"`
}

type JWT struct {
	SecretKey   string `yaml:"secret_key" mapstructure:"secret_key" env:"JWT_SECRET_KEY"`
	ExpiryHours int    `yaml:"expiry_hours" mapstructure:"expiry_hours" env:"JWT_EXPIRY_HOURS"`
}

func NewConfig() *Configuration {
	var config Configuration

	// Enable environment variable support with nested keys
	// This converts "db.url" -> "DB_URL", "oauth.google_client_id" -> "OAUTH_GOOGLE_CLIENT_ID"
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Try to read config file (optional for production)
	viper.SetConfigFile("config/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using environment variables only: %s", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}
	return &config
}
