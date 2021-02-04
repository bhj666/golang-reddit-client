package main

import (
	"aws-example/persistance"
	"aws-example/reddit"
	"aws-example/timeprovider"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

type MemeSearchHandler struct {
	TokenRepository persistance.TokenRepository
	RedditClient    reddit.Client
	TimeProvider    timeprovider.Provider
}

func NewMemeSearchHandler() http.Handler {
	return MemeSearchHandler{
		persistance.NewTokenRepository(),
		reddit.NewClient(),
		timeprovider.ProviderImpl{},
	}
}

type unauthorizedError struct {
}

func (unauthorizedError) Error() string {
	return "You need to authorize first"
}

func (h MemeSearchHandler) ServeHTTP(responseWriter http.ResponseWriter,
	request *http.Request) {
	log.Print("Get memes called")
	query := request.URL.Query()
	memeQuery := query.Get("query")
	from := query.Get("from")
	page, err := strconv.ParseInt(query.Get("page"), 0, 64)
	if err != nil {
		page = 0
	}
	pageSize, err := strconv.ParseInt(query.Get("pageSize"), 0, 64)
	if err != nil {
		pageSize = 25
	}
	body, code, err := HandleGetMemes(h, memeQuery, from, page, pageSize)
	if handleError(responseWriter, code, err) {
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(body))
}

func HandleGetMemes(h MemeSearchHandler, query, from string, page, pageSize int64) (string, int, error) {
	db := h.TokenRepository
	token := persistance.Token{}
	db.FindActive(&token, h.TimeProvider.GetCurrentSeconds())
	if token.AccessToken == "" {
		return "", http.StatusUnauthorized, unauthorizedError{}
	}
	after := ""
	posts := make([]reddit.PostResponseData, 0)
	for {
		response, err := h.RedditClient.FindMemes(query, from, after, token.AccessToken)
		if err != nil {
			return "", http.StatusInternalServerError, err
		}
		if response.StatusCode != http.StatusOK {
			return "Search return code", response.StatusCode, nil
		}
		resp := parseResponse(response)
		if resp == nil || len(resp.Data.Posts) == 0 {
			return "", 500, err
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
	body, err := json.Marshal(getPaginated(posts, page, pageSize))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	return string(body), http.StatusOK, nil
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

func parseResponse(response *http.Response) *reddit.SearchResponse {
	result := reddit.SearchResponse{}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil
	}
	return &result
}
