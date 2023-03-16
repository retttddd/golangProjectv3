package service

type storage interface {
	Read(key string) (string, error)
	Write(key string, value string) error
}
