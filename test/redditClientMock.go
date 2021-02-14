package test

import (
	"aws-example/persistance"
	"aws-example/reddit"
	"github.com/pkg/errors"
)

type RedditClientMock struct {
	ExchangeStatus int
	ExchangeBody   *persistance.Token
	ExchangeError  error
	RefreshBody    *persistance.Token
	RefreshError   error
	FindResults    map[string]FindResult
}

type FindResult struct {
	FindBody reddit.SearchResponse
}

func (*RedditClientMock) GetRedirectUrl(string) string {
	return "REDIRECT_URL"
}

func (c *RedditClientMock) ExchangeForToken(string) (*persistance.Token, error) {
	return c.ExchangeBody, c.ExchangeError
}

func (c *RedditClientMock) RefreshToken(string) (*persistance.Token, error) {
	return c.RefreshBody, c.RefreshError
}

func (c *RedditClientMock) FindMemes(subreddit, query, from, after, token string) (*reddit.SearchResponse, error) {
	response, ok := c.FindResults[after]
	if !ok {
		return nil, errors.New("Not found")
	}
	return &response.FindBody, nil
}
