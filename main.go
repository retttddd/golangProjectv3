package main

import (
	"awesomeProject3/service"
	storage "awesomeProject3/storage"
	"fmt"
)

func main() {
	srv := service.New(storage.New())
	value, _ := srv.ReadSecret("thisis32bitlongpassphraseimusing")
	fmt.Print(value)
	srv.WriteSecret("thisis32bitlongpassphraseimusing", "secret value")

}
