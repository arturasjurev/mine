package manager

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateCorrelationID() (string, error) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
