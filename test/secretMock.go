package test

import "aws-example/persistance"

type SecretRepositoryMock struct {
	data map[string]bool
}

func GetSecretRepositoryMock() persistance.SecretsRepository {
	return &SecretRepositoryMock{make(map[string]bool)}
}

func (m *SecretRepositoryMock) Save(secret persistance.Secret) {
	m.data[secret.Secret] = true
}

func (m *SecretRepositoryMock) Delete(secret persistance.Secret) {
	m.data[secret.Secret] = false
}

func (m *SecretRepositoryMock) Find(key string, secret *persistance.Secret) {
	if m.data[key] {
		secret.Secret = key
	}
}
