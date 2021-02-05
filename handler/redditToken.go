package main

import (
	errors "aws-example/error"
	"aws-example/persistance"
	"aws-example/reddit"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type TokenExchangeHandler struct {
	SecretsRepository persistance.SecretsRepository
	TokenRepository   persistance.TokenRepository
	RedditClient      reddit.Client
}

func NewTokenExchangeHandler() http.Handler {
	return TokenExchangeHandler{
		persistance.NewSecretsRepository(),
		persistance.NewTokenRepository(),
		reddit.NewClient(),
	}
}

func (h TokenExchangeHandler) ServeHTTP(responseWriter http.ResponseWriter,
	request *http.Request) {

	db := h.SecretsRepository
	query := request.URL.Query()
	err := query.Get("error")
	code := query.Get("code")
	state := query.Get("state")
	if err != "" {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		_, _ = responseWriter.Write([]byte(err))
		log.Error(err)
		return
	}
	secret := &persistance.Secret{}
	db.Find(state, secret)
	if secret.Secret == "" {
		responseWriter.WriteHeader(http.StatusNotFound)
		_, _ = responseWriter.Write([]byte("Such request not found"))
		log.Error(err)
		return
	}
	log.Printf("Secret: %v", secret)
	db.Delete(*secret)
	resp, er := h.RedditClient.ExchangeForToken(code)
	if errors.Handle(responseWriter, http.StatusInternalServerError, er) {
		return
	}
	response, er := ioutil.ReadAll(resp.Body)
	if errors.Handle(responseWriter, http.StatusInternalServerError, er) {
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Response code: %v", resp.StatusCode)
		responseWriter.WriteHeader(resp.StatusCode)
		responseWriter.Write(response)
		return
	}
	var token persistance.Token
	er = json.Unmarshal(response, &token)
	if errors.Handle(responseWriter, http.StatusInternalServerError, er) {
		return
	}
	now := time.Now()
	token.ExpiresAt = now.Unix() + token.ExpiresIn
	tokenDb := h.TokenRepository
	tokenDb.Save(token)
	responseWriter.WriteHeader(http.StatusCreated)
	responseWriter.Write([]byte("Congratulations!. Now you can use this app"))
}
