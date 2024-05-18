package logger

import (
	"log/slog"
	"path/filepath"
	"strings"
)

// GetLogger will return logger instance
func GetLogger() *slog.Logger {
	if logger == nil {
		panic("initialize logger first")
	}

	return logger
}

func getLevel(s string) slog.Level {
	s = strings.ToUpper(s)

	switch s {
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source, ok := a.Value.Any().(*slog.Source)
		if ok {
			source.File = filepath.Base(source.File)
		}
	}

	return a
}
