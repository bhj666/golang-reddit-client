package persistance

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type TokenRepository interface {
	Save(Token)
	Delete(Token)
	FindActive(*Token, int64)
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

func (Token) TableName() string {
	return "token"
}

func NewTokenRepository() TokenRepository {
	return TokenRepositoryImpl{openConnection()}
}

func (repo TokenRepositoryImpl) Save(message Token) {
	repo.db.Create(message)
}

func (repo TokenRepositoryImpl) Delete(message Token) {
	repo.db.Delete(message)
}

func (repo TokenRepositoryImpl) FindActive(response *Token, time int64) {
	repo.db.Where("expires_at > ?", time).Find(response)
}

type TokenRepositoryImpl struct {
	db *gorm.DB
}
