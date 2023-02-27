package main

import (
	"awesomeProject3/service"
	"awesomeProject3/service/de_encoding"
	"awesomeProject3/storage"
	"fmt"
)

func main() {
	passWord := "hello"
	srv := service.New(storage.New(), de_encoding.NewAESEncoder(de_encoding.PassToSecretKey(passWord)))
	//srv.WriteSecret("key18", "ThatsValue18")
	//srv.WriteSecret("key19", "THatsValue19")
	//srv.WriteSecret("key20", "ThatsValue20")
	//srv.WriteSecret("key21", "THatsValue21")
	value, _ := srv.ReadSecret("key19")
	fmt.Print(value)

}
