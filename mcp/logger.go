package main

import (
	"os"
	"path/filepath"

	"github.com/blacksheepaul/timelog/core/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var mcpLogger *zap.SugaredLogger

// InitMCPLogger creates a file-only logger for MCP debugging.
// This logger will never write to stdout to avoid breaking MCP protocol.
func InitMCPLogger(cfg *config.Config) *zap.SugaredLogger {
	if cfg == nil || !cfg.MCP.Enabled {
		mcpLogger = zap.NewNop().Sugar()
		return mcpLogger
	}

	if os.Getenv("MCP_DEBUG") == "false" {
		mcpLogger = zap.NewNop().Sugar()
		return mcpLogger
	}

	level, err := zap.ParseAtomicLevel(cfg.MCP.Level)
	if err != nil {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logPath := cfg.MCP.Path
	if logPath == "" {
		logPath = "logs/mcp.log"
	}

	logDir := filepath.Dir(logPath)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			mcpLogger = zap.NewNop().Sugar()
			return mcpLogger
		}
	}

	encoderConfig := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "timestamp",
		MessageKey:     "message",
		CallerKey:      "caller",
		NameKey:        "logger",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    cfg.Log.Rotation.MaxSize,
		MaxBackups: cfg.Log.Rotation.MaxBackups,
		MaxAge:     cfg.Log.Rotation.MaxAge,
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriter,
		level,
	)

	mcpLogger = zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	).Sugar()

	return mcpLogger
}

// LogMCPError logs MCP-specific errors for debugging.
func LogMCPError(operation string, err error, context map[string]interface{}) {
	if mcpLogger == nil {
		return
	}

	fields := []interface{}{
		"operation", operation,
		"error", err.Error(),
	}
	for k, v := range context {
		fields = append(fields, k, v)
	}

	mcpLogger.Errorw("MCP error occurred", fields...)
}

// LogMCPDebug logs debug information for troubleshooting.
func LogMCPDebug(message string, fields map[string]interface{}) {
	if mcpLogger == nil {
		return
	}

	logFields := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		logFields = append(logFields, k, v)
	}

	mcpLogger.Debugw(message, logFields...)
}
