package middlewares

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

type ctxKey int

const ridKey ctxKey = ctxKey(0)

func GetReqID(ctx context.Context) string {
	return ctx.Value(ridKey).(string)
}

func RequestId(next http.Handler) http.Handler {
	fn := func(rw http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.NewV4().String()
		}
		ctx := context.WithValue(r.Context(), ridKey, rid)
		rw.Header().Add("X-Request-ID", rid)
		next.ServeHTTP(rw, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
func Logger(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			var bodyBuf bytes.Buffer
			ww.Tee(&bodyBuf)
			start := time.Now()
			defer func() {
				if ww.Status() != http.StatusOK {
					logger.Info().
						Str("request-id", GetReqID(r.Context())).
						Int("status", ww.Status()).
						Int("bytes", ww.BytesWritten()).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Str("query", r.URL.RawQuery).
						Str("ip", r.RemoteAddr).
						Str("trace.id", trace.SpanFromContext(r.Context()).SpanContext().TraceID().String()).
						Str("user-agent", r.UserAgent()).
						Dur("latency", time.Since(start)).
						Str("resp_body", bodyBuf.String()).
						Msg("request completed")
				}
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
