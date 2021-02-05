package memes

import (
	errors "aws-example/error"
	"aws-example/persistance"
	"aws-example/reddit"
	testutils "aws-example/test"
	"encoding/json"
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
	if err != nil {
		test.Errorf("Error should not be returned, got %s", err.Error())
	}
	if len(resp.Data) != 1 {
		test.Errorf("Should return only one result got %v instead", len(resp.Data))
	}
	if resp.Data[0].Score != 15 {
		test.Error("Should properly sort response by score")
	}
}

func TestNoActiveTokenFlow(test *testing.T) {
	handler := getMemeHandlerWithOutdatedToken()

	_, err := handler.getMemes(request{
		Query:    "php",
		From:     "all",
		Page:     1,
		PageSize: 1,
	})
	if err == nil {
		test.Errorf("Error should be returned")
	}
	responseError, ok := err.(errors.ResponseError)
	if !ok {
		test.Errorf("Error should be of type ResponseError")
	}
	if responseError.Code() != http.StatusUnauthorized {
		test.Errorf("Code should be 401 got  %v instead", responseError.Code())
	}
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
		AccessToken:  "AccessToken",
		ExpiresIn:    3600,
		ExpiresAt:    122,
		RefreshToken: "RefreshToken",
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
		AccessToken:  "AccessToken",
		ExpiresIn:    3600,
		ExpiresAt:    120,
		RefreshToken: "RefreshToken",
	})
	return handler
}
