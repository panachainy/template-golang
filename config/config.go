package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server *Server
		Db     *Db
	}

	Server struct {
		Port int
	}

	Db struct {
		Host     string
		Port     int
		UserName string
		Password string
		DBName   string
		SSLMode  string
		TimeZone string
	}

	// Server struct {
	// 	Port int `mapstructure:"SERVER_PORT"`
	// }

	// Db struct {
	// 	Host     string `mapstructure:"DB_HOST"`
	// 	Port     int    `mapstructure:"DB_PORT"`
	// 	UserName string `mapstructure:"DB_USERNAME"`
	// 	Password string `mapstructure:"DB_PASSWORD"`
	// 	DBName   string `mapstructure:"DB_DBNAME"`
	// 	SSLMode  string `mapstructure:"DB_SSLMODE"`
	// 	TimeZone string `mapstructure:"DB_TIMEZONE"`
	// }
)

var (
	once           sync.Once
	configInstance *Config
)

func GetConfig() *Config {
	once.Do(func() {
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// Load config.yaml
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if err := viper.MergeInConfig(); err != nil {
			panic(err)
		}

		// Load .env file
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(".")
		if err := viper.MergeInConfig(); err != nil {
			fmt.Printf("Error loading .env file: %v\n", err)
		}

		// Load environment variables
		viper.AutomaticEnv()

		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}

		fmt.Printf("Config loaded successfully %+v\n", viper.AllKeys())
		for _, key := range viper.AllKeys() {
			fmt.Printf("Key: %s, Value: %v\n", key, viper.Get(key))
		}
	})

	return configInstance
}
