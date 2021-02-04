package test

import (
	"net/http"
	"strings"
)

type RedditClientMock struct {
	ExchangeStatus int
	ExchangeBody   string
	ExchangeError  error
	RefreshStatus  int
	RefreshBody    string
	RefreshError   error
	FindResults    map[string]FindResult
}

type FindResult struct {
	FindStatus int
	FindBody   string
	FindError  error
}

type bodyReader struct {
	reader *strings.Reader
}

func (r *bodyReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *bodyReader) Close() error {
	return nil
}

func (*RedditClientMock) GetRedirectUrl(string) string {
	return "REDIRECT_URL"
}

func (c *RedditClientMock) ExchangeForToken(string) (*http.Response, error) {
	return buildResponse(c.ExchangeStatus, c.ExchangeBody), c.ExchangeError
}

func (c *RedditClientMock) RefreshToken(string) (*http.Response, error) {
	return buildResponse(c.RefreshStatus, c.RefreshBody), c.RefreshError
}

func (c *RedditClientMock) FindMemes(query string, from string, after string, token string) (*http.Response, error) {
	response, ok := c.FindResults[after]
	if !ok {
		return buildResponse(http.StatusNotFound, "{}"), nil
	}
	return buildResponse(response.FindStatus, response.FindBody), response.FindError
}

func buildResponse(code int, body string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Body: &bodyReader{
			reader: strings.NewReader(body),
		},
	}
}
