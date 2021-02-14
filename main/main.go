package main

import (
	"aws-example/authorize"
	"aws-example/memes"
	"aws-example/produce"
	"aws-example/token"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	"net/http"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	router := mux.NewRouter()
	router.Handle("/api/memes", memes.Server()).Methods(http.MethodGet)
	router.Handle("/api/memes", produce.Server()).Methods(http.MethodPost)
	router.Handle("/api/token/", token.Server()).Methods(http.MethodGet)
	router.Handle("/api/authorize", authorize.Server()).Methods(http.MethodGet)
	adapter := gorillamux.New(router)
	r, err := adapter.Proxy(request)
	return r, err
}

func main() {
	lambda.Start(handler)
}
