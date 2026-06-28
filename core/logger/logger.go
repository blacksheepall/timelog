package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/requestid"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instance Logger
	once     sync.Once
)

type Logger interface {
	Debug(fields ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(fields ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(fields ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(fields ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatal(fields ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	// WithContext returns a logger that embeds the request/trace ID from the
	// context (if any). Use this in request handlers so downstream logs can be
	// correlated with the upstream request ID.
	WithContext(ctx context.Context) Logger
}

var ilogger Logger

func GetLogger() Logger {
	return ilogger
}

func SetZapLogger(cfg config.Config) Logger {
	// Setting log level
	level, err := zap.ParseAtomicLevel(cfg.Log.Level)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse log level: %s", err))
	}

	// Setting log output
	logDir := filepath.Dir(cfg.Log.Path)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, 0744)
		if err != nil {
			panic(fmt.Sprintf("Failed to create log directory: %s", err))
		}
	}

	// Setting log format
	encoderConfig := zapcore.EncoderConfig{
		LevelKey: "lvl", TimeKey: "ts",
		MessageKey: "msg", CallerKey: "src",
		NameKey: "logger", StacktraceKey: "stack", // not used for now
		//
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Setting log rotation
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Log.Path,
		MaxSize:    cfg.Log.Rotation.MaxSize, //MB
		MaxBackups: cfg.Log.Rotation.MaxBackups,
		MaxAge:     cfg.Log.Rotation.MaxAge, // days
	})

	c := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		level,
	)

	logger := zap.New(c,
		zap.AddCaller(),
	).
		Sugar()
	defer logger.Sync()

	ilogger = &zapLogger{logger}

	return ilogger
}

// zapLogger wraps zap.SugaredLogger and implements the Logger interface,
// including context-aware request ID injection.
type zapLogger struct {
	*zap.SugaredLogger
}

func (z *zapLogger) WithContext(ctx context.Context) Logger {
	if reqID := requestid.FromContext(ctx); reqID != "" {
		return &zapLogger{z.SugaredLogger.With("request_id", reqID)}
	}
	return z
}
