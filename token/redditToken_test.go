package token

import (
	"aws-example/encryption"
	errors "aws-example/error"
	"aws-example/persistance"
	testutils "aws-example/test"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func createHandler() tokenExchangeHandler {
	return tokenExchangeHandler{
		SecretsRepository: testutils.GetSecretRepositoryMock(),
		TokenRepository:   testutils.GetTokenRepositoryMock(),
		RedditClient:      redditClientMock,
	}
}

var redditClientMock = &testutils.RedditClientMock{}
var SECRET = "secret"

func TestHappyFlow(test *testing.T) {
	handler := createHandler()
	redditClientMock.ExchangeStatus = http.StatusOK
	token := persistance.Token{
		TokenType:    "Bearer",
		AccessToken:  encryption.EncryptedString{"AccessToken"},
		ExpiresIn:    3600,
		RefreshToken: encryption.EncryptedString{"refreshToken"},
	}
	redditClientMock.ExchangeBody = &token
	handler.SecretsRepository.Save(persistance.Secret{Secret: SECRET})

	_, err := handler.exchangeToken(request{
		State: SECRET,
		Code:  "code"})

	require.Nil(test, err)
}

func TestBadRequest(test *testing.T) {
	handler := createHandler()
	redditClientMock.ExchangeError = errors.GenericResponseError{ResponseCode: 400, Message: "Invalid"}
	handler.SecretsRepository.Save(persistance.Secret{Secret: SECRET})

	_, err := handler.exchangeToken(request{
		State: SECRET,
		Code:  "code",
	})

	require.NotNil(test, err)
	require.Equal(test, http.StatusInternalServerError, err.(errors.ResponseError).Code())
}

func TestSecretNotFound(test *testing.T) {
	handler := createHandler()

	_, err := handler.exchangeToken(request{
		State: SECRET,
		Code:  "code",
	})

	require.NotNil(test, err)
	require.Equal(test, http.StatusNotFound, err.(errors.ResponseError).Code())
}

func TestRedirectWithError(test *testing.T) {
	handler := createHandler()

	_, err := handler.exchangeToken(request{
		Error: "Wrong code",
	})

	require.NotNil(test, err)
	require.Equal(test, http.StatusInternalServerError, err.(errors.ResponseError).Code())
}
