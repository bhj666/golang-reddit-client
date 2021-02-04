package main

import (
	"aws-example/persistance"
	testutils "aws-example/test"
	"aws-example/timeprovider"
	"net/http"
	"net/url"
	"testing"
)

var authHandler = AuthorizationHandler{
	SecretsRepository: testutils.GetSecretRepositoryMock(),
	TokenRepository:   testutils.GetTokenRepositoryMock(),
	RedditClient:      &testutils.RedditClientMock{},
	TimeProvider:      timeprovider.ProviderImpl{},
}

var requestMock = &http.Request{
	URL: &url.URL{
		Path: "",
	},
}

func TestAuthorizeNoValidTokenFlow(test *testing.T) {
	var responseMock = testutils.ResponseWriterMock{}
	authHandler.ServeHTTP(&responseMock, requestMock)

	if responseMock.Code != http.StatusSeeOther {
		test.Error("User was not redirected on authorization flow")
	}
}

func TestAuthorizeValidTokenFlow(test *testing.T) {
	var responseMock = testutils.ResponseWriterMock{}
	authHandler.TokenRepository.Save(persistance.Token{
		AccessToken: "access_token",
		ExpiresAt:   authHandler.TimeProvider.GetCurrentSeconds() + 1000,
	})
	authHandler.ServeHTTP(&responseMock, requestMock)

	if responseMock.Code != http.StatusConflict {
		test.Error("Should not attempt authorization if there is already valid token")
	}
}
