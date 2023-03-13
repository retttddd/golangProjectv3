package service

type SecretService interface {
	ReadSecret(key string, password string) (string, error)
	WriteSecret(key string, value string, password string) error
}
