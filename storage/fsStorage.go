package storage

import (
	"encoding/json"
	"io/ioutil"
)

type FsStorage struct {
	path string
}

type myData struct {
	IntValue int
}

func (st FsStorage) Read(key string) (string, error) {
	return "got it", nil
}
func (st FsStorage) Write(key string, value string) {
	a, _ := json.Marshal(map[string]string{key: value})
	_ = ioutil.WriteFile("test.json", a, 0644)
}

func New() Storage {
	return FsStorage{
		path: "p",
	}
}
