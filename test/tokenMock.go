package test

import (
	"aws-example/persistance"
	"github.com/pkg/errors"
)

type TokenRepositoryMock struct {
	tokens []persistance.Token
}

func GetTokenRepositoryMock() persistance.TokenRepository {
	return &TokenRepositoryMock{make([]persistance.Token, 0)}
}

func (m *TokenRepositoryMock) Save(token persistance.Token) error {
	shouldAppend := true
	for i := range m.tokens {
		if m.tokens[i].AccessToken == token.AccessToken {
			shouldAppend = false
			break
		}
	}
	if shouldAppend {
		m.tokens = append(m.tokens, token)
	}
	return nil
}

func (m *TokenRepositoryMock) Delete(token persistance.Token) error {
	index := -1
	for i := range m.tokens {
		if m.tokens[i].AccessToken == token.AccessToken {
			index = i
			break
		}
	}
	if index >= 0 {
		m.tokens = append(m.tokens[:index], m.tokens[index+1:]...)
		return nil
	}
	return errors.New("Not found")
}

func (m *TokenRepositoryMock) FindActive(result *persistance.Token, time int64) error {
	for _, v := range m.tokens {
		if v.ExpiresAt > time {
			result.AccessToken = v.AccessToken
			result.RefreshToken = v.RefreshToken
			result.ExpiresAt = v.ExpiresAt
			return nil
		}
	}
	return errors.New("Not found")
}
