package main

import (
	"log"

	"github.com/TudorHulban/http-server/server"
)

func main() {
	certFile := "cert.pem" // Path to your certificate file
	keyFile := "key.pem"   // Path to your key file
	address := ":443"

	serverHTTPS, errCrServer := server.NewServer(certFile, keyFile)
	if errCrServer != nil {
		log.Fatal(errCrServer)
	}

	if errRun := serverHTTPS.Run(address); errRun != nil {
		log.Fatal(errRun)
	}
}
