.PHONY: build test fmt run

GOPATH = $(HOME)/gopath
export GOPATH


fmt:
	gofmt -l -w .

build: fmt
	go build -o stack cli/main.go

run: fmt
	go run cli/main.go analyze

test: fmt
	go test

converge:
	go run cli/main.go postgresql.json

