package service

import (
	"awesomeProject3/service/ciphering"
	"awesomeProject3/storage"
)

type SimpleSecretService struct {
	storage          storage.Storage
	encoder          ciphering.Encoder
	notrandomencoder ciphering.Encoder
}

func (ss *SimpleSecretService) ReadSecret(key string, password string) (string, error) {
	codedKey, err := ss.notrandomencoder.Encrypt(key, ciphering.PassToSecretKey(password)) //encode key without random aspect
	if err != nil {
		return "", err
	}
	encryptedVal, err := ss.storage.Read(codedKey)
	if err != nil {
		return "", err
	}
	decryptedVal, err := ss.encoder.Decrypt(encryptedVal, ciphering.PassToSecretKey(password))
	if err != nil {
		return "", err
	}
	return string(decryptedVal), nil
}

func (ss *SimpleSecretService) WriteSecret(key string, value string, password string) error {
	secretData, err := ss.encoder.Encrypt(value, ciphering.PassToSecretKey(password)) //encode value(standart way)
	if err != nil {
		return err
	}
	codedKey, err := ss.notrandomencoder.Encrypt(key, ciphering.PassToSecretKey(password)) //encode key without random aspect
	if err != nil {
		return err
	}
	return ss.storage.Write(codedKey, secretData)
}

func New(st storage.Storage, en ciphering.Encoder, er ciphering.Encoder) *SimpleSecretService {

	return &SimpleSecretService{
		storage:          st,
		encoder:          en,
		notrandomencoder: er,
	}
}
