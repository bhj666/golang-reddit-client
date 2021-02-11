package memes

import (
	errors "aws-example/error"
	"aws-example/persistance"
	"aws-example/reddit"
	"aws-example/timeprovider"
	"context"
	"github.com/go-kit/kit/endpoint"
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
		reddit.NewClient(*http.DefaultClient),
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
	err := db.FindActive(&token, h.TimeProvider.GetCurrentSeconds())
	if err != nil || token.AccessToken.StringValue == "" {
		return nil, errors.UnauthorizedError{}
	}
	after := ""
	posts := make([]reddit.PostResponseData, 0)
	for {
		searchResponse, err := h.RedditClient.FindMemes(request.Query, request.From, after, token.AccessToken.StringValue)
		if err != nil {
			return nil, errors.InternalError{
				Message: err.Error(),
			}
		}
		for _, p := range searchResponse.Data.Posts {
			posts = append(posts, p.Data)
		}
		if searchResponse.Data.After != "null" && searchResponse.Data.After != "" {
			after = searchResponse.Data.After
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

func makeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return newMemeSearchHandler().getMemes(r.(request))
	}
}
