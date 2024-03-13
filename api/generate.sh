#!/bin/sh -eu

if ! command -v protoc-gen-go &> /dev/null; then
  go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
fi

for proto_path in v1alpha1; do

  protoc \
    --go_opt=module=github.com/fillmore-labs/microbatch-lambda/api \
    --go_opt=paths=import \
    --go_out=. \
    proto/$proto_path/*.proto

done
