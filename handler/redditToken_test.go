package main

import (
	"aws-example/persistance"
	testutils "aws-example/test"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func createHandler() TokenExchangeHandler {
	return TokenExchangeHandler{
		SecretsRepository: testutils.GetSecretRepositoryMock(),
		TokenRepository:   testutils.GetTokenRepositoryMock(),
		RedditClient:      redditClientMock,
	}
}

var redditClientMock = &testutils.RedditClientMock{}

var requestWithErrorMock = &http.Request{
	URL: &url.URL{
		RawQuery: "error=unauthorized",
	},
}
var SECRET = "secret"

var requestCorrectMock = &http.Request{
	URL: &url.URL{
		RawQuery: fmt.Sprintf("code=asd&state=%s", SECRET),
	},
}

func TestHappyFlow(test *testing.T) {
	handler := createHandler()
	redditClientMock.ExchangeStatus = http.StatusOK
	body, _ := json.Marshal(persistance.Token{
		TokenType:    "Bearer",
		AccessToken:  "AccessToken",
		ExpiresIn:    3600,
		RefreshToken: "RefreshToken",
	})
	redditClientMock.ExchangeBody = string(body)
	var responseMock = testutils.ResponseWriterMock{}
	handler.SecretsRepository.Save(persistance.Secret{Secret: SECRET})

	handler.ServeHTTP(&responseMock, requestCorrectMock)

	if responseMock.Code != http.StatusCreated {
		test.Error("Token exchange should fail with the same error as exchange call")
	}
}

func TestBadRequest(test *testing.T) {
	handler := createHandler()
	redditClientMock.ExchangeStatus = http.StatusBadRequest
	redditClientMock.ExchangeBody = "{}"
	var responseMock = testutils.ResponseWriterMock{}
	handler.SecretsRepository.Save(persistance.Secret{Secret: SECRET})

	handler.ServeHTTP(&responseMock, requestCorrectMock)

	if responseMock.Code != http.StatusBadRequest {
		test.Error("Token exchange should fail with the same error as exchange call")
	}
}

func TestSecretNotFound(test *testing.T) {
	handler := createHandler()
	var responseMock = testutils.ResponseWriterMock{}

	handler.ServeHTTP(&responseMock, requestCorrectMock)

	if responseMock.Code != http.StatusNotFound {
		test.Error("Token exchange should fail with 500")
	}
}

func TestRedirectWithError(test *testing.T) {
	handler := createHandler()
	var responseMock = testutils.ResponseWriterMock{}

	handler.ServeHTTP(&responseMock, requestWithErrorMock)

	if responseMock.Code != http.StatusInternalServerError {
		test.Error("Token exchange should fail with 500")
	}
}
