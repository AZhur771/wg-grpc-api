BIN := "./bin/wg-grpc-api"

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
		  --openapiv2_opt allow_merge=true \
		  --openapiv2_opt merge_file_name=api \
		  api/proto/*.proto

protoc-version:
	which protoc && protoc --version

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

go-version:
	which go && go version

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" .

run: build
	$(BIN)

app-version: build
	$(BIN) version

test:
	go test -race ./...

test-proto:
	docker run --net=host --rm -v "$(pwd):/mount:ro" \
		ghcr.io/ktr0731/evans:latest \
		--port 51821 \
		-r \
		repl

cert:
	cd test_certs; ./gen.sh; cd ..

.PHONY: generate protoc-version install-lint-deps lint go-version
