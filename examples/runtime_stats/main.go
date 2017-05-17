package main

import (
	"log"
	"os"
	"time"

	"code.cloudfoundry.org/go-loggregator/runtimeemitter"
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
		v2.WithLogger(log.New(os.Stdout, "", log.LstdFlags)),
	)

	if err != nil {
		log.Fatal("Could not create client", err)
	}

	runtimeStats := runtimeemitter.New(
		client,
		runtimeemitter.WithInterval(10*time.Second),
	)

	runtimeStats.Run()
}
