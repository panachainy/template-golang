package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type (
	Config struct {
		// Note: The mapstructure:",squash" tag ensures that nested fields are treated as top-level environment variables.
		Server ServerConfig `mapstructure:",squash"`
		Db     DbConfig     `mapstructure:",squash"`
		Auth   AuthConfig   `mapstructure:",squash"`
	}

	ServerConfig struct {
		Port int    `mapstructure:"SERVER_PORT"`
		Mode string `mapstructure:"GIN_MODE"`
	}

	DbConfig struct {
		Host          string `mapstructure:"DB_HOST"`
		Port          int    `mapstructure:"DB_PORT"`
		UserName      string `mapstructure:"DB_USERNAME"`
		Password      string `mapstructure:"DB_PASSWORD"`
		DBName        string `mapstructure:"DB_DBNAME"`
		SSLMode       string `mapstructure:"DB_SSLMODE"`
		TimeZone      string `mapstructure:"DB_TIMEZONE"`
		MigrationPath string `mapstructure:"DB_MIGRATION_PATH"`
	}

	AuthConfig struct {
		PrivateKeyPath string `mapstructure:"PRIVATE_KEY_PATH"`

		LineClientID      string `mapstructure:"LINE_CLIENT_ID"`
		LineClientSecret  string `mapstructure:"LINE_CLIENT_SECRET"`
		LineCallbackURL   string `mapstructure:"LINE_CALLBACK_URL"`
		LineFECallbackURL string `mapstructure:"LINE_FE_CALLBACK_URL"`
	}
)

type ConfigOption struct {
	ConfigPath string
}

var (
	_once   sync.Once
	_config = &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Db: DbConfig{
			Host:          "0.0.0.0",
			Port:          5432,
			UserName:      "postgres",
			Password:      "postgres",
			DBName:        "postgres",
			SSLMode:       "disable",
			TimeZone:      "Asia/Bangkok",
			MigrationPath: "file://db/migrations",
		},
		Auth: AuthConfig{
			PrivateKeyPath: "private.pem",
		},
	}
)

func NewConfigOption(configPath string) *ConfigOption {
	return &ConfigOption{
		ConfigPath: configPath,
	}
}

func Provide(configOption *ConfigOption) *Config {
	_once.Do(func() {
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

		// Automatically override default values with environment variables
		viper.AutomaticEnv()

		fmt.Println("===================== Load .env =============================")

		// Determine which .env file to use
		var targetEnvFile string
		if gin.Mode() == gin.TestMode {
			fmt.Println("Test mode detected, loading .env.test file")
			targetEnvFile = ".env.test"
		} else {
			fmt.Println("loading .env file")
			targetEnvFile = ".env"
		}

		// Load .env file
		viper.SetConfigName(targetEnvFile)
		viper.SetConfigType("env")
		viper.AddConfigPath(configOption.ConfigPath)

		// Read the configuration file
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Fatal error loading config file: %s\n", err)
		}

		fmt.Println("======================================================")

		// Unmarshal the configuration into the Config struct
		if err := viper.Unmarshal(&_config); err != nil {
			panic(fmt.Errorf("unable to decode into struct: %v", err))
		}

		fmt.Println("Config loaded successfully")

		if _config.Server.Mode != "release" {
			fmt.Println("======================================================")
			fmt.Printf("[Loaded] Config: %+v\n", _config)
			fmt.Println("======================================================")
		}
	})

	return _config
}
