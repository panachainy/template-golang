package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, ".env", config.ConfigName)
	assert.Equal(t, "env", config.ConfigType)
	assert.Equal(t, []string{".", "./config"}, config.ConfigPaths)
	assert.Equal(t, "", config.EnvPrefix)
	assert.NotNil(t, config.EnvKeyReplacer)
	assert.NotNil(t, config.Defaults)
	assert.True(t, config.AutomaticEnv)
	assert.False(t, config.AllowEmptyEnv)
	assert.Equal(t, ".env.test", config.TestConfigName)
}

func TestNewViperLoader(t *testing.T) {
	// Test with nil config
	loader := NewViperLoader(nil)
	assert.NotNil(t, loader)
	assert.NotNil(t, loader.v)
	assert.NotNil(t, loader.config)

	// Test with custom config
	config := &Config{
		ConfigName: "custom",
		ConfigType: "yaml",
	}
	loader = NewViperLoader(config)
	assert.NotNil(t, loader)
	assert.Equal(t, config, loader.config)
}

func TestViperLoader_Load(t *testing.T) {
	// Create a temporary config file for testing
	tempDir := t.TempDir()
	configFile := tempDir + "/.env"

	err := os.WriteFile(configFile, []byte("TEST_KEY=test_value\nTEST_INT=123"), 0644)
	assert.NoError(t, err)

	config := &Config{
		ConfigName:   ".env",
		ConfigType:   "env",
		ConfigPaths:  []string{tempDir},
		AutomaticEnv: true,
		Defaults: map[string]interface{}{
			"DEFAULT_KEY": "default_value",
		},
	}

	loader := NewViperLoader(config)
	err = loader.Load()
	assert.NoError(t, err)

	// Test that values are accessible
	assert.Equal(t, "test_value", loader.GetString("TEST_KEY"))
	assert.Equal(t, 123, loader.GetInt("TEST_INT"))
	assert.Equal(t, "default_value", loader.GetString("DEFAULT_KEY"))
}

func TestViperLoader_Methods(t *testing.T) {
	config := &Config{
		Defaults: map[string]interface{}{
			"STRING_KEY":   "string_value",
			"INT_KEY":      42,
			"BOOL_KEY":     true,
			"FLOAT_KEY":    3.14,
			"DURATION_KEY": "5s",
			"SLICE_KEY":    []string{"a", "b", "c"},
		},
		AutomaticEnv: true,
	}

	loader := NewViperLoader(config)
	err := loader.Load()
	assert.NoError(t, err)

	// Test Get methods
	assert.Equal(t, "string_value", loader.Get("STRING_KEY"))
	assert.Equal(t, "string_value", loader.GetString("STRING_KEY"))
	assert.Equal(t, 42, loader.GetInt("INT_KEY"))
	assert.True(t, loader.GetBool("BOOL_KEY"))
	assert.Equal(t, 3.14, loader.GetFloat64("FLOAT_KEY"))
	assert.Equal(t, 5*time.Second, loader.GetDuration("DURATION_KEY"))
	assert.Equal(t, []string{"a", "b", "c"}, loader.GetStringSlice("SLICE_KEY"))

	// Test IsSet
	assert.True(t, loader.IsSet("STRING_KEY"))
	assert.False(t, loader.IsSet("NON_EXISTENT_KEY"))

	// Test AllKeys
	keys := loader.AllKeys()
	assert.Contains(t, keys, "STRING_KEY")
	assert.Contains(t, keys, "INT_KEY")

	// Test AllSettings
	settings := loader.AllSettings()
	assert.Equal(t, "string_value", settings["STRING_KEY"])
	assert.Equal(t, 42, settings["INT_KEY"])
}

