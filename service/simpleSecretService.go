package service

import (
	"awesomeProject3/service/ciphering"
	"awesomeProject3/storage"
)

type SimpleSecretService struct {
	storage storage.Storage
	encoder ciphering.Encoder
}

func (ss *SimpleSecretService) ReadSecret(key string, password string) (string, error) {
	encryptedVal, err := ss.storage.Read(key)
	if err != nil {
		return "", err
	}
	x, err2 := ss.encoder.Decrypt(encryptedVal, ciphering.PassToSecretKey(password))
	if err2 != nil {
		return "", err2
	}
	return string(x), nil
}

func (ss *SimpleSecretService) WriteSecret(key string, value string, password string) error {
	x, err := ss.encoder.Encrypt(value, ciphering.PassToSecretKey(password))
	if err != nil {
		return err
	}
	return ss.storage.Write(key, x)
}

func New(st storage.Storage, en ciphering.Encoder) SimpleSecretService {

	return SimpleSecretService{
		storage: st,
		encoder: en,
	}
}
