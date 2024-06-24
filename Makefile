export PATH := $(shell go env GOPATH)/bin:$(PATH)
export GOTOOLCHAIN=go1.22.3

build:
	go build -o app

test:
	go test -v ./...

generate:
	go generate ./...

