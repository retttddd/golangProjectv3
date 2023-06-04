package storage

import (
	"awesomeProject3/service"
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

func (st fsStorage) Read(key string) (*service.StorageModel, error) {
	file, err := ioutil.ReadFile(st.path)
	if err != nil {
		return nil, err
	}
	data := make(map[string]dataBox)
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	val, ok := data[key]
	if ok {
		return &service.StorageModel{Value: &val.Value}, nil
	}
	return nil, errors.New("item was not found")
}

func (st fsStorage) Write(key string, model *service.StorageModel) error {
	err := checkFile(st.path)
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile(st.path)
	if err != nil {
		return err
	}
	data := make(map[string]dataBox)
	if len(file) != 0 {
		err = json.Unmarshal(file, &data)
		if err != nil {
			return err
		}
	}

	data[key] = dataBox{Value: *model.Value}
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

func NewFsStorage(p string) *fsStorage {
	return &fsStorage{
		path: p,
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
