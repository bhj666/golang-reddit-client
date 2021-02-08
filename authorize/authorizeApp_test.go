package authorize

import (
	errors "aws-example/error"
	"aws-example/persistance"
	testutils "aws-example/test"
	"aws-example/timeprovider"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var authHandler = authorizationHandler{
	SecretsRepository: testutils.GetSecretRepositoryMock(),
	TokenRepository:   testutils.GetTokenRepositoryMock(),
	RedditClient:      &testutils.RedditClientMock{},
	TimeProvider:      timeprovider.ProviderImpl{},
}

func TestAuthorizeNoValidTokenFlow(test *testing.T) {
	resp, err := authHandler.authorize()

	require.Nil(test, err, "User was not redirected on authorization flow")
	require.NotEqual(test, "", resp.RedirectUrl, "Authorization url not returned")
}

func TestAuthorizeValidTokenFlow(test *testing.T) {
	authHandler.TokenRepository.Save(persistance.Token{
		AccessToken: "access_token",
		ExpiresAt:   authHandler.TimeProvider.GetCurrentSeconds() + 1000,
	})
	_, err := authHandler.authorize()
	require.NotNil(test, err, "Once there is already valid token endpoint should return error")
	require.Equal(test, http.StatusConflict, err.(errors.GenericResponseError).ResponseCode, "Should not attempt authorization if there is already valid token")
}
