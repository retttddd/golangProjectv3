package service

import (
	"awesomeProject3/service/ciphering"
)

type SimpleSecretService struct {
	storage    storage
	encoder    encoder
	keyEncoder encoder
}

func (ss *SimpleSecretService) ReadSecret(key string, password string) (string, error) {
	codedKey, err := ss.keyEncoder.Encrypt(key, ciphering.PassToSecretKey(password))
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
	return decryptedVal, nil
}

func (ss *SimpleSecretService) WriteSecret(key string, value string, password string) error {
	secretData, err := ss.encoder.Encrypt(value, ciphering.PassToSecretKey(password))
	if err != nil {
		return err
	}
	codedKey, err := ss.keyEncoder.Encrypt(key, ciphering.PassToSecretKey(password))
	if err != nil {
		return err
	}
	return ss.storage.Write(codedKey, secretData)
}

func New(st storage, en encoder, er encoder) *SimpleSecretService {
	//we create encoder and keyEncoder to separate encryption flow into key encryption(without random injection) and value encryption(with random part)
	return &SimpleSecretService{
		storage:    st,
		encoder:    en,
		keyEncoder: er,
	}
}
