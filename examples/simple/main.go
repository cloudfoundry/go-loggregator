package main

import (
	"log"
	"os"

	"code.cloudfoundry.org/go-loggregator"
)

func main() {
	tlsConfig, err := loggregator.NewIngressTLSConfig(
		os.Getenv("CA_CERT_PATH"),
		os.Getenv("CERT_PATH"),
		os.Getenv("KEY_PATH"),
	)
	if err != nil {
		log.Fatal("Could not create TLS config", err)
	}

	client, err := loggregator.NewIngressClient(
		tlsConfig,
		loggregator.WithAddr("localhost:3458"),
	)

	if err != nil {
		log.Fatal("Could not create client", err)
	}

	client.EmitLog("some log goes here")
}
