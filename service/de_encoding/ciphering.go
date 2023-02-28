package de_encoding

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type Encoder interface {
	Encrypt(plaintext string) string
	Decrypt(ct string) string
}

type AESEncoder struct {
	key []byte
}

func (en AESEncoder) Decrypt(ct string) string {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(en.key)
	checkError(err)
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println(err)
	}

	return string(plaintext[:])
}

func (en AESEncoder) Encrypt(plaintext string) string {
	c, err := aes.NewCipher(en.key)
	checkError(err)
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(plaintext), nil))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func NewAESEncoder(aesKey []byte) Encoder {

	return &AESEncoder{
		key: aesKey,
	}
}
