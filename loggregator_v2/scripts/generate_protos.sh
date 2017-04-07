#!/usr/bin/env bash

LOGGREGATOR_API=$GOPATH/src/code.cloudfoundry.org/loggregator-api

rm $GOPATH/bin/protoc-gen-gogoslick
go install -v github.com/gogo/protobuf/protoc-gen-gogoslick

protoc --proto_path=$LOGGREGATOR_API/v2 \
       --go_out=plugins=grpc:. \
       $LOGGREGATOR_API/v2/*.proto
