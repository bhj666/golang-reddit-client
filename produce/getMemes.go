package produce

import (
	errors "aws-example/error"
	"aws-example/persistance"
	"aws-example/reddit"
	"aws-example/slack"
	"aws-example/timeprovider"
	"context"
	"fmt"
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
	Query     string
	Subreddit string
	From      string
	PageSize  int64
}

type response struct {
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
		searchResponse, err := h.RedditClient.FindMemes(request.Subreddit, request.Query, request.From, after, token.AccessToken.StringValue)
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
	memes := getPaginated(posts, request.PageSize)
	slackClient := slack.NewClient(*http.DefaultClient)
	slackClient.SendMessage(fmt.Sprintf("Memes from %s subreddit and from last %s for query %s", request.Subreddit, request.From, request.Query))
	for _, m := range memes {
		slackClient.SendMessage(m.Url)
	}
	return nil, nil
}

func getPaginated(response []reddit.PostResponseData, pageSize int64) []reddit.PostResponseData {
	size := len(response)
	max := int(pageSize)
	if max > size {
		return response
	}
	return response[0:max]
}

func makeEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		return newMemeSearchHandler().getMemes(r.(request))
	}
}
