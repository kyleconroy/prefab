.PHONY: build test fmt run stack

test: fmt
	cd stackgo && go test

fmt:
	gofmt -l -w . stackgo

build: fmt
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o stack

run: fmt
	go run software/postgresql/manifest.json

stack.linux64.tar.gz: build
	tar -czf stack.linux64.tar.gz stack


