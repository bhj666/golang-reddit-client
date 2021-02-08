package errors

import (
	"context"
	"net/http"
)

type ResponseError interface {
	Code() int
	error
}

type UnauthorizedError struct {
}

func (UnauthorizedError) Error() string {
	return "You need to authorize first"
}

func (UnauthorizedError) Code() int {
	return http.StatusUnauthorized
}

type InternalError struct {
	Message string
}

func (e InternalError) Error() string {
	return e.Message
}

func (InternalError) Code() int {
	return http.StatusInternalServerError
}

type GenericResponseError struct {
	Message      string
	ResponseCode int
}

func (e GenericResponseError) Error() string {
	return e.Message
}

func (e GenericResponseError) Code() int {
	return e.ResponseCode
}

func Encoder(ctx context.Context, err error, w http.ResponseWriter) {
	//logger := ctx.Value("Logger").(log.Logger)
	//logger.Log(err)
	if responseError, ok := err.(ResponseError); ok {
		w.WriteHeader(responseError.Code())
		w.Write([]byte(responseError.Error()))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

}
