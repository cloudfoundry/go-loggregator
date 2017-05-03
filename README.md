# go-loggregator

This is a golang client library for the [Loggregator v2 API](https://github.com/cloudfoundry/loggregator-api).

**WARNING:** It is in alpha release and unstable.

## Usage

This repository should be imported as:

`import "code.cloudfoundry.org/go-loggregator/loggregator_v2"`

## Example

Example implementation of the client is provided in `examples/main.go`.

Build the example client by running `go build -o client main.go`

Collocate the `client` with a metron agent and set the following environment
variables: `CA_CERT_PATH`, `CERT_PATH`, `KEY_PATH`
