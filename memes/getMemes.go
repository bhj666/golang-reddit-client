package memes

import (
	errors "aws-example/error"
	"aws-example/persistance"
	"aws-example/reddit"
	"aws-example/timeprovider"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"io/ioutil"
	"net/http"
	"sort"
)

type memeSearchHandler struct {
	TokenRepository persistance.TokenRepository
	RedditClient    reddit.Client
	TimeProvider    timeprovider.Provider
}

func newMemeSearchHandler() memeSearchHandler {
	return memeSearchHandler{
		persistance.NewTokenRepository(),
		reddit.NewClient(),
		timeprovider.ProviderImpl{},
	}
}

type request struct {
	Query    string
	From     string
	Page     int64
	PageSize int64
}

type response struct {
	Data []reddit.PostResponseData
}

func (h memeSearchHandler) getMemes(request request) (*response, error) {
	db := h.TokenRepository
	token := persistance.Token{}
	db.FindActive(&token, h.TimeProvider.GetCurrentSeconds())
	if token.AccessToken.StringValue == "" {
		return nil, errors.UnauthorizedError{}
	}
	after := ""
	posts := make([]reddit.PostResponseData, 0)
	for {
		response, err := h.RedditClient.FindMemes(request.Query, request.From, after, token.AccessToken.StringValue)
		if err != nil {
			return nil, errors.InternalError{
				Message: err.Error(),
			}
		}
		if response.StatusCode != http.StatusOK {
			return nil, errors.InternalError{
				Message: fmt.Sprintf("During fetching memes  request ended with %v code", response.StatusCode),
			}
		}
		resp, err := parseResponse(response)
		if err != nil {
			return nil, errors.InternalError{
				Message: fmt.Sprintf("Parsing response thrown error %v", err),
			}
		}
		for _, p := range resp.Data.Posts {
			posts = append(posts, p.Data)
		}
		if resp.Data.After != "null" && resp.Data.After != "" {
			after = resp.Data.After
		} else {
			break
		}
	}
	sort.SliceStable(posts, func(i, j int) bool {
		return posts[i].Score > posts[j].Score
	})
	return &response{Data: getPaginated(posts, request.Page, request.PageSize)}, nil
}

func getPaginated(response []reddit.PostResponseData, page, pageSize int64) []reddit.PostResponseData {
	min := int(page * pageSize)
	size := len(response)
	if min > size {
		return make([]reddit.PostResponseData, 0)
	}
	max := min + int(pageSize)
	if max > size {
		return response[min:]
	}
	return response[min:max]
}

func parseResponse(response *http.Response) (*reddit.SearchResponse, error) {
	result := reddit.SearchResponse{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func makeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return newMemeSearchHandler().getMemes(r.(request))
	}
}
