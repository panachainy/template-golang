package logger

import (
	"context"
	"testing"

	pkgContext "template-golang/pkg/context"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   bool
	}{
		{
			name:   "nil config uses default",
			config: nil,
			want:   true,
		},
		{
			name: "custom config",
			config: &Config{
				Level:       "debug",
				Format:      "json",
				Development: true,
				OutputPaths: []string{"stdout"},
			},
			want: true,
		},
		{
			name: "invalid level uses info",
			config: &Config{
				Level:  "invalid",
				Format: "console",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewLogger(tt.config)

			if tt.want {
				assert.NoError(t, err)
				assert.NotNil(t, logger)
				assert.NotNil(t, logger.Logger)
				assert.NotNil(t, logger.Sugar())
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "info", config.Level)
	assert.Equal(t, "console", config.Format)
	assert.False(t, config.Development)
	assert.Equal(t, []string{"stdout"}, config.OutputPaths)
}

func TestGetDefault(t *testing.T) {
	logger1 := GetDefault()
	logger2 := GetDefault()

	// Should return the same instance (singleton)
	assert.Equal(t, logger1, logger2)
	assert.NotNil(t, logger1.Logger)
	assert.NotNil(t, logger1.Sugar())
}

func TestWithContext(t *testing.T) {
	// Create a test logger with observer
	core, recorded := observer.New(zap.InfoLevel)
	logger := &Logger{
		Logger: zap.New(core),
	}
	logger.sugar = logger.Logger.Sugar()

	// Create context with values
	ctx := context.Background()
	ctx = context.WithValue(ctx, pkgContext.RequestIDKey, "test-request-123")
	ctx = context.WithValue(ctx, pkgContext.UserIDKey, "user-456")
	ctx = context.WithValue(ctx, pkgContext.TraceIDKey, "trace-789")

	// Create logger with context
	contextLogger := logger.WithContext(ctx)
	contextLogger.Info("test message")

	// Verify fields were added
	entries := recorded.All()
	assert.Len(t, entries, 1)

	entry := entries[0]
	assert.Equal(t, "test message", entry.Message)

	// Check that context fields were added
	fieldMap := make(map[string]interface{})
	for _, field := range entry.Context {
		fieldMap[field.Key] = field.Interface
	}

	assert.Equal(t, "test-request-123", fieldMap[string(pkgContext.RequestIDKey)])
	assert.Equal(t, "user-456", fieldMap[string(pkgContext.UserIDKey)])
	assert.Equal(t, "trace-789", fieldMap[string(pkgContext.TraceIDKey)])
}

func TestWithFields(t *testing.T) {
	// Create a test logger with observer
	core, recorded := observer.New(zap.InfoLevel)
	logger := &Logger{
		Logger: zap.New(core),
	}
	logger.sugar = logger.Logger.Sugar()

	// Add fields
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	fieldsLogger := logger.WithFields(fields)
	fieldsLogger.Info("test message")

	// Verify fields were added
	entries := recorded.All()
	assert.Len(t, entries, 1)

	entry := entries[0]
	assert.Equal(t, "test message", entry.Message)

	// Check that fields were added
	fieldMap := make(map[string]interface{})
	for _, field := range entry.Context {
		fieldMap[field.Key] = field.Interface
	}

	assert.Equal(t, "value1", fieldMap["key1"])
	assert.Equal(t, 123, fieldMap["key2"])
	assert.Equal(t, true, fieldMap["key3"])
}

func TestWithField(t *testing.T) {
	// Create a test logger with observer
	core, recorded := observer.New(zap.InfoLevel)
	logger := &Logger{
		Logger: zap.New(core),
	}
	logger.sugar = logger.Logger.Sugar()

	// Add single field
	fieldLogger := logger.WithField("test_key", "test_value")
	fieldLogger.Info("test message")

	// Verify field was added
	entries := recorded.All()
	assert.Len(t, entries, 1)

	entry := entries[0]
	assert.Equal(t, "test message", entry.Message)

	// Check that field was added
	fieldMap := make(map[string]interface{})
	for _, field := range entry.Context {
		fieldMap[field.Key] = field.Interface
	}

	assert.Equal(t, "test_value", fieldMap["test_key"])
}

func TestWithError(t *testing.T) {
	// Create a test logger with observer
	core, recorded := observer.New(zap.InfoLevel)
	logger := &Logger{
		Logger: zap.New(core),
	}
	logger.sugar = logger.Logger.Sugar()

	// Test with error
	err := assert.AnError
	errorLogger := logger.WithError(err)
	errorLogger.Info("test message")

	// Verify error was added
	entries := recorded.All()
	assert.Len(t, entries, 1)

	entry := entries[0]
	assert.Equal(t, "test message", entry.Message)

	// Check that error was added
	fieldMap := make(map[string]interface{})
	for _, field := range entry.Context {
		fieldMap[field.Key] = field.Interface
	}

	assert.Equal(t, err.Error(), fieldMap["error"])
}

func TestWithError_NilError(t *testing.T) {
	// Create a test logger with observer
	core, recorded := observer.New(zap.InfoLevel)
	logger := &Logger{
		Logger: zap.New(core),
	}
	logger.sugar = logger.Logger.Sugar()

	// Test with nil error - should return same logger
	errorLogger := logger.WithError(nil)
	assert.Equal(t, logger, errorLogger)
}

func TestGlobalFunctions(t *testing.T) {
	// Test that global functions don't panic
	assert.NotPanics(t, func() {
		Debug("debug message")
		Info("info message")
		Warn("warn message")
		Error("error message")

		Debugf("debug message %s", "formatted")
		Infof("info message %s", "formatted")
		Warnf("warn message %s", "formatted")
		Errorf("error message %s", "formatted")
	})
}
