package kitserver

import (
	errors "aws-example/error"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"
)

func GetOptions() []kithttp.ServerOption {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logInterface := logrus.NewLogrusLogger(logger)
	enrichedLogger := TraceLogger{logger: logInterface}
	return []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(enrichedLogger)),
		kithttp.ServerBefore(ExtractTraceInfo(LogRequest(enrichedLogger))),
		kithttp.ServerAfter(InsertTraceInfo(LogResponse())),
		kithttp.ServerErrorEncoder(errors.Encoder),
	}
}
