package service

type SecretService interface {
	ReadSecret(key string) (string, error)
}
