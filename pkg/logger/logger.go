// Package logger provides structured logging for the application.
package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 定義日誌記錄介面
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sync() error
}

// zapLogger 實作 Logger 介面
type logger struct {
	logger *zap.Logger
}

// Info 記錄資訊級別日誌
func (l *logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Error 記錄錯誤級別日誌
func (l *logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Warn 記錄警告級別日誌
func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Debug 記錄除錯級別日誌
func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Panic 記錄並觸發 panic
func (l *logger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

// Fatal 記錄並終止程式
func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Sync 同步日誌緩衝區
func (l *logger) Sync() error {
	return l.logger.Sync()
}

// NewLogger 建立新的 Logger 實例
func NewLogger() (Logger, error) {
	mode := os.Getenv("GIN_MODE")

	var cfg zap.Config
	if strings.EqualFold(mode, "release") {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	if lvl := strings.TrimSpace(os.Getenv("LOG_LEVEL")); lvl != "" {
		var level zapcore.Level
		if err := level.Set(strings.ToLower(lvl)); err != nil {
			return nil, err
		}
		cfg.Level = zap.NewAtomicLevelAt(level)
	}

	zapLog, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &logger{logger: zapLog}, nil
}
