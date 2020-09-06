GOPATH?=$(HOME)/go
GO_VERSION:=$(shell go version)
GO_VERSION_SHORT:=$(shell echo $(GO_VERSION)|sed -E 's/.* go(.*) .*/\1/g')
BIN?=./bin/sadwave-events-tg

build:
	$(info #Building...)
	go build -race -o $(BIN) ./cmd/main.go
