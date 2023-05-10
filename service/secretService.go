package service

type SecretService interface {
	ReadSecret(key string, password string) (*SecretServiceModel, error)
	WriteSecret(key string, model *SecretServiceModel, password string) error
}
