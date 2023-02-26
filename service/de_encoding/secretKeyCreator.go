package de_encoding

import (
	"crypto/sha256"
)

func PassToSecretKey(string2 string) []byte {
	h := sha256.New()
	h.Write([]byte(string2))
	encryptedKeyWord := (h.Sum(nil))
	return encryptedKeyWord[:32]

}
