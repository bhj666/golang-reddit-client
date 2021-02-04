package persistance

import (
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type SecretsRepository interface {
	Save(Secret)
	Delete(Secret)
	Find(string, *Secret)
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

func (repo SecretsRepositoryImpl) Save(message Secret) {
	repo.db.Create(message)
}

func (repo SecretsRepositoryImpl) Delete(message Secret) {
	repo.db.Delete(message)
}

func (repo SecretsRepositoryImpl) Find(secret string, response *Secret) {
	repo.db.Where("secret = ?", secret).Find(response)
}

type SecretsRepositoryImpl struct {
	db *gorm.DB
}
