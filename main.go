package main

import (
	"awesomeProject3/cli"
	"os"
)

func main() {
	if err := cli.Execute(); err != nil {
		//fmt.Fprintf(os.Stderr, "'%s'", err)
		os.Exit(1)
	}
}
