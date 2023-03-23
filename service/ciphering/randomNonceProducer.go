package ciphering

import (
	"crypto/rand"
	"io"
)

type randomNonceProducer struct {
	rndReader io.Reader
}

func NewRandomNonceProducer(rnd io.Reader) *randomNonceProducer {
	return &randomNonceProducer{
		rndReader: rnd,
	}
}
func (*randomNonceProducer) generate(size int) (string, error) {
	nonce := make([]byte, size)

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return string(nonce[:]), nil
}
