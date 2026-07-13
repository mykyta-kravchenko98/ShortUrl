package observability

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

type traceContextHandler struct {
	next slog.Handler
}

func newTraceContextHandler(next slog.Handler) *traceContextHandler {
	return &traceContextHandler{next: next}
}

func (h *traceContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *traceContextHandler) Handle(ctx context.Context, record slog.Record) error {
	if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
		record.AddAttrs(
			slog.String("trace_id", sc.TraceID().String()),
			slog.String("span_id", sc.SpanID().String()),
		)
	}
	return h.next.Handle(ctx, record)
}

func (h *traceContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newTraceContextHandler(h.next.WithAttrs(attrs))
}

func (h *traceContextHandler) WithGroup(name string) slog.Handler {
	return newTraceContextHandler(h.next.WithGroup(name))
}
