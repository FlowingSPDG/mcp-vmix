package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GetLogFilePath は%appdata%/RSLT/vmix-mcp/logs/にログファイルのパスを生成します
func GetLogFilePath() (string, error) {
	appData, _ := os.UserConfigDir()
	logDir := filepath.Join(appData, "RSLT", "vmix-mcp", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create log directory: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	logPath := filepath.Join(logDir, fmt.Sprintf("vmix_mcp_log_%s.log", timestamp))
	return logPath, nil
}

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
	if _, err := l.f.WriteString(fmt.Sprintf("[DEBUG] %s\n", msg)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write debug log: %v\n", err)
	}
}

// Error implements Logger.
func (l *fileLogger) Error(msg string) {
	if _, err := l.f.WriteString(fmt.Sprintf("[ERROR] %s\n", msg)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write error log: %v\n", err)
	}
}

// Info implements Logger.
func (l *fileLogger) Info(msg string) {
	if _, err := l.f.WriteString(fmt.Sprintf("[INFO] %s\n", msg)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write info log: %v\n", err)
	}
}

// Warn implements Logger.
func (l *fileLogger) Warn(msg string) {
	if _, err := l.f.WriteString(fmt.Sprintf("[WARN] %s\n", msg)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write warning log: %v\n", err)
	}
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
