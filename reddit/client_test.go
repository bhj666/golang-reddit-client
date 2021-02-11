package reddit

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func getClientForTokenSuccess(t *testing.T) Client {
	client := getSuccessClient("{\"access_token\":\"a_token\", \"token_type\":\"bearer\", \"expires_in\":3600, \"refresh_token\":\"r_token\"}", t)
	return ClientImpl{Client: *client}
}

func getSuccessClient(body string, t *testing.T) *http.Client {
	client := NewTestClient(func(req *http.Request) *http.Response {
		require.Equal(t, "POST", req.Method)
		reqBody, _ := ioutil.ReadAll(req.Body)
		require.Contains(t, string(reqBody), "code=anyCode&grant_type=authorization_code")
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header:     make(http.Header),
		}
	})
	return client
}

func getUnauthorizedClient() *http.Client {
	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 403,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
			Header:     make(http.Header),
		}
	})
	return client
}

func TestParsingTokenCorrectResponse(t *testing.T) {
	client := getClientForTokenSuccess(t)

	token, err := client.ExchangeForToken("anyCode")

	require.Nil(t, err, "There should be no error")
	require.Equal(t, "a_token", token.AccessToken.StringValue)
}

func TestParsingTokenErrorResponse(t *testing.T) {
	client := ClientImpl{Client: *getUnauthorizedClient()}

	token, err := client.ExchangeForToken("anyCode")

	require.NotNil(t, err, "There should be some error")
	require.Nil(t, token)
}
