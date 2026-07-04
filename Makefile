.PHONY: build test lint release clean

build:
	go build ./cmd/agents-lint

test:
	go test ./...

lint:
	go vet ./...

# Cross-compiled, stripped release binaries (NFR-DIST-01, NFR-PERF-02).
# CGO_ENABLED=0 keeps them static; -trimpath keeps builds reproducible;
# -ldflags="-s -w" strips debug info to help stay under the 15MB ceiling.
release:
	mkdir -p dist
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o dist/agents-lint-linux-amd64 ./cmd/agents-lint
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o dist/agents-lint-darwin-arm64 ./cmd/agents-lint

clean:
	rm -rf dist
