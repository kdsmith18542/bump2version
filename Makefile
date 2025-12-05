# bumpx Makefile

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.Version=$(VERSION)

.PHONY: build test lint clean install release

# Default target
all: build

# Build the binary (development)
build:
	go build -o bumpx ./cmd/bumpx/

# Build optimized release binary (static, stripped)
release:
	CGO_ENABLED=0 go build -ldflags='$(LDFLAGS)' -o bumpx ./cmd/bumpx/

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -f bumpx bumpx.exe bumpx-*
	rm -f coverage.out coverage.html

# Install to GOPATH/bin
install:
	go install -ldflags='$(LDFLAGS)' ./cmd/bumpx/

# Build for multiple platforms (release builds)
build-all: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='$(LDFLAGS)' -o bumpx-linux-amd64 ./cmd/bumpx/
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags='$(LDFLAGS)' -o bumpx-linux-arm64 ./cmd/bumpx/
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags='$(LDFLAGS)' -o bumpx-darwin-amd64 ./cmd/bumpx/
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags='$(LDFLAGS)' -o bumpx-darwin-arm64 ./cmd/bumpx/
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags='$(LDFLAGS)' -o bumpx-windows-amd64.exe ./cmd/bumpx/
