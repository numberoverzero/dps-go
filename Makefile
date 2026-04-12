.PHONY: all _init build test tidy

all: _init build tidy test

_init:
	@git config core.hooksPath .githooks

build:
	go build ./...

test:
	go test ./...

tidy:
	go mod tidy	
	go vet ./...
	golangci-lint run --fix ./...
