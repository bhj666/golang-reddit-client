package reddit

import (
	"aws-example/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	GetRedirectUrl(secret string) string
	ExchangeForToken(code string) (*http.Response, error)
	RefreshToken(refreshToken string) (*http.Response, error)
	FindMemes(query, from, after, token string) (*http.Response, error)
}

type ClientImpl struct {
}

func NewClient() Client {
	return ClientImpl{}
}

func (ClientImpl) GetRedirectUrl(secret string) string {
	result := fmt.Sprintf(config.REDDIT_AUTHORIZE_URL_TEMPLATE, config.REDDIT_APP_ID, secret, config.REDDIT_REDIRECT_URL, config.REDDIT_SCOPE)
	log.Printf("Responding with url %s", result)
	return result

}

func (c ClientImpl) ExchangeForToken(code string) (*http.Response, error) {
	values := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {config.REDDIT_REDIRECT_URL}}
	data := strings.NewReader(values.Encode())
	client := http.DefaultClient
	request, err := http.NewRequest("POST", config.REDDIT_TOKEN_EXCHANGE_URL, data)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(config.REDDIT_APP_ID, config.REDDIT_APP_SECRET)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	enrichRequest(request)
	return client.Do(request)
}

func (c ClientImpl) RefreshToken(refreshToken string) (*http.Response, error) {
	values := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken}}
	data := strings.NewReader(values.Encode())
	client := http.DefaultClient
	request, err := http.NewRequest("POST", config.REDDIT_TOKEN_EXCHANGE_URL, data)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(config.REDDIT_APP_ID, config.REDDIT_APP_SECRET)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	enrichRequest(request)
	return client.Do(request)
}

func (c ClientImpl) FindMemes(query, from, after, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", config.REDDIT_SEARCH_URL, nil)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	q := req.URL.Query()
	q.Add("q", query)
	q.Add("restrict_sr", "1")
	q.Add("from", from)
	q.Add("after", after)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+token)
	fmt.Println(req.URL.String())
	enrichRequest(req)
	return client.Do(req)
}

func enrichRequest(request *http.Request) {
	request.Header.Set("User-agent", config.REDDIT_USER_AGENT)
}
