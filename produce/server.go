package produce

import (
	"aws-example/kitserver"
	"context"
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
	subreddit := query.Get("subreddit")
	if memeQuery == "" || from == "" {

	}
	pageSize, err := strconv.ParseInt(query.Get("pageSize"), 0, 64)
	if err != nil {
		pageSize = 25
	}

	return request{
		Query:     memeQuery,
		From:      from,
		Subreddit: subreddit,
		PageSize:  pageSize,
	}, nil
}

func encodeResponse(
	ctx context.Context, w http.ResponseWriter, response interface{},
) error {
	w.WriteHeader(201)
	return nil
}
