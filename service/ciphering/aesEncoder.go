package ciphering

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

type aesEncoder struct {
	np nonceProducer
}
type nonceProducer interface {
	generate(size int) (string, error)
}

func (en aesEncoder) Decrypt(ct string, aesKey []byte) (string, error) {
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
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext[:]), nil
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

	nonce, err := en.np.generate(gcm.NonceSize())

	if err != nil {
		return "", err
	}
	encryptedText = hex.EncodeToString(gcm.Seal([]byte(nonce), []byte(nonce), []byte(plaintext), nil))
	return encryptedText, nil

}

func NewAESEncoder(np nonceProducer) *aesEncoder {

	return &aesEncoder{
		np,
	}
}
