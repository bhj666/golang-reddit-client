package main

import (
	"aws-example/persistance"
	testutils "aws-example/test"
	"encoding/json"
	"testing"
)

var refreshedToken = persistance.Token{
	AccessToken: "RefreshedAccessToken",
	ExpiresIn:   3600,
	TokenType:   "Bearer",
}

var currentToken = persistance.Token{
	AccessToken:  "CurrentToken",
	RefreshToken: "RefreshToken",
	ExpiresIn:    3600,
	ExpiresAt:    2,
	TokenType:    "Bearer",
}

var invalidToken = persistance.Token{
	AccessToken:  "CurrentToken",
	RefreshToken: "RefreshToken",
	ExpiresIn:    3600,
	ExpiresAt:    0,
	TokenType:    "Bearer",
}

func getHandler(currentToken, refreshedToken persistance.Token) RefreshHandlerImpl {
	tokenMock := testutils.GetTokenRepositoryMock()
	refreshResponse, _ := json.Marshal(refreshedToken)
	tokenMock.Save(currentToken)
	return RefreshHandlerImpl{
		TokenRepository: tokenMock,
		RedditClient:    &testutils.RedditClientMock{RefreshError: nil, RefreshStatus: 200, RefreshBody: string(refreshResponse)},
		TimeProvider:    testutils.TimeProviderMock{Time: 1},
	}
}

func TestHappyPathFlow(test *testing.T) {
	handler := getHandler(currentToken, refreshedToken)

	handler.RefreshToken()()

	newToken := &persistance.Token{}
	handler.TokenRepository.FindActive(newToken, 3600)
	if newToken.AccessToken != "RefreshedAccessToken" {
		test.Error("Token was not refreshed")
	}
	if newToken.RefreshToken != "RefreshToken" {
		test.Error("Refresh token was not propagated")
	}

}

func TestNoTokenToRefreshFlow(test *testing.T) {
	handler := getHandler(invalidToken, refreshedToken)

	handler.RefreshToken()()

	newToken := &persistance.Token{}
	handler.TokenRepository.FindActive(newToken, 3600)
	if newToken.AccessToken != "" {
		test.Error("Token should not be refreshed")
	}
}
