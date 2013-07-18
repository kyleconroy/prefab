.PHONY: build test fmt run

test: fmt
	cd prefab && go test

fmt:
	gofmt -l -w . prefab

build: fmt
	mkdir -p bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/prefab

run: fmt
	go run software/postgresql/manifest.json

prefab.linux64.tar.gz: build
	tar -czf prefab.linux64.tar.gz pf


