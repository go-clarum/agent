package logging

import (
	"fmt"
	"github.com/go-clarum/agent/infrastructure/config"
	"github.com/go-clarum/agent/infrastructure/sync"
	"log"
	"log/slog"
	"os"
	"strings"
)

var activeLogLevel slog.Level
var internalLogger *log.Logger
var defaultLogger *Logger
var LogEmitter *sync.Emitter

func init() {
	activeLogLevel = parseLevel(config.LoggingLevel())
	internalLogger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
	defaultLogger = NewLogger("")
	LogEmitter = sync.NewEmitter()
}

type Logger struct {
	name string
}

func NewLogger(name string) *Logger {
	return &Logger{name}
}

func (logger *Logger) Info(message string) {
	if activeLogLevel <= slog.LevelInfo {
		logLine := "INFO " + logger.name + " " + message
		internalLogger.Println(logLine)
		LogEmitter.Send(logLine)
	}
}

func (logger *Logger) Infof(format string, a ...any) {
	if activeLogLevel <= slog.LevelInfo {
		logLine := "INFO " + logger.name + " " + fmt.Sprintf(format, a...)
		internalLogger.Println(logLine)
		LogEmitter.Send(logLine)
	}
}

func (logger *Logger) Debug(message string) {
	if activeLogLevel <= slog.LevelDebug {
		logLine := "DEBUG " + logger.name + " " + message
		internalLogger.Println(logLine)
		LogEmitter.Send(logLine)
	}
}

func (logger *Logger) Debugf(format string, a ...any) {
	if activeLogLevel <= slog.LevelDebug {
		logLine := "DEBUG " + logger.name + " " + fmt.Sprintf(format, a...)
		internalLogger.Println(logLine)
		LogEmitter.Send(logLine)
	}
}

func (logger *Logger) Warn(message string) {
	if activeLogLevel <= slog.LevelWarn {
		logLine := "WARN " + logger.name + " " + message
		internalLogger.Println(logLine)
		LogEmitter.Send(logLine)
	}
}

func (logger *Logger) Warnf(format string, a ...any) {
	if activeLogLevel <= slog.LevelWarn {
		logLine := "WARN " + logger.name + " " + fmt.Sprintf(format, a...)
		internalLogger.Println(logLine)
		LogEmitter.Send(logLine)
	}
}

func (logger *Logger) Error(message string) {
	if activeLogLevel <= slog.LevelError {
		logLine := "ERROR " + logger.name + " " + message
		internalLogger.Println(logLine)
		LogEmitter.Send(logLine)
	}
}

func (logger *Logger) Errorf(format string, a ...any) {
	if activeLogLevel <= slog.LevelError {
		logLine := "ERROR " + logger.name + " " + fmt.Sprintf(format, a...)
		internalLogger.Println(logLine)
		LogEmitter.Send(logLine)
	}
}

func (logger *Logger) Name() string {
	return logger.name
}

// calls on the default logger

func Info(message string) {
	defaultLogger.Info(message)
}

func Infof(format string, a ...any) {
	defaultLogger.Infof(format, a...)
}

func Debug(message string) {
	defaultLogger.Debug(message)
}

func Debugf(format string, a ...any) {
	defaultLogger.Debugf(format, a...)
}

func Warn(message string) {
	defaultLogger.Warn(message)
}

func Warnf(format string, a ...any) {
	defaultLogger.Warnf(format, a...)
}

func Error(message string) {
	defaultLogger.Error(message)
}

func Errorf(format string, a ...any) {
	defaultLogger.Errorf(format, a...)
}

func parseLevel(level string) slog.Level {
	lcLevel := strings.ToLower(level)
	var result slog.Level

	switch lcLevel {
	case "error":
		result = slog.LevelError
	case "warn":
		result = slog.LevelWarn
	case "debug":
		result = slog.LevelDebug
	default:
		result = slog.LevelInfo
	}

	return result
}
