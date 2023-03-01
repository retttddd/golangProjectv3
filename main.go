package main

import (
	"awesomeProject3/cli"
)

func main() {

	//passWord := "hello2"
	//srv := service.New(storage.New(), de_encoding.NewAESEncoder(de_encoding.PassToSecretKey(passWord)))
	////srv.WriteSecret("key18", "ThatsValue18") //hello
	////srv.WriteSecret("key19", "THatsValue19") //hello
	////srv.WriteSecret("key20", "ThatsValue20") //hello
	////srv.WriteSecret("key21", "THatsValue21") //hello
	////srv.WriteSecret("key22", "THatsValue22") //hello2
	//value, _ := srv.ReadSecret("key22")
	//fmt.Print(value)
	cli.Execute()

}
