package authorize

import (
	errors "aws-example/error"
	"aws-example/persistance"
	"aws-example/reddit"
	"aws-example/timeprovider"
	"github.com/google/uuid"
)

type response struct {
	RedirectUrl string
}

type authorizationHandler struct {
	SecretsRepository persistance.SecretsRepository
	TokenRepository   persistance.TokenRepository
	RedditClient      reddit.Client
	TimeProvider      timeprovider.Provider
}

func newAuthorizationHandler() authorizationHandler {
	return authorizationHandler{
		persistance.NewSecretsRepository(),
		persistance.NewTokenRepository(),
		reddit.NewClient(),
		timeprovider.ProviderImpl{},
	}
}

func (h authorizationHandler) authorize() (*response, error) {
	tokensDb := h.TokenRepository
	token := persistance.Token{}
	tokensDb.FindActive(&token, h.TimeProvider.GetCurrentSeconds())
	if token.AccessToken != "" {
		return nil, errors.GenericResponseError{
			Message:      "Active token already exists",
			ResponseCode: 409,
		}
	}
	db := h.SecretsRepository
	secret := uuid.New().String()
	db.Save(persistance.Secret{Secret: secret})
	return &response{RedirectUrl: h.RedditClient.GetRedirectUrl(secret)}, nil
}
