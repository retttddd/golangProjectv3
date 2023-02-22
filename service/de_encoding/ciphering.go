package de_encoding

import (
	"crypto/aes"
	"encoding/hex"
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

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	return string(pt[:])
}

func (en AESEncoder) Encrypt(plaintext string) string {
	c, err := aes.NewCipher(en.key)
	checkError(err)

	out := make([]byte, len(plaintext))

	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
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
