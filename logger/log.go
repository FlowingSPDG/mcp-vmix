package logger

import (
	"fmt"
	"os"
)

type Logger interface {
	Close() error
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

type fileLogger struct {
	f *os.File
}

// Debug implements Logger.
func (l *fileLogger) Debug(msg string) {
	l.f.WriteString(fmt.Sprintf("[DEBUG] %s\n", msg))
}

// Error implements Logger.
func (l *fileLogger) Error(msg string) {
	l.f.WriteString(fmt.Sprintf("[ERROR] %s\n", msg))
}

// Info implements Logger.
func (l *fileLogger) Info(msg string) {
	l.f.WriteString(fmt.Sprintf("[INFO] %s\n", msg))
}

// Warn implements Logger.
func (l *fileLogger) Warn(msg string) {
	l.f.WriteString(fmt.Sprintf("[WARN] %s\n", msg))
}

func NewFileLogger(path string) (Logger, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &fileLogger{f: f}, nil
}

func (l *fileLogger) Close() error {
	return l.f.Close()
}
