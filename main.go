package main

import "log"

func main() {
	certFile := "cert.pem" // Path to your certificate file
	keyFile := "key.pem"   // Path to your key file
	address := ":443"

	server, errCrServer := NewServer(certFile, keyFile)
	if errCrServer != nil {
		log.Fatal(errCrServer)
	}

	if errRun := server.Run(address); errRun != nil {
		log.Fatal(errRun)
	}
}
