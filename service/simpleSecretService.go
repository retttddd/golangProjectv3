package service

import "awesomeProject3/storage"

type SimpleSecretService struct {
	storage storage.Storage
}

func (ss *SimpleSecretService) ReadSecret(key string) (string, error) {
	return ss.storage.Read(key)

}
func New(st storage.Storage) SecretService {
	return &SimpleSecretService{
		storage: st,
	}

}
