.PHONY: build test fmt

GOPATH = $(HOME)/gopath
export GOPATH


fmt:
	gofmt -l -w .

build:
	go build -o stack cli/main.go

run:
	go run cli/main.go analyze

