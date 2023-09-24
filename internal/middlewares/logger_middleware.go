package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/hytkgami/slog-tracer/internal/logger"
)

type ResponseWriterWrapper struct {
	rw     http.ResponseWriter
	writer io.Writer
}

func NewResponseWriterWrapper(rw http.ResponseWriter, buf io.Writer) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		rw:     rw,
		writer: io.MultiWriter(rw, buf),
	}
}

// Header implements http.ResponseWriter.
func (rww *ResponseWriterWrapper) Header() http.Header {
	return rww.rw.Header()
}

// Write implements http.ResponseWriter.
func (rww *ResponseWriterWrapper) Write(b []byte) (int, error) {
	return rww.writer.Write(b)
}

// WriteHeader implements http.ResponseWriter.
func (rww *ResponseWriterWrapper) WriteHeader(statusCode int) {
	rww.rw.WriteHeader(statusCode)
}

var _ http.ResponseWriter = (*ResponseWriterWrapper)(nil)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				case slog.MessageKey:
					a = slog.Attr{
						Key:   "message",
						Value: a.Value,
					}
				}
				return a
			},
		})
		// request body
		l := slog.New(logHandler)
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			l.ErrorContext(r.Context(), "failed to read request body", "error", err)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		// response body
		respBuf := bytes.NewBuffer(nil)
		rww := NewResponseWriterWrapper(w, respBuf)
		next.ServeHTTP(rww, r)
		go func() {
			if len(reqBody) == 0 {
				return
			}
			buf := bytes.NewBuffer(nil)
			err := json.Compact(buf, reqBody)
			if err != nil {
				l.ErrorContext(r.Context(), "failed to compact request body", "error", err)
			}
			l.InfoContext(r.Context(), buf.String(), "method", r.Method, "path", r.URL.Path)
		}()
		go func() {
			l.InfoContext(r.Context(), respBuf.String(), "method", r.Method, "path", r.URL.Path)
		}()
	})
}
