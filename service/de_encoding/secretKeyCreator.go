package de_encoding

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"math/rand"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func shaEncrypt(fileNameFunc string) []byte {
	h := sha256.New()
	h.Write([]byte(fileNameFunc))
	encryptedKeyWord := h.Sum(nil)
	return encryptedKeyWord
}

func CheckFile(name string) []byte {
	randomKey, _ := GenerateRandomBytes(32)
	fileName := shaEncrypt(name)
	data, err := ioutil.ReadFile(string(fileName))
	if err != nil {
		fmt.Println(err)
	}
	if data != nil {
		return data
	} else {
		err := ioutil.WriteFile(string(fileName), randomKey, 0777)

		if err != nil {
			fmt.Println(err)
		}
		return randomKey
	}
}
