package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler := NewRefreshHandler()
	lambda.Start(handler.RefreshToken())
}
