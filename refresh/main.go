package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	handler := newRefreshHandler()
	lambda.Start(handler.refreshToken())
}
