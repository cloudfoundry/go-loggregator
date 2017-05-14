# go-loggregator [![slack.cloudfoundry.org][slack-badge]][loggregator-slack]

This is a golang client library for the [Loggregator v2 API][loggregator-api].

## Usage

This repository should be imported as:

`import loggregator "code.cloudfoundry.org/go-loggregator"`

## Example

Example implementation of the client is provided in `examples/main.go`.

Build the example client by running `go build -o client main.go`

Collocate the `client` with a metron agent and set the following environment
variables: `CA_CERT_PATH`, `CERT_PATH`, `KEY_PATH`

[slack-badge]:              https://slack.cloudfoundry.org/badge.svg
[loggregator-slack]:        https://cloudfoundry.slack.com/archives/loggregator
[loggregator-api]:          https://github.com/cloudfoundry/loggregator-api
