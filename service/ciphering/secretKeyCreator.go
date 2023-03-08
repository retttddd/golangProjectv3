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
		dif := keySize - h.Size()
		return createArray(dif, encryptedKeyWord)
	}
	return encryptedKeyWord[:keySize]

}
func createArray(diff int, arr []byte) []byte {
	bigSlice := make([]byte, diff)
	for j := 0; j < len(bigSlice); j++ {
		bigSlice[j] = 65
	}
	arr = append(arr, bigSlice...)
	return arr
}
