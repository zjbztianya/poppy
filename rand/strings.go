package rand

import (
	"crypto/rand"
	"encoding/base64"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}

const RememberTokenBytes  = 32

func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
