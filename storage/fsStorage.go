package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type FsStorage struct {
	path string
}
type MyStruct struct {
	StructData map[string]string
}

func (st FsStorage) Read(key string) (string, error) {
	return "got it", nil
}
func (st FsStorage) Write(key string, value string) {
	filename := "test.json"
	err := checkFile(filename)
	if err != nil {
		log.Println("action failed: ", err)
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("action failed: ", err)
	}

	data := []MyStruct{}

	// Here the magic happens!
	json.Unmarshal(file, &data)

	newStruct := &MyStruct{
		StructData: map[string]string{key: value},
	}

	data = append(data, *newStruct)

	// Preparing the data to be marshalled and written.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("action failed: ", err)
	}

	err = ioutil.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		log.Println("action failed: ", err)
	}
	//checkFile("test.json")
	//var gotData = map[string]string{key: value}
	//data, err := ioutil.ReadFile("test.json")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//var slice []map[string]string
	//err = json.Unmarshal(data, &slice)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//slice = append(slice, gotData)
	//a, _ := json.Marshal(slice)
	//_ = ioutil.WriteFile("test.json", a, 0644)

}

func New() Storage {
	return FsStorage{
		path: "p",
	}
}

func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}
