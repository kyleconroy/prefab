.PHONY: build test fmt run

fmt:
	gofmt -l -w .

build: fmt
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o stack cli/main.go

run: fmt
	go run cli/main.go postgresql/manifest.json

test: fmt
	go test

