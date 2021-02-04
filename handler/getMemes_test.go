package main

import (
	"aws-example/persistance"
	"aws-example/reddit"
	testutils "aws-example/test"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

var memeRequestMock = &http.Request{
	URL: &url.URL{
		Path: "",
	},
}

func TestDefaultPaginationFlow(test *testing.T) {
	handler := getMemeHandlerWithValidTokenAndData()
	var responseMock = testutils.ResponseWriterMock{}
	memeRequestMock.URL.RawQuery = ""

	handler.ServeHTTP(&responseMock, memeRequestMock)

	response := make([]reddit.PostResponseData, 0)
	if responseMock.Code != http.StatusOK {
		test.Errorf("Should return correct result got %v instead", responseMock.Code)
	}
	_ = json.Unmarshal([]byte(responseMock.Data), &response)
	if len(response) != 6 {
		test.Errorf("Should return only one result got %s instead", responseMock.Data)
	}
	if response[0].Score != 100 {
		test.Error("Should properly sort response by score")
	}
}

func TestCustomPaginationFlow(test *testing.T) {
	handler := getMemeHandlerWithValidTokenAndData()
	var responseMock = testutils.ResponseWriterMock{}
	memeRequestMock.URL.RawQuery = "page=1&pageSize=1"

	handler.ServeHTTP(&responseMock, memeRequestMock)

	response := make([]reddit.PostResponseData, 0)
	if responseMock.Code != http.StatusOK {
		test.Errorf("Should return correct result got %v instead", responseMock.Code)
	}
	_ = json.Unmarshal([]byte(responseMock.Data), &response)
	if len(response) != 1 {
		test.Errorf("Should return only one result got %s instead", responseMock.Data)
	}
	if response[0].Score != 15 {
		test.Error("Should properly sort response by score")
	}
}

func TestNoActiveTokenFlow(test *testing.T) {
	handler := getMemeHandlerWithOutdatedToken()
	var responseMock = testutils.ResponseWriterMock{}

	handler.ServeHTTP(&responseMock, requestMock)

	if responseMock.Code != http.StatusUnauthorized {
		test.Errorf("Should return unauthorized error got %v instead", responseMock.Code)
	}
}

func memeHandler() *MemeSearchHandler {
	return &MemeSearchHandler{
		TokenRepository: testutils.GetTokenRepositoryMock(),
		RedditClient:    &testutils.RedditClientMock{FindResults: make(map[string]testutils.FindResult)},
		TimeProvider:    testutils.TimeProviderMock{Time: 121},
	}
}

func getMemeHandlerWithValidTokenAndData() *MemeSearchHandler {
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

func getMemeHandlerWithOutdatedToken() *MemeSearchHandler {
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
