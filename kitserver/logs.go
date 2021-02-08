package kitserver

import (
	"context"
	"github.com/go-kit/kit/log"
	"net/http"
)

func LogRequest(logger log.Logger) func(context.Context, *http.Request) context.Context {
	return func(ctx context.Context, request *http.Request) context.Context {
		traceLogger := logger.(TraceLogger)
		ctx = context.WithValue(ctx, "logger", traceLogger)
		traceLogger.LogWithContext(ctx, "Path", request.URL.Path,
			"Query", request.URL.RawQuery,
			"Method", request.Method,
			"msg", "Received call")
		return ctx
	}
}

func LogResponse() func(context.Context, http.ResponseWriter) context.Context {
	return func(ctx context.Context, w http.ResponseWriter) context.Context {
		logger := ctx.Value("logger").(TraceLogger)
		logger.LogWithContext(ctx,
			"Headers", w.Header(),
			"msg", "Responded with")
		return ctx
	}
}
