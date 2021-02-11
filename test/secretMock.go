package test

import (
	"aws-example/persistance"
	"github.com/pkg/errors"
)

type SecretRepositoryMock struct {
	data map[string]bool
}

func GetSecretRepositoryMock() persistance.SecretsRepository {
	return &SecretRepositoryMock{make(map[string]bool)}
}

func (m *SecretRepositoryMock) Save(secret persistance.Secret) error {
	m.data[secret.Secret] = true
	return nil
}

func (m *SecretRepositoryMock) Delete(secret persistance.Secret) error {
	m.data[secret.Secret] = false
	return nil
}

func (m *SecretRepositoryMock) Find(key string, secret *persistance.Secret) error {
	if m.data[key] {
		secret.Secret = key
		return nil
	}
	return errors.New("Not found")
}
