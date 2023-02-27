package service

import (
	"awesomeProject3/service/de_encoding"
	"awesomeProject3/storage"
	"strings"
)

const trash = "1234567890123456"

type SimpleSecretService struct {
	storage storage.Storage
	encoder de_encoding.Encoder
}

func (ss *SimpleSecretService) ReadSecret(key string) (string, error) {
	encryptedVal, err := ss.storage.Read(key)
	checkError(err)
	return ss.encoder.Decrypt(gettingRidOfTrash(encryptedVal)), nil
}

func (ss *SimpleSecretService) WriteSecret(key string, value string) {
	ss.storage.Write(key, ss.encoder.Encrypt(addingTrash(value)))
}
func addingTrash(val string) string {
	return val + trash
}
func gettingRidOfTrash(val string) string {
	return strings.TrimRight(val, trash)
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
