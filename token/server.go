package token

import (
	"aws-example/kitserver"
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
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
	err := query.Get("error")
	code := query.Get("code")
	state := query.Get("state")

	return request{
		Error: err,
		Code:  code,
		State: state,
	}, nil
}

func encodeResponse(
	ctx context.Context, w http.ResponseWriter, response interface{},
) error {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Congratulations!. Now you can use this app"))
	return nil
}

func makeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return newTokenExchangeHandler().exchangeToken(r.(request))
	}
}
