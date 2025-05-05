package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		// Note: The mapstructure:",squash" tag ensures that nested fields are treated as top-level environment variables.
		Server ServerConfig `mapstructure:",squash"`
		Db     DbConfig     `mapstructure:",squash"`
	}

	ServerConfig struct {
		Port int `mapstructure:"SERVER_PORT"`
	}

	DbConfig struct {
		Host     string `mapstructure:"DB_HOST"`
		Port     int    `mapstructure:"DB_PORT"`
		UserName string `mapstructure:"DB_USERNAME"`
		Password string `mapstructure:"DB_PASSWORD"`
		DBName   string `mapstructure:"DB_DBNAME"`
		SSLMode  string `mapstructure:"DB_SSLMODE"`
		TimeZone string `mapstructure:"DB_TIMEZONE"`
	}
)

var (
	_once   sync.Once
	_config = &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Db: DbConfig{
			Host:     "0.0.0.0",
			Port:     5432,
			UserName: "postgres",
			Password: "postgres",
			DBName:   "postgres",
			SSLMode:  "disable",
			TimeZone: "Asia/Bangkok",
		},
	}
)

func Provide() *Config {
	_once.Do(func() {
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// Automatically override default values with environment variables
		viper.AutomaticEnv()

		// Load .env file
		viper.SetConfigName(".env")
		// Set the configuration file type
		viper.SetConfigType("env")
		viper.AddConfigPath(".")

		// Read the configuration file
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Fatal error loading config file: %s\n", err)
		}

		// Unmarshal the configuration into the Config struct
		if err := viper.Unmarshal(&_config); err != nil {
			panic(fmt.Errorf("unable to decode into struct: %v", err))
		}

		fmt.Println("Config loaded successfully")
	})

	return _config
}
