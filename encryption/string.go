package encryption

import (
	"aws-example/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"reflect"
)

type EncryptedString struct {
	StringValue string
}

func (es *EncryptedString) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		decrypted, err := Decrypt(v)
		if err != nil {
			return err
		}
		es.StringValue = string(decrypted)
		return nil
	case []byte:
		encrypted, err := Decrypt(string(v))
		if err != nil {
			return err
		}
		es.StringValue = string(encrypted)
		return nil
	default:
		return fmt.Errorf("failed to scan type %+v for value", reflect.TypeOf(value))
	}
}

func (es EncryptedString) Value() (driver.Value, error) {
	encrypted, err := Encrypt(es.StringValue)
	if err != nil {
		return "", err
	}
	return encrypted, nil
}

func (es EncryptedString) GormDataType() string {
	return "string"
}

func (es *EncryptedString) MarshalJSON() ([]byte, error) {
	return json.Marshal(es.StringValue)
}

func (es *EncryptedString) UnmarshalJSON(b []byte) error {
	var target string
	if err := json.Unmarshal(b, &target); err != nil {
		return err
	}
	es.StringValue = target
	return nil
}

func Encrypt(plaintext string) ([]byte, error) {
	c, err := aes.NewCipher([]byte(config.ENCRYPTION_SALT))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, []byte(plaintext), nil), nil
}

func Decrypt(ciphertext string) ([]byte, error) {
	c, err := aes.NewCipher([]byte(config.ENCRYPTION_SALT))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
}
