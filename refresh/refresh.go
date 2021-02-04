package main

import (
	"aws-example/persistance"
	"aws-example/reddit"
	"aws-example/timeprovider"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

type RefreshHandler interface {
	RefreshToken() func()
}

type RefreshHandlerImpl struct {
	TokenRepository persistance.TokenRepository
	RedditClient    reddit.Client
	TimeProvider    timeprovider.Provider
}

func NewRefreshHandler() RefreshHandler {
	return RefreshHandlerImpl{
		persistance.NewTokenRepository(),
		reddit.NewClient(),
		timeprovider.ProviderImpl{},
	}
}

func (h RefreshHandlerImpl) RefreshToken() func() {
	return func() {
		log.Printf("Scheduled function called")
		db := h.TokenRepository
		token := persistance.Token{}
		db.FindActive(&token, h.TimeProvider.GetCurrentSeconds())
		if token.AccessToken == "" {
			return
		}
		resp, er := h.RedditClient.RefreshToken(token.RefreshToken)
		if er != nil || resp.StatusCode != 200 {
			log.Print("No active tokens to refresh")
			return
		}
		var newToken persistance.Token
		response, er := ioutil.ReadAll(resp.Body)
		er = json.Unmarshal(response, &newToken)
		if er != nil {
			log.Errorf("Error when parsing token refresh response %v", er)
			return
		}
		newToken.ExpiresAt = h.TimeProvider.GetCurrentSeconds() + newToken.ExpiresIn
		newToken.RefreshToken = token.RefreshToken
		db.Delete(token)
		db.Save(newToken)
		log.Print("Successfully refreshed token")
	}
}
