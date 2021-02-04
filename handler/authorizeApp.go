package main

import (
	"aws-example/persistance"
	"aws-example/reddit"
	"aws-example/timeprovider"
	"github.com/google/uuid"
	"net/http"
)

type AuthorizationHandler struct {
	SecretsRepository persistance.SecretsRepository
	TokenRepository   persistance.TokenRepository
	RedditClient      reddit.Client
	TimeProvider      timeprovider.Provider
}

func NewAuthorizationHandler() http.Handler {
	return AuthorizationHandler{
		persistance.NewSecretsRepository(),
		persistance.NewTokenRepository(),
		reddit.NewClient(),
		timeprovider.ProviderImpl{},
	}
}

func (h AuthorizationHandler) ServeHTTP(responseWriter http.ResponseWriter,
	request *http.Request) {
	tokensDb := h.TokenRepository
	token := persistance.Token{}
	tokensDb.FindActive(&token, h.TimeProvider.GetCurrentSeconds())
	if token.AccessToken != "" {
		responseWriter.WriteHeader(409)
		responseWriter.Write([]byte("Active token already exists"))
		return
	}
	db := h.SecretsRepository
	secret := uuid.New().String()
	db.Save(persistance.Secret{Secret: secret})
	redirectUrl := h.RedditClient.GetRedirectUrl(secret)
	http.Redirect(responseWriter, request, redirectUrl, http.StatusSeeOther)

}
