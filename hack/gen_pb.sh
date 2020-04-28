#! /usr/bin/env bash

# This script auto-generates protobuf related files. It is intended to be run manually when either
# API types are added/modified, or server gRPC calls are added. The generated files should then
# be checked into source control.

set -x
set -o errexit
set -o nounset
set -o pipefail

PROJECT_ROOT=$(cd $(dirname ${BASH_SOURCE})/..; pwd)
PATH="${PROJECT_ROOT}/dist:${PATH}"

# Generate server/<service>/(<service>.pb.go|<service>.pb.gw.go)
PROTO_FILES=proto/*

# generate the api objects
GOOGLE_PROTO_API_PATH=${PROJECT_ROOT}/common-protos
GOGO_PROTOBUF_PATH=${PROJECT_ROOT}/common-protos/github.com/gogo/protobuf


protoc \
        -I${PROJECT_ROOT} \
        -I/usr/local/include \
        -I./vendor \
        -I$GOPATH/src \
        -I${GOOGLE_PROTO_API_PATH} \
        proto/prediction_server.proto --go_out=plugins=grpc:$GOPATH/src





