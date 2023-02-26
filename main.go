package main

import (
	"awesomeProject3/service"
	"awesomeProject3/service/de_encoding"
	"awesomeProject3/storage"
)

func main() {
	passWord := "hello"
	srv := service.New(storage.New(), de_encoding.NewAESEncoder(de_encoding.PassToSecretKey(passWord)))
	srv.WriteSecret("key10", "123356389125345678")
	//value, _ := srv.ReadSecret("key1")

}
