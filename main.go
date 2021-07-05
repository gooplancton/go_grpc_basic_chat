package main

import (
	"os"
)

func main() {
	args := os.Args[1:]

	if args[0] == "server" {
		mainServer()
	} else {
		mainClient(args[0])
	}
}
