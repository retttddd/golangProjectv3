package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const filename = "test.json"

type FsStorage struct {
	path string
}
type MyStruct struct {
	Key   string
	Value string
}

func (st FsStorage) Read(key string) (string, error) {
	allDataFromJson := []MyStruct{}
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("action failed: ", err)
	}
	json.Unmarshal(file, &allDataFromJson)
	if err != nil {
		log.Println("action failed: ", err)
	}
	for _, v := range allDataFromJson {
		if v.Key == key {
			return v.Value, nil
		}
	}
	return "", nil
}
func (st FsStorage) Write(key string, value string) {
	err := checkFile(filename)
	if err != nil {
		log.Println("action failed: ", err)
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("action failed: ", err)
	}

	data := []MyStruct{}

	json.Unmarshal(file, &data)

	newStruct := &MyStruct{
		Key:   key,
		Value: value,
	}

	data = append(data, *newStruct)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("action failed: ", err)
	}

	err = ioutil.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		log.Println("action failed: ", err)
	}
}

func New() Storage {
	return FsStorage{
		path: "p",
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
		if err != nil {
			return err
		}
	}
	return nil
}
