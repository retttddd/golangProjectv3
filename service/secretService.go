package service

type SecretService interface {
	ReadSecret(key []byte) (string, error)
	WriteSecret(key []byte, value string)
}
