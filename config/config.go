package config

import (
	"fmt"
	"reflect"
	"sync"

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
	// TODO: impl config here
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

func NewConfig(configOption *ConfigOption) *Config {
	_once.Do(func() {
		// Automatically override default values with environment variables
		viper.AutomaticEnv()

		fmt.Println("======================================================")

		// Bind every leaf key in Config to env
		BindEnvsFromStruct("", _config)

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Warning: unable to read config file: %v\n", err)
		}

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

// BindEnvsFromStruct binds environment variables for all fields in the given struct using viper.
// It recursively traverses nested structs and binds each field's mapstructure tag as the env key.
func BindEnvsFromStruct(prefix string, s any) {
	// Use reflection to traverse struct fields
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag == "" || tag == ",squash" {
			// If squash, recurse into nested struct
			BindEnvsFromStruct(prefix, val.Field(i).Interface())
			continue
		}
		// Compose env key
		envKey := tag
		if prefix != "" {
			envKey = prefix + "_" + envKey
		}
		viper.BindEnv(envKey)
	}
}
