package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(logLevel string) *zap.Logger {
	logTemplate := "{\"level\":\"%s\",\"timestamp\":\"%s\",\"function\":\"github.com/nighostchris/everytrack-backend/internal/logger.New\",\"message\":\"%s\"}\n"
	config := zap.NewProductionConfig()
	fmt.Printf(logTemplate, "info", time.Now().Format(time.RFC3339Nano), "initializing logger")

	// Set the log level
	var level zapcore.Level
	switch logLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.ErrorLevel
	}
	config.Level = zap.NewAtomicLevelAt(level)

	// Set log format
	config.Encoding = "json"

	// Fine tune timestamp, function and message fields in log
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoderConfig.CallerKey = ""
	encoderConfig.MessageKey = "message"
	encoderConfig.FunctionKey = "function"
	config.EncoderConfig = encoderConfig

	// Create logger instance
	zapLogger, buildLoggerError := config.Build()
	if buildLoggerError != nil {
		fmt.Printf(logTemplate, "error", time.Now().Format(time.RFC3339Nano), buildLoggerError.Error())
		os.Exit(1)
	}
	defer zapLogger.Sync()
	fmt.Printf(logTemplate, "info", time.Now().Format(time.RFC3339Nano), "initialized logger")

	return zapLogger
}
