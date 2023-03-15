package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type fsStorage struct {
	path string
}
type dataBox struct {
	Value string `json:"value"`
}

func (st fsStorage) Read(key string) (string, error) {
	file, err := ioutil.ReadFile(st.path)
	if err != nil {
		return "", err
	}
	data := make(map[string]dataBox)
	err = json.Unmarshal(file, &data)
	if err != nil {
		return "", err
	}

	val, ok := data[key]
	if ok {
		return val.Value, nil
	}
	return "", errors.New("item was not found")
}

func (st fsStorage) Write(key string, value string) error {
	err := checkFile(st.path)
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(st.path)
	if err != nil {
		return err
	}
	data := make(map[string]dataBox)
	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	data[key] = dataBox{Value: value}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(st.path, dataBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func NewFsStorage() *fsStorage {
	return &fsStorage{
		path: "./data/test.json",
	}
}

func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()
		return nil
	}

	return err
}
