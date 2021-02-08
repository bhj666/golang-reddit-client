package errors

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Handle(w http.ResponseWriter, statusCode int, err error) bool {
	if err != nil {
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(err.Error()))
		log.Error(err)
	}
	return err != nil
}

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
	if responseError, ok := err.(ResponseError); ok {
		w.WriteHeader(responseError.Code())
		w.Write([]byte(responseError.Error()))
	} else {
		w.WriteHeader(responseError.Code())
		w.Write([]byte(err.Error()))
	}

}
