package utils

import (
	"log"
	"os"
)

// Logger represents a logger
type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

// NewLogger creates a new logger
func NewLogger() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		warn:  log.New(os.Stdout, "[WARN] ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.info.Printf(msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.warn.Printf(msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.error.Printf(msg, args...)
}
