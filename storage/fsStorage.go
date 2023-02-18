package storage

import "fmt"

type FsStorage struct {
	path string
}

func (st FsStorage) Read(key string) (string, error) {
	return "got it", nil
}
func (st FsStorage) Write(key string, value string) {
	fmt.Print(value)
}

func New() Storage {
	return FsStorage{
		path: "p",
	}
}
