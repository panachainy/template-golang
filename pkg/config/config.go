package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Loader interface for configuration loading
type Loader interface {
	Load() error
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
	Unmarshal(rawVal interface{}) error
	UnmarshalKey(key string, rawVal interface{}) error
	IsSet(key string) bool
	AllKeys() []string
	AllSettings() map[string]interface{}
}

// Config holds configuration loading options
type Config struct {
	ConfigName      string            // config file name (without extension)
	ConfigType      string            // config file type (yaml, json, toml, etc.)
	ConfigPaths     []string          // paths to search for config file
	EnvPrefix       string            // prefix for environment variables
	EnvKeyReplacer  *strings.Replacer // replacer for environment keys
	Defaults        map[string]interface{} // default values
	AutomaticEnv    bool              // automatically bind environment variables
	AllowEmptyEnv   bool              // allow empty environment variables
	TestConfigName  string            // config file name for test environment
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		ConfigName:     ".env",
		ConfigType:     "env",
		ConfigPaths:    []string{".", "./config"},
		EnvPrefix:      "",
		EnvKeyReplacer: strings.NewReplacer(".", "_", "-", "_"),
		Defaults:       make(map[string]interface{}),
		AutomaticEnv:   true,
		AllowEmptyEnv:  false,
		TestConfigName: ".env.test",
	}
}

// ViperLoader implements Loader using Viper
type ViperLoader struct {
	v      *viper.Viper
	config *Config
	once   sync.Once
}

// NewViperLoader creates a new Viper-based configuration loader
func NewViperLoader(config *Config) *ViperLoader {
	if config == nil {
		config = DefaultConfig()
	}

	v := viper.New()

	return &ViperLoader{
		v:      v,
		config: config,
	}
}

// Load loads the configuration
func (l *ViperLoader) Load() error {
	var loadErr error

	l.once.Do(func() {
		// Set up environment variable handling
		if l.config.EnvPrefix != "" {
			l.v.SetEnvPrefix(l.config.EnvPrefix)
		}

		if l.config.EnvKeyReplacer != nil {
			l.v.SetEnvKeyReplacer(l.config.EnvKeyReplacer)
		}

		if l.config.AutomaticEnv {
			l.v.AutomaticEnv()
		}

		// Set defaults
		for key, value := range l.config.Defaults {
			l.v.SetDefault(key, value)
		}

		// Determine config file name based on environment
		configName := l.config.ConfigName
		if isTestEnvironment() && l.config.TestConfigName != "" {
			configName = l.config.TestConfigName
		}

		// Set up config file
		l.v.SetConfigName(configName)
		l.v.SetConfigType(l.config.ConfigType)

		// Add config paths
		for _, path := range l.config.ConfigPaths {
			l.v.AddConfigPath(path)
		}

		// Try to read config file
		if err := l.v.ReadInConfig(); err != nil {
			// Only return error if it's not a "file not found" error
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				loadErr = fmt.Errorf("failed to read config file: %w", err)
				return
			}
		}
	})

	return loadErr
}

// Get retrieves a value by key
func (l *ViperLoader) Get(key string) interface{} {
	return l.v.Get(key)
}

// GetString retrieves a string value by key
func (l *ViperLoader) GetString(key string) string {
	return l.v.GetString(key)
}

// GetInt retrieves an integer value by key
func (l *ViperLoader) GetInt(key string) int {
	return l.v.GetInt(key)
}

// GetBool retrieves a boolean value by key
func (l *ViperLoader) GetBool(key string) bool {
	return l.v.GetBool(key)
}

// GetFloat64 retrieves a float64 value by key
func (l *ViperLoader) GetFloat64(key string) float64 {
	return l.v.GetFloat64(key)
}

// GetDuration retrieves a duration value by key
func (l *ViperLoader) GetDuration(key string) time.Duration {
	return l.v.GetDuration(key)
}

// GetStringSlice retrieves a string slice value by key
func (l *ViperLoader) GetStringSlice(key string) []string {
	return l.v.GetStringSlice(key)
}

// Unmarshal unmarshals the config into a struct
func (l *ViperLoader) Unmarshal(rawVal interface{}) error {
	return l.v.Unmarshal(rawVal)
}

// UnmarshalKey unmarshals a specific key into a struct
func (l *ViperLoader) UnmarshalKey(key string, rawVal interface{}) error {
	return l.v.UnmarshalKey(key, rawVal)
}

// IsSet checks if a key is set
func (l *ViperLoader) IsSet(key string) bool {
	return l.v.IsSet(key)
}

// AllKeys returns all keys
func (l *ViperLoader) AllKeys() []string {
	return l.v.AllKeys()
}

// AllSettings returns all settings
func (l *ViperLoader) AllSettings() map[string]interface{} {
	return l.v.AllSettings()
}

// Helper functions

// isTestEnvironment checks if we're in a test environment
func isTestEnvironment() bool {
	return strings.Contains(strings.ToLower(os.Args[0]), "test") ||
		   os.Getenv("GO_ENV") == "test" ||
		   os.Getenv("GIN_MODE") == "test"
}

// LoadWithDefaults loads configuration with default settings
func LoadWithDefaults(defaults map[string]interface{}) (Loader, error) {
	config := DefaultConfig()
	config.Defaults = defaults

	loader := NewViperLoader(config)
	err := loader.Load()
	return loader, err
}

// LoadFromPath loads configuration from a specific path
func LoadFromPath(configPath, configName, configType string) (Loader, error) {
	config := DefaultConfig()
	config.ConfigPaths = []string{configPath}
	config.ConfigName = configName
	config.ConfigType = configType

	loader := NewViperLoader(config)
	err := loader.Load()
	return loader, err
}

// LoadYAMLConfig loads a YAML configuration file
func LoadYAMLConfig(configPath, configName string) (Loader, error) {
	return LoadFromPath(configPath, configName, "yaml")
}

// LoadJSONConfig loads a JSON configuration file
func LoadJSONConfig(configPath, configName string) (Loader, error) {
	return LoadFromPath(configPath, configName, "json")
}

// LoadEnvConfig loads environment configuration (.env file)
func LoadEnvConfig(configPath string) (Loader, error) {
	config := DefaultConfig()
	if configPath != "" {
		config.ConfigPaths = []string{configPath}
	}

	loader := NewViperLoader(config)
	err := loader.Load()
	return loader, err
}

// MustLoad loads configuration and panics on error
func MustLoad(config *Config) Loader {
	loader := NewViperLoader(config)
	if err := loader.Load(); err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	return loader
}

// LoadAndUnmarshal loads configuration and unmarshals it into a struct
func LoadAndUnmarshal(config *Config, target interface{}) error {
	loader := NewViperLoader(config)
	if err := loader.Load(); err != nil {
		return err
	}
	return loader.Unmarshal(target)
}

// MustLoadAndUnmarshal loads configuration, unmarshals it, and panics on error
func MustLoadAndUnmarshal(config *Config, target interface{}) {
	if err := LoadAndUnmarshal(config, target); err != nil {
		panic(fmt.Sprintf("failed to load and unmarshal configuration: %v", err))
	}
}
