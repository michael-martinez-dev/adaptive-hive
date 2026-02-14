.PHONY: fmt build test lint clean

fmt:
	go fmt ./...

build:
	go build -o bin/agent ./cmd/agent/
	go build -o bin/hivectl ./cmd/hivectl/

test:
	go test -race -v ./...

cover:
	go test -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

test-int:
	go test -race -v -tags=int ./test/...

