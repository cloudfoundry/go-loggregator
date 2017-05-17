package main

import (
	"log"
	"os"

	"code.cloudfoundry.org/go-loggregator/v2"
)

func main() {
	tlsConfig, err := v2.NewTLSConfig(
		os.Getenv("CA_CERT_PATH"),
		os.Getenv("CERT_PATH"),
		os.Getenv("KEY_PATH"),
	)
	if err != nil {
		log.Fatal("Could not create TLS config", err)
	}

	client, err := v2.NewClient(
		tlsConfig,
		v2.WithPort(3458),
	)

	if err != nil {
		log.Fatal("Could not create client", err)
	}

	client.EmitLog("some log goes here")
}
