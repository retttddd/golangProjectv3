package ciphering

import (
	"crypto/sha256"
)

const keySize = 32 // length of secret key which is minimal to complete AES Encoding

func PassToSecretKey(pass string) []byte {

	h := sha256.New()
	h.Write([]byte(pass))

	encryptedKeyWord := h.Sum(nil)
	if h.Size() < keySize {
		panic("key size is incorrect ")
	}
	return encryptedKeyWord[:keySize]

}
