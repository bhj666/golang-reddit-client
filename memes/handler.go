package memes

import (
	errors "aws-example/error"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func Handler() *kithttp.Server {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logrus.NewLogrusLogger(logger))),
		kithttp.ServerErrorEncoder(errors.ErrorEncoder),
	}

	handler := kithttp.NewServer(
		MakeEndpoint(),
		decodeRequest,
		encodeResponse,
		opts...,
	)
	return handler
}

func decodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	query := r.URL.Query()
	memeQuery := query.Get("query")
	from := query.Get("from")
	if memeQuery == "" || from == "" {

	}
	page, err := strconv.ParseInt(query.Get("page"), 0, 64)
	if err != nil {
		page = 0
	}
	pageSize, err := strconv.ParseInt(query.Get("pageSize"), 0, 64)
	if err != nil {
		pageSize = 25
	}

	return request{
		Query:    memeQuery,
		From:     from,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func encodeResponse(
	ctx context.Context, w http.ResponseWriter, response interface{},
) error {
	body, err := json.Marshal(response)
	if err != nil {
		return errors.InternalError{
			Message: fmt.Sprintf("Parsing response thrown error %v", err),
		}
	}
	w.Write(body)
	w.WriteHeader(200)
	return nil
}
