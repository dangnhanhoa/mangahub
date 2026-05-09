package utils

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// NewLogger creates a structured JSON logger writing to both stdout and a log file.
// level: "debug" | "info" | "warn" | "error"
func NewLogger(level, logDir string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: lvl}

	writers := []io.Writer{os.Stdout}
	if logDir != "" {
		if err := os.MkdirAll(logDir, 0o755); err == nil {
			f, err := os.OpenFile(
				filepath.Join(logDir, "server.log"),
				os.O_CREATE|os.O_APPEND|os.O_WRONLY,
				0o644,
			)
			if err == nil {
				writers = append(writers, f)
			}
		}
	}

	w := io.MultiWriter(writers...)
	return slog.New(slog.NewJSONHandler(w, opts))
}
