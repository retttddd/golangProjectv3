package service

import (
	"awesomeProject3/service/de_encoding"
	"awesomeProject3/storage"
)

const trash = "1234567890123456"

type SimpleSecretService struct {
	storage storage.Storage
	encoder de_encoding.Encoder
}

func (ss *SimpleSecretService) ReadSecret(key string) (string, error) {
	encryptedVal, err := ss.storage.Read(key)
	checkError(err)
	return ss.encoder.Decrypt(encryptedVal), nil
}

func (ss *SimpleSecretService) WriteSecret(key string, value string) {
	ss.storage.Write(key, ss.encoder.Encrypt(value))
}

func New(st storage.Storage, en de_encoding.Encoder) *SimpleSecretService {

	return &SimpleSecretService{
		storage: st,
		encoder: en,
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
