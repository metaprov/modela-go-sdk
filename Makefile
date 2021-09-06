.PHONY: build build-alpine clean test help default gen

PROTO_ROOT       =https://raw.githubusercontent.com/metaprov/modela-api/master/modela-api
VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
RELEASE_VERSION= 0.1.13

.PHONY: all
all: help

.PHONY: tidy
tidy: ## Updates the go modules
	go mod tidy

.PHONY: test
test: tidy ## Tests the entire project
	go test -count=1 \
			-race \
			-coverprofile=coverage.txt \
			-covermode=atomic \
			./...

.PHONY: spell
spell: ## Checks spelling across the entire project
	@command -v misspell > /dev/null 2>&1 || go get github.com/client9/misspell/cmd/misspell
	@misspell -locale US -error go=golang client/**/* example/**/* .

.PHONY: cover
cover: tidy ## Displays test coverage in the client and service packages
	go test -coverprofile=cover-client.out ./client && go tool cover -html=cover-client.out
	go test -coverprofile=cover-grpc.out ./protos/prediction-server/v1  && go tool cover -html=cover-grpc.out

.PHONY: lint
lint: ## Lints the entire project
	golangci-lint run --timeout=3m

.PHONY: tag
tag: ## Creates release tag
	git tag $(RELEASE_VERSION)
	git push origin $(RELEASE_VERSION)

.PHONY: clean
clean: ## Cleans the generated files
	go clean
	rm -rf ./proto/prediction-server/v1/*


.PHONY: gen
gen: ## Downloads proto files from modela/modela-api master and generates gRPC proto clients
	go install github.com/gogo/protobuf/gogoreplace

	rm -rf ./protos/*

	mkdir -p ./protos/prediction-server/v1/

	wget -q https://raw.githubusercontent.com/metaprov/modela-api/main/prediction-server/v1/prediction_server.proto -O ./protos/prediction-server/v1/prediction_server.proto
	gogoreplace 'option go_package = "github.com/metaprov/modela/pkg/proto/predictionserver/v1;prediction-server";' \
		'option go_package = "github.com/metaprov/modela-go-sdk/modela/proto/predictionserver/v1;prediction-server";' \
		./protos/prediction-server/v1/prediction_server.proto

	protoc -Icommon-protos \
		   -I. \
		   --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		   ./protos/prediction-server/v1/prediction_server.proto

	rm -f ./protos/prediction_server/v1/prediction_server.proto


.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk \
		'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: install-tools
install-tools:
	go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get google.golang.org/grpc@v1.32.0