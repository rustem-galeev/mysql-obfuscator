package main

import (
	"obfuscator/httpServer"
)

func main() {
	err := httpServer.InitServer()
	if err != nil {
		panic(err)
	}
}
