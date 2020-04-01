package rand

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/zjbztianya/poppy/conf"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}

func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}

func RememberToken() (string, error) {
	return String(conf.Conf.App.RememberTokenBytes)
}
