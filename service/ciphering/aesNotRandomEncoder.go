package ciphering

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

type regularEncoder struct {
}

func (er regularEncoder) Decrypt(ct string, aesKey []byte) (string, error) {
	ciphertext, err := hex.DecodeString(ct)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext len is too short")
	}
	plaintext, err := gcm.Open(nil, make([]byte, nonceSize), ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext[:]), nil
}

func (er regularEncoder) Encrypt(plaintext string, aesKey []byte) (encryptedText string, err error) {
	c, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err

	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err

	}
	nonce := make([]byte, gcm.NonceSize())
	encryptedText = hex.EncodeToString(gcm.Seal(nil, nonce, []byte(plaintext), nil))
	return encryptedText, nil
}

func NewRegularEncoder() *regularEncoder {

	return &regularEncoder{}
}
