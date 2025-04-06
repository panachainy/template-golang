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
			Port: 8083,
		},
		Db: DbConfig{
			Host: "localhost",
		},
	}
)

func GetConfig() *Config {
	_once.Do(func() {
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// Load .env file
		viper.SetConfigName(".env")
		// Set the configuration file type
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		if err := viper.MergeInConfig(); err != nil {
			fmt.Printf("Error loading .env file: %v\n", err)
		}

		// Load environment variables
		viper.AutomaticEnv()

		if err := viper.Unmarshal(&_config); err != nil {
			panic(err)
		}

		fmt.Println("=================================")

		for _, key := range viper.AllKeys() {
			fmt.Printf("Key: %s, Value: %v\n", key, viper.Get(key))
		}
		fmt.Println("=================================")
		fmt.Printf("1. Config loaded successfully %+v\n", _config.Db)
		fmt.Printf("2. Config loaded successfully %+v\n", _config.Server)
		fmt.Printf("3. Config loaded successfully %+v\n", _config)
	})

	return _config
}
