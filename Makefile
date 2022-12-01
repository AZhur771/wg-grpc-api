BIN := "./bin/wgapi"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	protoc --go_out=gen/ \
		  --go-grpc_out=gen/ \
		  --grpc-gateway_out=gen/ \
		  --grpc-gateway_opt=allow_delete_body=true \
		  --grpc-gateway_opt generate_unbound_methods=true \
  		  --proto_path=third_party/ \
  		  --proto_path=api/proto/ \
  		  --openapiv2_out=third_party/openapiv2 \
  		  --openapiv2_opt allow_delete_body=true \
		  api/proto/peer_service.proto

protoc-version:
	which protoc && protoc --version

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

go-version:
	which go && go version

.PHONY: generate protoc-version install-lint-deps lint go-version
