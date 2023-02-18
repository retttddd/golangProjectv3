package storage

type Storage interface {
	Read(key []byte) (string, error)
	Write(key []byte, value string)
}
