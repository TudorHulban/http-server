package main

import (
	"log"
)

func main() {
	server := NewServer()

	log.Fatal(server.Run(":8080"))
}
