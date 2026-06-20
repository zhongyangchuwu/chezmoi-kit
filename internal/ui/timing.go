package ui

import (
	"log/slog"
	"os"
	"time"
)

type syncTimingLogger struct {
	logger *slog.Logger
	file   *os.File
	path   string
}

func newSyncTimingLogger(enabled bool) (*syncTimingLogger, error) {
	if !enabled {
		return nil, nil
	}
	file, err := os.CreateTemp("", "cm-sync-*.log")
	if err != nil {
		return nil, err
	}
	handler := slog.NewTextHandler(file, &slog.HandlerOptions{Level: slog.LevelInfo})
	return &syncTimingLogger{logger: slog.New(handler), file: file, path: file.Name()}, nil
}

func (l *syncTimingLogger) Enabled() bool {
	return l != nil && l.file != nil
}

func (l *syncTimingLogger) Path() string {
	if l == nil {
		return ""
	}
	return l.path
}

func (l *syncTimingLogger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	return l.file.Close()
}

func (l *syncTimingLogger) Info(message string, attrs ...any) {
	if l == nil || l.logger == nil {
		return
	}
	l.logger.Info(message, attrs...)
}

func elapsed(start time.Time) time.Duration {
	return time.Since(start).Round(time.Millisecond)
}
