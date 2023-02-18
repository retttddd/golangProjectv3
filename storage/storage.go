package storage

type Storage interface {
	Read(key string) (string, error)
	Write(key string, value string)
}
