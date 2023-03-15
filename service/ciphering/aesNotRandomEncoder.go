package ciphering

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

type aesNotRandomEncoder struct {
}

func (er aesNotRandomEncoder) Decrypt(ct string, aesKey []byte) (plaintext []byte, err error) {
	ciphertext, _ := hex.DecodeString(ct)
	c, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext len is too short")
	}
	plaintext, err = gcm.Open(nil, make([]byte, nonceSize), ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext[:], nil
}

func (er aesNotRandomEncoder) Encrypt(plaintext string, aesKey []byte) (encryptedText string, err error) {
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

func NewAESNotRandomEncoder() aesNotRandomEncoder {

	return aesNotRandomEncoder{}
}
