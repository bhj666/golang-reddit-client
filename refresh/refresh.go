package main

import (
	"aws-example/persistance"
	"aws-example/reddit"
	"aws-example/timeprovider"
	log "github.com/sirupsen/logrus"
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
		reddit.NewClient(*http.DefaultClient),
		timeprovider.ProviderImpl{},
	}
}

func (h refreshHandlerImpl) refreshToken() func() {
	return func() {
		log.Printf("Scheduled function called")
		db := h.TokenRepository
		token := persistance.Token{}
		err := db.FindActive(&token, h.TimeProvider.GetCurrentSeconds())
		if err != nil || token.AccessToken.StringValue == "" {
			return
		}
		newToken, er := h.RedditClient.RefreshToken(token.RefreshToken.StringValue)
		if er != nil {
			return
		}
		newToken.ExpiresAt = h.TimeProvider.GetCurrentSeconds() + newToken.ExpiresIn
		newToken.RefreshToken = token.RefreshToken
		err = db.Delete(token)
		if err != nil {
			log.Error(err)
		}
		err = db.Save(*newToken)
		if err != nil {
			log.Error(err)
		}
		log.Print("Successfully refreshed token")
	}
}
