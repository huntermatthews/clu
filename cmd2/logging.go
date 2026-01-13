package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// CustomHandler is a simple log handler with TIME LEVEL: message format
type CustomHandler struct {
	level slog.Level
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	fmt.Fprintf(os.Stderr, "%s %s: %s\n",
		r.Time.Format("15:04:05"),
		r.Level.String(),
		r.Message,
	)
	return nil
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}
