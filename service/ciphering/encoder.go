package ciphering

type Encoder interface {
	Encrypt(plaintext string, aesKey []byte) (string, error)
	Decrypt(ct string, aesKey []byte) ([]byte, error)
}
