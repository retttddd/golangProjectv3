package service

type encoder interface {
	Encrypt(plaintext string, cipherKey []byte) (string, error)
	Decrypt(ct string, cipherKey []byte) (string, error)
}
