package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger with additional functionality
type Logger struct {
	*zap.Logger
	sugar *zap.SugaredLogger
}

// Config holds logger configuration
type Config struct {
	Level       string `mapstructure:"LOG_LEVEL"`       // debug, info, warn, error
	Format      string `mapstructure:"LOG_FORMAT"`      // json, console
	Development bool   `mapstructure:"LOG_DEVELOPMENT"` // development mode
	OutputPaths []string `mapstructure:"LOG_OUTPUT_PATHS"` // output file paths
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:       "info",
		Format:      "console",
		Development: false,
		OutputPaths: []string{"stdout"},
	}
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// NewLogger creates a new logger instance
func NewLogger(config *Config) (*Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Configure encoder
	var encoderConfig zapcore.EncoderConfig
	if config.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Configure encoder format
	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure output
	var outputs []zapcore.WriteSyncer
	if len(config.OutputPaths) == 0 {
		outputs = append(outputs, zapcore.AddSync(os.Stdout))
	} else {
		for _, path := range config.OutputPaths {
			if path == "stdout" {
				outputs = append(outputs, zapcore.AddSync(os.Stdout))
			} else if path == "stderr" {
				outputs = append(outputs, zapcore.AddSync(os.Stderr))
			} else {
				file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					return nil, err
				}
				outputs = append(outputs, zapcore.AddSync(file))
			}
		}
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(outputs...),
		level,
	)

	// Build logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		Logger: zapLogger,
		sugar:  zapLogger.Sugar(),
	}, nil
}

// GetDefault returns the default logger instance (singleton)
func GetDefault() *Logger {
	once.Do(func() {
		logger, err := NewLogger(DefaultConfig())
		if err != nil {
			panic("failed to create default logger: " + err.Error())
		}
		defaultLogger = logger
	})
	return defaultLogger
}

// SetDefault sets the default logger
func SetDefault(logger *Logger) {
	defaultLogger = logger
}

// Sugar returns the sugared logger for formatted logging
func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.sugar
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	logger := l.Logger

	// Add request ID if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		logger = logger.With(zap.Any("request_id", requestID))
	}

	// Add user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		logger = logger.With(zap.Any("user_id", userID))
	}

	// Add trace ID if available
	if traceID := ctx.Value("trace_id"); traceID != nil {
		logger = logger.With(zap.Any("trace_id", traceID))
	}

	return &Logger{
		Logger: logger,
		sugar:  logger.Sugar(),
	}
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	var zapFields []zap.Field
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}

	logger := l.Logger.With(zapFields...)
	return &Logger{
		Logger: logger,
		sugar:  logger.Sugar(),
	}
}

// WithField adds a single field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	logger := l.Logger.With(zap.Any(key, value))
	return &Logger{
		Logger: logger,
		sugar:  logger.Sugar(),
	}
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}
	logger := l.Logger.With(zap.Error(err))
	return &Logger{
		Logger: logger,
		sugar:  logger.Sugar(),
	}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Global convenience functions using default logger
func Debug(msg string, fields ...zap.Field) {
	GetDefault().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	GetDefault().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetDefault().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetDefault().Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetDefault().Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	GetDefault().Panic(msg, fields...)
}

// Formatted logging functions
func Debugf(template string, args ...interface{}) {
	GetDefault().Sugar().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	GetDefault().Sugar().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	GetDefault().Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	GetDefault().Sugar().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	GetDefault().Sugar().Fatalf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	GetDefault().Sugar().Panicf(template, args...)
}
