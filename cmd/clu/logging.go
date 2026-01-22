package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
)

// ConsoleHandler is a simple log handler with TIME LEVEL: message format
type ConsoleHandler struct {
	level slog.Level
}

func setupLogging(debug bool) {
	var handler *ConsoleHandler
	if debug {
		handler = &ConsoleHandler{level: slog.LevelDebug}
	} else {
		handler = &ConsoleHandler{level: slog.LevelInfo}
	}

	// Replace the default logger
	slog.SetDefault(slog.New(handler))

	slog.Debug("debug logging is enabled")

}

func (h *ConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add source info if debug level
	if r.Level == slog.LevelDebug && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		fmt.Fprintf(os.Stderr, "[%s:%d] ", filepath.Base(f.File), f.Line)
	}

	fmt.Fprintf(os.Stderr, "%s %s: %s",
		r.Time.Format("15:04:05"),
		r.Level.String(),
		r.Message,
	)

	// Handle key-value pairs
	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(os.Stderr, " %s=%v", a.Key, a.Value)
		return true
	})

	fmt.Fprintf(os.Stderr, "\n")
	return nil
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *ConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}
