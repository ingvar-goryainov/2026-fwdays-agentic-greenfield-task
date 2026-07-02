.PHONY: build test lint

build:
	go build ./cmd/agents-lint

test:
	go test ./...

lint:
	go vet ./...
