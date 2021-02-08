package authorize

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
		kithttp.NopRequestDecoder,
		encodeResponse,
		kitserver.GetOptions()...,
	)
	return server
}

func encodeResponse(
	ctx context.Context, w http.ResponseWriter, r interface{},
) error {
	resp, _ := r.(response)
	http.Redirect(w, &http.Request{Method: "GET"}, resp.RedirectUrl, http.StatusSeeOther)
	return nil
}

func makeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return newAuthorizationHandler().authorize()
	}
}
