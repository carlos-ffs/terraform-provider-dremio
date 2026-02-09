default: fmt lint install generate

BINARY_NAME=terraform-provider-dremio
GOBIN?=$(shell go env GOPATH)/bin

build:
	go build -v -o $(BINARY_NAME) .

install: build
	mv $(BINARY_NAME) $(GOBIN)/$(BINARY_NAME)

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate
