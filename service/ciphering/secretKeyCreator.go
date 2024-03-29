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
		panic("can not covert password into AES key. Incorrect length")
	}
	return encryptedKeyWord[:keySize]

}
