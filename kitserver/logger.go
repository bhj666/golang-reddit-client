package kitserver

import (
	"context"
	"github.com/go-kit/kit/log"
)

type TraceLogger struct {
	logger log.Logger
}

func (l TraceLogger) LogWithContext(ctx context.Context, keyvals ...interface{}) error {
	return l.Log("Correlation-id", ctx.Value(CorrelationIdHeader),
		"Request-id", ctx.Value(RequestIdHeader),
		keyvals)
}

func (l TraceLogger) Log(keyvals ...interface{}) error {
	return l.logger.Log(keyvals)
}
