#! /usr/bin/env bash

# This script auto-generates protobuf related files. It is intended to be run manually when either
# API types are added/modified, or server gRPC calls are added. The generated files should then
# be checked into source control.

set -x
set -o errexit
set -o nounset
set -o pipefail

PROJECT_ROOT=$(cd $(dirname ${BASH_SOURCE})/..; pwd)
CODEGEN_PKG=${CODEGEN_PKG:-$(cd ${PROJECT_ROOT}; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}
PATH="${PROJECT_ROOT}/dist:${PATH}"

# Generate server/<service>/(<service>.pb.go|<service>.pb.gw.go)

# generate the api objects
GOOGLE_PROTO_API_PATH=${PROJECT_ROOT}/common-protos
GOGO_PROTOBUF_PATH=${PROJECT_ROOT}/common-protos/github.com/gogo/protobuf


protoc \
   -I${PROJECT_ROOT} \
   -I/usr/local/include \
   -Iproto \
   -I${GOOGLE_PROTO_API_PATH} \
   -I${GOGO_PROTOBUF_PATH} \
   --gogo_out=plugins=grpc:sdk/golang/predictorsdk \
    "prediction_server.proto"
        #--go_out=plugins=grpc:$GOPATH/src \