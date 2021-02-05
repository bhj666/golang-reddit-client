package main

import (
	"aws-example/memes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("Received call")
	log.Printf("Path: %s", request.Path)
	log.Printf("Query: %s", request.QueryStringParameters)
	log.Printf("Method: %s", request.HTTPMethod)
	log.Print(request.Body)
	router := mux.NewRouter()
	router.Handle("/api/memes", memes.Handler()).Methods(http.MethodGet)
	router.Handle("/api/token", NewTokenExchangeHandler()).Methods(http.MethodGet)
	router.Handle("/api/authorize", NewAuthorizationHandler()).Methods(http.MethodGet)
	adapter := gorillamux.New(router)
	r, err := adapter.Proxy(request)
	return r, err
}

func main() {
	lambda.Start(handler)
}
