package logger

import (
	"context"
	"io"
	"log/slog"
)

type (
	Handler struct {
		handler slog.Handler
	}
	Options struct {
		AddSource   bool
		Level       slog.Level
		Writer      io.Writer
		ReplaceAttr func(groups []string, a slog.Attr) slog.Attr
	}
)

func NewHandler(o *Options) *Handler {
	handlerOptions := slog.HandlerOptions{
		AddSource:   o.AddSource,
		Level:       o.Level,
		ReplaceAttr: o.ReplaceAttr,
	}
	return &Handler{handler: slog.NewJSONHandler(o.Writer, &handlerOptions)}
}

// Enabled implements slog.Handler.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler.
func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(
		slog.Bool("logging.googleapis.com/trace_sampled", true),
	)
	return h.handler.Handle(ctx, record)
}

// WithAttrs implements slog.Handler.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

// WithGroup implements slog.Handler.
func (h *Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}

var _ slog.Handler = (*Handler)(nil)
