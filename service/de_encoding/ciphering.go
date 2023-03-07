package de_encoding

type Encoder interface {
	Encrypt(plaintext string) string
	Decrypt(ct string) string
}