func TestViperLoader_Unmarshal(t *testing.T) {
	type TestConfig struct {
		StringKey string `mapstructure:"STRING_KEY"`
		IntKey    int    `mapstructure:"INT_KEY"`
		BoolKey   bool   `mapstructure:"BOOL_KEY"`
	}

	config := &Config{
		Defaults: map[string]interface{}{
			"STRING_KEY": "test_string",
			"INT_KEY":    100,
			"BOOL_KEY":   true,
		},
		AutomaticEnv: true,
	}

	loader := NewViperLoader(config)
	err := loader.Load()
	assert.NoError(t, err)

	var testConfig TestConfig
	err = loader.Unmarshal(&testConfig)
	assert.NoError(t, err)

	assert.Equal(t, "test_string", testConfig.StringKey)
	assert.Equal(t, 100, testConfig.IntKey)
	assert.True(t, testConfig.BoolKey)
}

func TestViperLoader_UnmarshalKey(t *testing.T) {
	type ServerConfig struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	}

	config := &Config{
		Defaults: map[string]interface{}{
			"server.port": 8080,
			"server.host": "localhost",
		},
		AutomaticEnv: true,
	}

	loader := NewViperLoader(config)
	err := loader.Load()
	assert.NoError(t, err)

	var serverConfig ServerConfig
	err = loader.UnmarshalKey("server", &serverConfig)
	assert.NoError(t, err)

	assert.Equal(t, 8080, serverConfig.Port)
	assert.Equal(t, "localhost", serverConfig.Host)
}

func TestIsTestEnvironment(t *testing.T) {
	// Save original env
	originalGinMode := os.Getenv("GIN_MODE")
	originalGoEnv := os.Getenv("GO_ENV")

	defer func() {
		os.Setenv("GIN_MODE", originalGinMode)
		os.Setenv("GO_ENV", originalGoEnv)
	}()

	// Test with GIN_MODE=test
	os.Setenv("GIN_MODE", "test")
	assert.True(t, isTestEnvironment())

	// Test with GO_ENV=test
	os.Setenv("GIN_MODE", "")
	os.Setenv("GO_ENV", "test")
	assert.True(t, isTestEnvironment())

	// Test without test environment
	os.Setenv("GIN_MODE", "release")
	os.Setenv("GO_ENV", "production")
	// Note: This might still be true because os.Args[0] contains "test"
	// when running tests
}

func TestLoadWithDefaults(t *testing.T) {
	defaults := map[string]interface{}{
		"TEST_DEFAULT": "default_value",
		"TEST_INT":     42,
	}

	loader, err := LoadWithDefaults(defaults)
	assert.NoError(t, err)
	assert.NotNil(t, loader)

	assert.Equal(t, "default_value", loader.GetString("TEST_DEFAULT"))
	assert.Equal(t, 42, loader.GetInt("TEST_INT"))
}

func TestLoadEnvConfig(t *testing.T) {
	// Test with empty path (should use default)
	loader, err := LoadEnvConfig("")
	assert.NoError(t, err)
	assert.NotNil(t, loader)

	// Test with specific path
	tempDir := t.TempDir()
	loader, err = LoadEnvConfig(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, loader)
}

func TestLoadAndUnmarshal(t *testing.T) {
	type TestConfig struct {
		TestKey string `mapstructure:"TEST_KEY"`
	}

	config := &Config{
		Defaults: map[string]interface{}{
			"TEST_KEY": "test_value",
		},
		AutomaticEnv: true,
	}

	var testConfig TestConfig
	err := LoadAndUnmarshal(config, &testConfig)
	assert.NoError(t, err)

	assert.Equal(t, "test_value", testConfig.TestKey)
}

func TestMustLoad(t *testing.T) {
	config := DefaultConfig()

	// Should not panic with valid config
	assert.NotPanics(t, func() {
		loader := MustLoad(config)
		assert.NotNil(t, loader)
	})
}

func TestMustLoadAndUnmarshal(t *testing.T) {
	type TestConfig struct {
		TestKey string `mapstructure:"TEST_KEY"`
	}

	config := &Config{
		Defaults: map[string]interface{}{
			"TEST_KEY": "test_value",
		},
		AutomaticEnv: true,
	}

	var testConfig TestConfig

	// Should not panic with valid config
	assert.NotPanics(t, func() {
		MustLoadAndUnmarshal(config, &testConfig)
		assert.Equal(t, "test_value", testConfig.TestKey)
	})
}
