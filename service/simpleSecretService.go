package service

import (
	"awesomeProject3/storage"
	"crypto/aes"
	"encoding/hex"
)

const secretkey = "thisis32bitlongpassphraseimusing"

type SimpleSecretService struct {
	storage storage.Storage
	keyAES  []byte
}

func encryptAES(key []byte, plaintext string) string {
	c, err := aes.NewCipher(key)
	checkError(err)

	out := make([]byte, len(plaintext))

	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
}

func decryptAES(key []byte, ct string) string {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	checkError(err)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	return string(pt[:])
}

func (ss *SimpleSecretService) ReadSecret(key string) (string, error) {
	encryptedVal, err := ss.storage.Read(key)
	checkError(err)
	return decryptAES(ss.keyAES, encryptedVal), nil
}

func (ss *SimpleSecretService) WriteSecret(key string, value string) {
	ss.storage.Write(key, encryptAES(ss.keyAES, value))
}

func New(st storage.Storage) *SimpleSecretService {
	//look for a file located on the file system
	// if there is no file --> creates file --> writes down random key
	data := CheckFile("secretkey.txt")
	return &SimpleSecretService{
		storage: st,
		keyAES:  data,
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
