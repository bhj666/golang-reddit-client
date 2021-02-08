package kitserver

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const CorrelationIdHeader = "Correlation-id"
const RequestIdHeader = "Request-id"

func ExtractTraceInfo(next func(context.Context, *http.Request) context.Context) func(context.Context, *http.Request) context.Context {
	return func(ctx context.Context, request *http.Request) context.Context {
		if correlationId := request.Header.Get(CorrelationIdHeader); correlationId != "" {
			ctx = context.WithValue(ctx, CorrelationIdHeader, correlationId)
		} else {
			ctx = context.WithValue(ctx, CorrelationIdHeader, uuid.New().String())
		}
		if requestId := request.Header.Get(RequestIdHeader); requestId != "" {
			ctx = context.WithValue(ctx, RequestIdHeader, requestId)
		} else {
			ctx = context.WithValue(ctx, RequestIdHeader, uuid.New().String())
		}

		return next(ctx, request)
	}
}

func InsertTraceInfo(next func(context.Context, http.ResponseWriter) context.Context) func(context.Context, http.ResponseWriter) context.Context {
	return func(ctx context.Context, w http.ResponseWriter) context.Context {
		correlationId := ctx.Value(CorrelationIdHeader)
		w.Header().Set(CorrelationIdHeader, correlationId.(string))
		return next(ctx, w)
	}
}
