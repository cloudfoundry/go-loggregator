package main

import (
	"log"
	"os"
	"time"

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

	for i := 0; i < 50; i++ {
		client.EmitLog("some log goes here")
		client.EmitEvent("event-title", "event-body")

		time.Sleep(10 * time.Millisecond)
	}

	client.CloseSend()
}
