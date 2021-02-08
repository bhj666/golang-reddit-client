package token

import (
	errors "aws-example/error"
	"aws-example/persistance"
	"aws-example/reddit"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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
		reddit.NewClient(),
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
	db.Find(r.State, secret)
	if secret.Secret == "" {
		return nil, errors.GenericResponseError{
			Message:      "Such request not found",
			ResponseCode: http.StatusNotFound,
		}
	}
	log.Printf("Secret: %v", secret)
	db.Delete(*secret)
	resp, er := h.RedditClient.ExchangeForToken(r.Code)
	if er != nil {
		return nil, errors.InternalError{Message: er.Error()}
	}
	response, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		return nil, errors.InternalError{Message: er.Error()}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.InternalError{Message: fmt.Sprintf("Response code: %v", resp.StatusCode)}
	}
	var token persistance.Token
	er = json.Unmarshal(response, &token)
	if er != nil {
		return nil, errors.InternalError{Message: er.Error()}
	}
	now := time.Now()
	token.ExpiresAt = now.Unix() + token.ExpiresIn
	tokenDb := h.TokenRepository
	tokenDb.Save(token)
	return nil, nil
}
