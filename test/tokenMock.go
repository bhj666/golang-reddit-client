package test

import "aws-example/persistance"

type TokenRepositoryMock struct {
	tokens []persistance.Token
}

func GetTokenRepositoryMock() persistance.TokenRepository {
	return &TokenRepositoryMock{make([]persistance.Token, 0)}
}

func (m *TokenRepositoryMock) Save(token persistance.Token) {
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
}

func (m *TokenRepositoryMock) Delete(token persistance.Token) {
	index := -1
	for i := range m.tokens {
		if m.tokens[i].AccessToken == token.AccessToken {
			index = i
			break
		}
	}
	if index >= 0 {
		m.tokens = append(m.tokens[:index], m.tokens[index+1:]...)
	}
}

func (m *TokenRepositoryMock) FindActive(result *persistance.Token, time int64) {
	for _, v := range m.tokens {
		if v.ExpiresAt > time {
			result.AccessToken = v.AccessToken
			result.RefreshToken = v.RefreshToken
			result.ExpiresAt = v.ExpiresAt
			break
		}
	}
}
