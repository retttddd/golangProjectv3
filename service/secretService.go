package service

type SecretService interface {
	ReadSecret(key string) (string, error)
	WriteSecret(key, value string) error
}
