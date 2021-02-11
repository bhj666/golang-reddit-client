package persistance

import (
	"aws-example/encryption"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type TokenRepository interface {
	Save(Token) error
	Delete(Token) error
	FindActive(*Token, int64) error
}

type Token struct {
	AccessToken  encryption.EncryptedString `json:"access_token"`
	TokenType    string                     `json:"token_type"`
	ExpiresIn    int64                      `json:"expires_in"`
	Scope        string                     `json:"scope"`
	RefreshToken encryption.EncryptedString `json:"refresh_token"`
	ExpiresAt    int64                      `json:"expires_at"`
}

func (Token) TableName() string {
	return "token"
}

func NewTokenRepository() TokenRepository {
	return TokenRepositoryImpl{openConnection()}
}

func (repo TokenRepositoryImpl) Save(message Token) error {
	return repo.db.Create(message).Error
}

func (repo TokenRepositoryImpl) Delete(message Token) error {
	return repo.db.Delete(message).Error
}

func (repo TokenRepositoryImpl) FindActive(response *Token, time int64) error {
	return repo.db.Where("expires_at > ?", time).Find(response).Error
}

type TokenRepositoryImpl struct {
	db *gorm.DB
}
