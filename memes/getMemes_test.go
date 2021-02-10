package memes

import (
	"aws-example/encryption"
	errors "aws-example/error"
	"aws-example/persistance"
	"aws-example/reddit"
	testutils "aws-example/test"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestCustomPaginationFlow(test *testing.T) {
	handler := getMemeHandlerWithValidTokenAndData()

	resp, err := handler.getMemes(request{
		Query:    "php",
		From:     "all",
		Page:     1,
		PageSize: 1,
	})

	require.Nil(test, err, fmt.Sprintf("Error should not be returned"))
	require.Equal(test, 1, len(resp.Data))
	require.Equal(test, int64(15), resp.Data[0].Score, "Should properly sort response by score")
}

func TestNoActiveTokenFlow(test *testing.T) {
	handler := getMemeHandlerWithOutdatedToken()

	_, err := handler.getMemes(request{
		Query:    "php",
		From:     "all",
		Page:     1,
		PageSize: 1,
	})

	require.NotNil(test, err, fmt.Sprintf("Error should be returned"))
	responseError, ok := err.(errors.ResponseError)
	require.True(test, ok, "Error should be of type ResponseError")
	require.Equal(test, http.StatusUnauthorized, responseError.Code())
}

func memeHandler() *memeSearchHandler {
	return &memeSearchHandler{
		TokenRepository: testutils.GetTokenRepositoryMock(),
		RedditClient:    &testutils.RedditClientMock{FindResults: make(map[string]testutils.FindResult)},
		TimeProvider:    testutils.TimeProviderMock{Time: 121},
	}
}

func getMemeHandlerWithValidTokenAndData() *memeSearchHandler {
	handler := memeHandler()
	handler.TokenRepository.Save(persistance.Token{
		TokenType:    "Bearer",
		AccessToken:  encryption.EncryptedString{"AccessToken"},
		ExpiresIn:    3600,
		ExpiresAt:    122,
		RefreshToken: encryption.EncryptedString{"refreshToken"},
	})
	body1 := reddit.SearchResponse{
		Data: reddit.SearchResponseData{
			After: "after1",
			Posts: []reddit.PostResponse{
				{
					Data: reddit.PostResponseData{
						Score: 1,
						Url:   "1",
					},
				},
				{
					Data: reddit.PostResponseData{
						Score: 2,
						Url:   "2",
					},
				},
				{
					Data: reddit.PostResponseData{
						Score: 3,
						Url:   "3",
					},
				},
			},
		},
	}
	body1String, _ := json.Marshal(body1)
	redditMock := &testutils.RedditClientMock{FindResults: make(map[string]testutils.FindResult)}
	redditMock.FindResults[""] = testutils.FindResult{
		FindStatus: 200,
		FindError:  nil,
		FindBody:   string(body1String),
	}
	body2 := reddit.SearchResponse{
		Data: reddit.SearchResponseData{
			After: "",
			Posts: []reddit.PostResponse{
				{
					Data: reddit.PostResponseData{
						Score: 0,
						Url:   "0",
					},
				},
				{
					Data: reddit.PostResponseData{
						Score: 100,
						Url:   "100",
					},
				},
				{
					Data: reddit.PostResponseData{
						Score: 15,
						Url:   "15",
					},
				},
			},
		},
	}
	body2String, _ := json.Marshal(body2)
	redditMock.FindResults["after1"] = testutils.FindResult{
		FindStatus: 200,
		FindError:  nil,
		FindBody:   string(body2String),
	}
	handler.RedditClient = redditMock
	return handler
}

func getMemeHandlerWithOutdatedToken() *memeSearchHandler {
	handler := memeHandler()
	handler.TokenRepository.Save(persistance.Token{
		TokenType:    "Bearer",
		AccessToken:  encryption.EncryptedString{"AccessToken"},
		ExpiresIn:    3600,
		ExpiresAt:    120,
		RefreshToken: encryption.EncryptedString{"refreshToken"},
	})
	return handler
}
