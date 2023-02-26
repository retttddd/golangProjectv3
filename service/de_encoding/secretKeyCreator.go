package de_encoding

import (
	"crypto/sha256"
)

func PassToSecretKey(string2 string) []byte {
	return shaEncrypt(string2)
}

func shaEncrypt(fileNameFunc string) []byte {
	h := sha256.New()
	h.Write([]byte(fileNameFunc))
	encryptedKeyWord := h.Sum(nil)
	return encryptedKeyWord
}
