package service

import "awesomeProject3/storage"

type SimpleSecretService struct {
	storage storage.Storage
}

func (ss *SimpleSecretService) ReadSecret(key []byte) (string, error) {
	return ss.storage.Read(key)
}
func (as *SimpleSecretService) WriteSecret(key []byte, value string) {
	as.storage.Write(key, value)
}
func New(st storage.Storage) *SimpleSecretService {
	return &SimpleSecretService{
		storage: st,
	}

}
