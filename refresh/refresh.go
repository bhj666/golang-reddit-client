package main

import (
	"aws-example/persistance"
	"aws-example/reddit"
	"aws-example/timeprovider"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type refreshHandler interface {
	refreshToken() func()
}

type refreshHandlerImpl struct {
	TokenRepository persistance.TokenRepository
	RedditClient    reddit.Client
	TimeProvider    timeprovider.Provider
}

func newRefreshHandler() refreshHandler {
	return refreshHandlerImpl{
		persistance.NewTokenRepository(),
		reddit.NewClient(),
		timeprovider.ProviderImpl{},
	}
}

func (h refreshHandlerImpl) refreshToken() func() {
	return func() {
		log.Printf("Scheduled function called")
		db := h.TokenRepository
		token := persistance.Token{}
		db.FindActive(&token, h.TimeProvider.GetCurrentSeconds())
		if token.AccessToken == "" {
			return
		}
		resp, er := h.RedditClient.RefreshToken(token.RefreshToken)
		if er != nil || resp.StatusCode != http.StatusOK {
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
