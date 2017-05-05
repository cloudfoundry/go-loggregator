#!/usr/bin/env bash

LOGGREGATOR_API=$GOPATH/src/code.cloudfoundry.org/loggregator-api

go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

protoc --proto_path=$LOGGREGATOR_API/v2 \
       --go_out=plugins=grpc:. \
       $LOGGREGATOR_API/v2/ingress.proto \
       $LOGGREGATOR_API/v2/envelope.proto
