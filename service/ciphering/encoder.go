package ciphering

type Encoder interface {
	Encrypt(plaintext string, aesKey []byte, useRandom bool) (string, error)
	Decrypt(ct string, aesKey []byte, useRandom bool) ([]byte, error)
}
