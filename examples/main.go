package main

import (
	"log"
	"os"

	loggregator "code.cloudfoundry.org/go-loggregator"
)

func main() {
	metronCfg := loggregator.Config{
		UseV2API:      true,
		APIPort:       3458,
		CACertPath:    os.Getenv("CA_CERT_PATH"),
		CertPath:      os.Getenv("CERT_PATH"),
		KeyPath:       os.Getenv("KEY_PATH"),
		JobDeployment: "example-deployment",
		JobName:       "example-job",
		JobIndex:      "example-index",
		JobIP:         "0.0.0.0",
		JobOrigin:     "example-deployment",
	}

	client, err := loggregator.NewClient(metronCfg)

	if err != nil {
		log.Fatal("Could not create client", err)
	}

	client.SendMetric("some-metric-name", 1234)
}
