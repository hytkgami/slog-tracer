package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
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
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		trace := fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"), sc.TraceID().String())
		record.AddAttrs(
			slog.String("logging.googleapis.com/trace", trace),
			slog.String("logging.googleapis.com/spanId", sc.SpanID().String()),
			slog.Bool("logging.googleapis.com/trace_sampled", sc.TraceFlags().IsSampled()),
		)
	}
	record.AddAttrs(
		slog.Group("logging.googleapis.com/labels",
			slog.String("uid", "1234567890"),
		),
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
