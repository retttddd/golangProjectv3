package main

import (
	"awesomeProject3/service"
	de_encoding "awesomeProject3/service/de_encoding"
	storage "awesomeProject3/storage"
	"fmt"
)

func main() {
	passWord := "secretkey"
	srv := service.New(storage.New(), de_encoding.NewAESEncoder(de_encoding.CheckFile(passWord+".txt")))
	srv.WriteSecret("key1", "value")
	value, _ := srv.ReadSecret("key1")
	fmt.Print(value)

}
