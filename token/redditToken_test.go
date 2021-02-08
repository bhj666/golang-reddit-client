package token

import (
	errors "aws-example/error"
	"aws-example/persistance"
	testutils "aws-example/test"
	"encoding/json"
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
	body, _ := json.Marshal(persistance.Token{
		TokenType:    "Bearer",
		AccessToken:  "AccessToken",
		ExpiresIn:    3600,
		RefreshToken: "refreshToken",
	})
	redditClientMock.ExchangeBody = string(body)
	handler.SecretsRepository.Save(persistance.Secret{Secret: SECRET})

	_, err := handler.exchangeToken(request{
		State: SECRET,
		Code:  "code"})

	require.Nil(test, err)
}

func TestBadRequest(test *testing.T) {
	handler := createHandler()
	redditClientMock.ExchangeStatus = http.StatusBadRequest
	redditClientMock.ExchangeBody = "{}"
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
