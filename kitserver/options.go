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
	return []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logrus.NewLogrusLogger(logger))),
		kithttp.ServerErrorEncoder(errors.Encoder),
	}
}
