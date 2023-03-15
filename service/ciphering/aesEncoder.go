package ciphering

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

type aesEncoder struct {
}

func (en aesEncoder) Decrypt(ct string, aesKey []byte) (plaintext []byte, err error) {
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
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext[:], nil
}

func (en aesEncoder) Encrypt(plaintext string, aesKey []byte) (encryptedText string, err error) {
	c, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err

	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err

	}
	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	encryptedText = hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(plaintext), nil))
	return encryptedText, nil

}

func NewAESEncoder() aesEncoder {

	return aesEncoder{}
}
