package middlewares

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/hytkgami/slog-tracer/internal/logger"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go func() {
			logHandler := logger.NewHandler(&logger.Options{
				AddSource: true,
				Level:     slog.LevelInfo,
				Writer:    os.Stdout,
			})
			l := slog.New(logHandler)
			l.InfoContext(r.Context(), "request", "method", r.Method, "path", r.URL.Path)
		}()
		next.ServeHTTP(w, r)
	})
}
