package de_encoding

import (
	"crypto/sha256"
)

func PassToSecretKey(pass string) []byte {
	h := sha256.New()
	h.Write([]byte(pass))
	encryptedKeyWord := h.Sum(nil)
	return encryptedKeyWord[:32]

}
