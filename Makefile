.PHONY: build test fmt run

test: fmt
	go test

fmt:
	gofmt -l -w .

stack:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o stack cli/main.go

build: fmt
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o stack cli/main.go

run: fmt
	go run cli/main.go postgresql/manifest.json

stack.linux64.tar.gz: stack
	tar -czf stack.linux64.tar.gz stack


