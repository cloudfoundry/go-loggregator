#!/usr/bin/env bash
#
# This script generates all the data used for testing.
# Usage: `scripts/update-testdata.sh`.

set -eu

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
DATA_PATH="${REPO_ROOT}/testdata"

echo "Installing certstrap..."
go install github.com/square/certstrap

echo "Cleaning up old generated data..."
rm -rf "${DATA_PATH}"

echo "Generating new data..."

# Ensure there's a directory there
mkdir -p "${DATA_PATH}"

# Create Certificate Authority, including certificate, key and extra information file
certstrap --depot-path "${DATA_PATH}" init --passphrase '' --cn loggregator

# Create and sign certificate for metron
certstrap --depot-path "${DATA_PATH}" request-cert --passphrase '' --cn 'metron' --domain 'metron'
certstrap --depot-path "${DATA_PATH}" sign --CA loggregator metron

# Create and sign certificate for the server
certstrap --depot-path "${DATA_PATH}" request-cert --passphrase '' --cn 'reverselogproxy' --domain 'reverselogproxy'
certstrap --depot-path "${DATA_PATH}" sign --CA loggregator reverselogproxy
