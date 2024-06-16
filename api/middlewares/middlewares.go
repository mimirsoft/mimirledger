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
	a, ok := ctx.Value(ridKey).(string)
	if !ok {
		panic("requestID not found in context - failed type assertion")
	}

	return a
}

func RequestId(next http.Handler) http.Handler {
	handlerFn := func(res http.ResponseWriter, req *http.Request) {
		rid := req.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.NewV4().String()
		}

		ctx := context.WithValue(req.Context(), ridKey, rid)

		res.Header().Add("X-Request-ID", rid)
		next.ServeHTTP(res, req.WithContext(ctx))
	}

	return http.HandlerFunc(handlerFn)
}
func Logger(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handlerFn := func(res http.ResponseWriter, req *http.Request) {
			start := time.Now()
			wrapWriter := middleware.NewWrapResponseWriter(res, req.ProtoMajor)

			var bodyBuf bytes.Buffer

			wrapWriter.Tee(&bodyBuf)
			defer func() {
				if wrapWriter.Status() != http.StatusOK {
					logger.Info().
						Str("request-id", GetReqID(req.Context())).
						Int("status", wrapWriter.Status()).
						Int("bytes", wrapWriter.BytesWritten()).
						Str("method", req.Method).
						Str("path", req.URL.Path).
						Str("query", req.URL.RawQuery).
						Str("ip", req.RemoteAddr).
						Str("trace.id", trace.SpanFromContext(req.Context()).SpanContext().TraceID().String()).
						Str("user-agent", req.UserAgent()).
						Dur("latency", time.Since(start)).
						Str("resp_body", bodyBuf.String()).
						Msg("request completed")
				}
			}()

			next.ServeHTTP(wrapWriter, req)
		}

		return http.HandlerFunc(handlerFn)
	}
}
