package token

import (
	errors "aws-example/error"
	"aws-example/persistance"
	"aws-example/reddit"
	"net/http"
	"time"
)

type tokenExchangeHandler struct {
	SecretsRepository persistance.SecretsRepository
	TokenRepository   persistance.TokenRepository
	RedditClient      reddit.Client
}

func newTokenExchangeHandler() tokenExchangeHandler {
	return tokenExchangeHandler{
		persistance.NewSecretsRepository(),
		persistance.NewTokenRepository(),
		reddit.NewClient(*http.DefaultClient),
	}
}

type request struct {
	Error string
	Code  string
	State string
}

type response struct {
}

func (h tokenExchangeHandler) exchangeToken(r request) (*response, error) {
	db := h.SecretsRepository
	if r.Error != "" {
		return nil, errors.InternalError{Message: r.Error}
	}
	secret := &persistance.Secret{}
	err := db.Find(r.State, secret)
	if err != nil || secret.Secret == "" {
		return nil, errors.GenericResponseError{
			Message:      "Such request not found",
			ResponseCode: http.StatusNotFound,
		}
	}
	err = db.Delete(*secret)
	if err != nil {
		return nil, errors.GenericResponseError{
			Message:      err.Error(),
			ResponseCode: 500,
		}
	}
	token, er := h.RedditClient.ExchangeForToken(r.Code)
	if er != nil {
		return nil, errors.InternalError{Message: er.Error()}
	}
	now := time.Now()
	token.ExpiresAt = now.Unix() + token.ExpiresIn
	tokenDb := h.TokenRepository
	err = tokenDb.Save(*token)
	if err != nil {
		return nil, errors.GenericResponseError{
			Message:      err.Error(),
			ResponseCode: 500,
		}
	}
	return nil, nil
}
