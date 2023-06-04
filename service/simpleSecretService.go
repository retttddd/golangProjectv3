package service

import (
	"awesomeProject3/service/ciphering"
)

type SimpleSecretService struct {
	storage    storage
	encoder    encoder
	keyEncoder encoder
}

func (ss *SimpleSecretService) ReadSecret(key string, password string) (*SecretServiceModel, error) {
	codedKey, err := ss.keyEncoder.Encrypt(key, ciphering.PassToSecretKey(password))
	if err != nil {
		return nil, err
	}
	encryptedData, err := ss.storage.Read(codedKey)
	if err != nil {
		return nil, err
	}
	decryptedVal, err := ss.encoder.Decrypt(*encryptedData.Value, ciphering.PassToSecretKey(password))
	if err != nil {
		return nil, err
	}
	return &SecretServiceModel{
		Value: &decryptedVal,
	}, nil
}

func (ss *SimpleSecretService) WriteSecret(key string, model *SecretServiceModel, password string) error {
	secretData, err := ss.encoder.Encrypt(*model.Value, ciphering.PassToSecretKey(password))
	if err != nil {
		return err
	}
	codedKey, err := ss.keyEncoder.Encrypt(key, ciphering.PassToSecretKey(password))
	if err != nil {
		return err
	}
	return ss.storage.Write(codedKey, &StorageModel{Value: &secretData})
}

func New(st storage, en encoder, er encoder) *SimpleSecretService {
	//we create encoder and keyEncoder to separate encryption flow into key encryption(without random injection)
	//and value encryption(with random part)
	return &SimpleSecretService{
		storage:    st,
		encoder:    en,
		keyEncoder: er,
	}
}
