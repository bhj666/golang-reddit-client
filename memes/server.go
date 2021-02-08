package memes

import (
	errors "aws-example/error"
	"aws-example/kitserver"
	"context"
	"encoding/json"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"strconv"
)

func Server() *kithttp.Server {
	server := kithttp.NewServer(
		makeEndpoint(),
		decodeRequest,
		encodeResponse,
		kitserver.GetOptions()...,
	)
	return server
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
