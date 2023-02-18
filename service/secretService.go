package service

type SecretService interface {
	ReadSecret(key string) (string, error)
	WriteSecret(key string, value string)
}
