package authorize

import (
	"aws-example/kitserver"
	"context"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
)

const RequestKey = "request"

func Server() *kithttp.Server {
	server := kithttp.NewServer(
		makeEndpoint(),
		kithttp.NopRequestDecoder,
		encodeResponse,
		GetOptionsEnriched()...,
	)
	return server
}

func GetOptionsEnriched() []kithttp.ServerOption {
	options := kitserver.GetOptions()
	options = append(options, kithttp.ServerBefore(putRequestInCtx))
	return options
}

func putRequestInCtx(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, RequestKey, r)
}

func encodeResponse(
	ctx context.Context, w http.ResponseWriter, r interface{},
) error {
	resp, _ := r.(*response)
	request := ctx.Value(RequestKey).(*http.Request)
	http.Redirect(w, request, resp.RedirectUrl, http.StatusSeeOther)
	return nil
}

func makeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return newAuthorizationHandler().authorize()
	}
}
