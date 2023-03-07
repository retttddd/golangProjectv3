package de_encoding

import (
	"crypto/sha256"
)

const keySize = 32 // length of secret key which is minimal to complete AES Encoding

func PassToSecretKey(pass string) []byte {

	h := sha256.New()
	h.Write([]byte(pass))

	encryptedKeyWord := h.Sum(nil)
	if h.Size() < 32 {
		dif := 32 - h.Size()
		return CreateArray(dif, encryptedKeyWord)
	}
	return encryptedKeyWord[:keySize]

}
func CreateArray(diff int, arr []byte) []byte {
	bigSlice := make([]byte, diff)
	for j := 0; j < len(bigSlice); j++ {
		bigSlice[j] = 65
	}
	arr = append(arr, bigSlice...)
	return arr
}
