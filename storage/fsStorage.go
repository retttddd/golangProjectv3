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
type myStruct struct {
	Key   string
	Value string
}

func (st fsStorage) Read(key string) (string, error) {
	allDataFromJson := []myStruct{}
	file, err := ioutil.ReadFile(st.path)
	if err != nil {
		return "", err
	}
	json.Unmarshal(file, &allDataFromJson)
	if err != nil {
		return "", err
	}
	for _, v := range allDataFromJson {
		if v.Key == key {
			return v.Value, nil
		}
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

	data := []myStruct{}

	json.Unmarshal(file, &data)

	newStruct := &myStruct{
		Key:   key,
		Value: value,
	}

	data = append(data, *newStruct)
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

func NewFsStorage() fsStorage {
	return fsStorage{
		path: "test.json",
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
	}
	return nil
}
