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
)

var (
	_once   sync.Once
	_config = &Config{
		Server: &Server{
			Port: 8083,
		},
	}
)

func GetConfig() *Config {
	_once.Do(func() {
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		// Load .env file
		viper.SetConfigName(".env")
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
		fmt.Printf("Config loaded successfully %+v\n", _config.Db)
		fmt.Printf("Config loaded successfully %+v\n", _config.Server)
	})

	return _config
}
