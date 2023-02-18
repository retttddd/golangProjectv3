package service

import (
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

func CheckFile(name string) []byte {
	randomKey, _ := GenerateRandomBytes(32)
	data, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err)
	}
	if data != nil {
		return data
	} else {
		err := ioutil.WriteFile(name, randomKey, 0777)

		if err != nil {
			fmt.Println(err)
		}
		return randomKey
	}
}
