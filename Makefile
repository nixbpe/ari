.PHONY: build test lint coverage bench setup

build:
	go build ./cmd/ari

test:
	go test -race --count=1 -coverprofile=coverage.out -covermode=atomic ./...

lint:
	golangci-lint run ./...

coverage:
	go test -race --count=1 -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

bench:
	go test -bench=./... -benchmem -run=^$ ./...

setup:
	go mod download
	go install github.com/vladopajic/go-test-coverage/v2@latest
