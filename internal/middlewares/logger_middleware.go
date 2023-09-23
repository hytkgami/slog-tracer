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
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					switch a.Key {
					case slog.LevelKey:
						a = slog.Attr{
							Key:   "severity",
							Value: a.Value,
						}
					case slog.SourceKey:
						a = slog.Attr{
							Key:   "logging.googleapis.com/sourceLocation",
							Value: a.Value,
						}
					}
					return a
				},
			})
			l := slog.New(logHandler)
			l.InfoContext(r.Context(), "request", "method", r.Method, "path", r.URL.Path)
		}()
		next.ServeHTTP(w, r)
	})
}
