package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		span := trace.SpanFromContext(r.Context()).SpanContext()

		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rw.status),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
		}

		if span.IsValid() {
			attrs = append(attrs,
				slog.String("trace_id", span.TraceID().String()),
				slog.String("span_id", span.SpanID().String()),
			)
		}

		slog.LogAttrs(r.Context(), slog.LevelInfo, "request", attrs...)
	})
}
