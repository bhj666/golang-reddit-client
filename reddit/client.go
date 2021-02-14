package reddit

import (
	"aws-example/config"
	errors "aws-example/error"
	"aws-example/persistance"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	GetRedirectUrl(secret string) string
	ExchangeForToken(code string) (*persistance.Token, error)
	RefreshToken(refreshToken string) (*persistance.Token, error)
	FindMemes(subreddit, query, from, after, token string) (*SearchResponse, error)
}

type ClientImpl struct {
	Client http.Client
}

func NewClient(client http.Client) Client {
	return ClientImpl{Client: client}
}

func (ClientImpl) GetRedirectUrl(secret string) string {
	result := fmt.Sprintf(config.REDDIT_AUTHORIZE_URL_TEMPLATE, config.REDDIT_APP_ID, secret, config.REDDIT_REDIRECT_URL, config.REDDIT_SCOPE)
	log.Printf("Responding with url %s", result)
	return result

}

func (c ClientImpl) ExchangeForToken(code string) (*persistance.Token, error) {
	values := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {config.REDDIT_REDIRECT_URL}}
	data := strings.NewReader(values.Encode())
	client := c.Client
	request, err := http.NewRequest("POST", config.REDDIT_TOKEN_EXCHANGE_URL, data)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(config.REDDIT_APP_ID, config.REDDIT_APP_SECRET)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	enrichRequest(request)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	response, er := ioutil.ReadAll(resp.Body)
	if er != nil {
		return nil, errors.InternalError{Message: er.Error()}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.InternalError{Message: fmt.Sprintf("Response code: %v", resp.StatusCode)}
	}
	var token persistance.Token
	er = json.Unmarshal(response, &token)
	if er != nil {
		return nil, errors.InternalError{Message: er.Error()}
	}
	return &token, nil
}

func (c ClientImpl) RefreshToken(refreshToken string) (*persistance.Token, error) {
	values := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken}}
	data := strings.NewReader(values.Encode())
	client := c.Client
	request, err := http.NewRequest("POST", config.REDDIT_TOKEN_EXCHANGE_URL, data)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(config.REDDIT_APP_ID, config.REDDIT_APP_SECRET)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	enrichRequest(request)
	resp, er := client.Do(request)
	if er != nil {
		log.Print("Request for refresh token failed")
		return nil, er
	}
	if resp.StatusCode != http.StatusOK {
		log.Print("No active tokens to refresh")
		return nil, errors.GenericResponseError{Message: "Request for refresh token failed", ResponseCode: http.StatusInternalServerError}
	}
	var newToken persistance.Token
	response, er := ioutil.ReadAll(resp.Body)
	er = json.Unmarshal(response, &newToken)
	if er != nil {
		log.Errorf("Error when parsing token refresh response %v", er)
		return nil, errors.GenericResponseError{Message: "Parsing token failed", ResponseCode: http.StatusInternalServerError}
	}
	return &newToken, nil
}

func (c ClientImpl) FindMemes(subreddit, query, from, after, token string) (*SearchResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(config.REDDIT_SEARCH_URL_TEMPLATE, subreddit), nil)
	if err != nil {
		return nil, err
	}
	client := c.Client
	q := req.URL.Query()
	q.Add("q", query)
	q.Add("restrict_sr", "1")
	q.Add("t", from)
	q.Add("after", after)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+token)
	fmt.Println(req.URL.String())
	enrichRequest(req)
	body, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return parseResponse(body)
}

func parseResponse(response *http.Response) (*SearchResponse, error) {
	if response.StatusCode != http.StatusOK {
		return nil, errors.InternalError{
			Message: fmt.Sprintf("During fetching memes  request ended with %v code", response.StatusCode),
		}
	}
	result := SearchResponse{}
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

func enrichRequest(request *http.Request) {
	request.Header.Set("User-agent", config.REDDIT_USER_AGENT)
}
