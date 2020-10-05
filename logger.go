package redis

import (
	"context"
	"fmt"
	"log"
	"os"
)

// Logger is a minimal interface for a logger
type Logger interface {
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
}

// LogWrapper wraps a standard Logger to implement the Logger interface.
type LogWrapper struct {
	Logger *log.Logger
}

// NewLogWrapper returns a new LogWrapper.
func NewLogWrapper() *LogWrapper {
	return &LogWrapper{
		Logger: log.New(os.Stdout, "", log.Flags()),
	}
}

// Info implements the Logger interface.
func (l *LogWrapper) Info(ctx context.Context, format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("INFO  %s", format), args...)
}

// Warn implements the Logger interface.
func (l *LogWrapper) Warn(ctx context.Context, format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("WARN  %s", format), args...)
}

// Error implements the Logger interface.
func (l *LogWrapper) Error(ctx context.Context, format string, args ...interface{}) {
	log.Printf(fmt.Sprintf("ERROR %s", format), args...)
}
