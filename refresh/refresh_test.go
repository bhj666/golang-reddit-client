package main

import (
	"aws-example/encryption"
	"aws-example/persistance"
	testutils "aws-example/test"
	"github.com/stretchr/testify/require"
	"testing"
)

var refreshedToken = persistance.Token{
	AccessToken: encryption.EncryptedString{"RefreshedAccessToken"},
	ExpiresIn:   3600,
	TokenType:   "Bearer",
}

var currentToken = persistance.Token{
	AccessToken:  encryption.EncryptedString{"CurrentToken"},
	RefreshToken: encryption.EncryptedString{"refreshToken"},
	ExpiresIn:    3600,
	ExpiresAt:    2,
	TokenType:    "Bearer",
}

var invalidToken = persistance.Token{
	AccessToken:  encryption.EncryptedString{"CurrentToken"},
	RefreshToken: encryption.EncryptedString{"refreshToken"},
	ExpiresIn:    3600,
	ExpiresAt:    0,
	TokenType:    "Bearer",
}

func getHandler(currentToken, refreshedToken persistance.Token) refreshHandlerImpl {
	tokenMock := testutils.GetTokenRepositoryMock()
	tokenMock.Save(currentToken)
	return refreshHandlerImpl{
		TokenRepository: tokenMock,
		RedditClient:    &testutils.RedditClientMock{RefreshError: nil, RefreshBody: &refreshedToken},
		TimeProvider:    testutils.TimeProviderMock{Time: 1},
	}
}

func TestHappyPathFlow(test *testing.T) {
	handler := getHandler(currentToken, refreshedToken)

	handler.refreshToken()()

	newToken := &persistance.Token{}
	handler.TokenRepository.FindActive(newToken, 3600)
	require.Equal(test, "RefreshedAccessToken", newToken.AccessToken.StringValue)
	require.Equal(test, "refreshToken", newToken.RefreshToken.StringValue)
}

func TestNoTokenToRefreshFlow(test *testing.T) {
	handler := getHandler(invalidToken, refreshedToken)

	handler.refreshToken()()

	newToken := &persistance.Token{}
	handler.TokenRepository.FindActive(newToken, 3600)
	require.Equal(test, "", newToken.AccessToken.StringValue)
}
