package main

import (
	"awesomeProject3/service"
	storage "awesomeProject3/storage"
	"fmt"
)

func main() {
	srv := service.New(storage.New())
	value, _ := srv.ReadSecret("qweqweqwq")
	fmt.Print(value)
	srv.WriteSecret("weqweqeqwe", "qweqweqwe")

}
