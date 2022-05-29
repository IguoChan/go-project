#!/bin/bash

set -ex

# gen demo_app
protoc -I"$1" \
  -I"$1"/.. \
  --go_out="$2" --go_opt paths=source_relative \
  --go-grpc_out="$2" --go-grpc_opt paths=source_relative \
  --grpc-gateway_out="$2" --grpc-gateway_opt paths=source_relative --grpc-gateway_opt logtostderr=true \
  "$1"/**/*.proto