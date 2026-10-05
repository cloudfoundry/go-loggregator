#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]:-$0}")" &> /dev/null && pwd)"
source "$SCRIPT_DIR/.protobuf-versions"
TMP_DIR="$(mktemp -d)"

cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

git init --quiet "$TMP_DIR/loggregator-api"
git -C "$TMP_DIR/loggregator-api" remote add origin https://github.com/cloudfoundry/loggregator-api.git
git -C "$TMP_DIR/loggregator-api" fetch --depth 1 origin "$LOGGREGATOR_API_REF"
git -C "$TMP_DIR/loggregator-api" checkout --quiet FETCH_HEAD

printf 'protoc: '
protoc --version
printf 'protoc-gen-go: '
protoc-gen-go --version
printf 'protoc-gen-go-grpc: '
protoc-gen-go-grpc --version
printf 'loggregator-api: %s\n' "$LOGGREGATOR_API_REF"

pushd "$SCRIPT_DIR/../.." > /dev/null
    protoc -I="$TMP_DIR" \
        --go_out=. \
        --go-grpc_out=. \
        "$TMP_DIR"/loggregator-api/v2/*.proto
    mv code.cloudfoundry.org/go-loggregator/v10/rpc/loggregator_v2/* rpc/loggregator_v2/
    rm -rf code.cloudfoundry.org
popd > /dev/null
