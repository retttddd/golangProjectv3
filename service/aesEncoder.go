package service

import (
	"awesomeProject3/service/de_encoding"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type aesEncoder struct {
	key []byte
}

func (en aesEncoder) Decrypt(ct string) string {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(en.key)
	checkError(err)
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		panic(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext[:])
}

func (en aesEncoder) Encrypt(plaintext string) string {
	c, err := aes.NewCipher(en.key)
	checkError(err)
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)

	}
	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}
	return hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(plaintext), nil))
}

func NewAESEncoder(aesKey []byte) de_encoding.Encoder {

	return &aesEncoder{
		key: aesKey,
	}
}
