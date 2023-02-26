package main

import (
	"awesomeProject3/service"
	de_encoding "awesomeProject3/service/de_encoding"
	storage "awesomeProject3/storage"
)

func main() {
	passWord := "secretkey"
	srv := service.New(storage.New(), de_encoding.NewAESEncoder(de_encoding.PassToSecretKey(passWord)))
	srv.WriteSecret("key1", "value")
	//value, _ := srv.ReadSecret("key1")
	
}
