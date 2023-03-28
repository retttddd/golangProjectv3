package ciphering

import (
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
func (rnp *randomNonceProducer) generate(size int) (string, error) {
	nonce := make([]byte, size)

	if _, err := io.ReadFull(rnp.rndReader, nonce); err != nil {
		return "", err
	}
	return string(nonce[:]), nil
}
