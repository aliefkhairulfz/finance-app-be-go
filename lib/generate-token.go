package lib

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSessionToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	return token[:length], nil
}

func GenerateCsrfToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(bytes)
	return token[:length], nil
}
