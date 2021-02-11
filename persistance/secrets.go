package persistance

import (
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type SecretsRepository interface {
	Save(Secret) error
	Delete(Secret) error
	Find(string, *Secret) error
}

type Secret struct {
	Secret string `json:"secret" validate:"required"`
}

func (Secret) TableName() string {
	return "secrets"
}

func (m *Secret) Validate() error {
	validate := validator.New()
	return validate.Struct(m)
}

func NewSecretsRepository() SecretsRepository {
	return SecretsRepositoryImpl{openConnection()}
}

func (repo SecretsRepositoryImpl) Save(message Secret) error {
	return repo.db.Create(message).Error
}

func (repo SecretsRepositoryImpl) Delete(message Secret) error {
	return repo.db.Delete(message).Error
}

func (repo SecretsRepositoryImpl) Find(secret string, response *Secret) error {
	return repo.db.Where("secret = ?", secret).Find(response).Error
}

type SecretsRepositoryImpl struct {
	db *gorm.DB
}
