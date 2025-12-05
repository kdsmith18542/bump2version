# bumpx Makefile

.PHONY: build test lint clean install

# Default target
all: build

# Build the binary
build:
	go build -o bumpx ./cmd/bumpx/

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
	rm -f bumpx bumpx.exe
	rm -f coverage.out coverage.html

# Install to GOPATH/bin
install:
	go install ./cmd/bumpx/

# Build for multiple platforms
build-all: clean
	GOOS=linux GOARCH=amd64 go build -o bumpx-linux-amd64 ./cmd/bumpx/
	GOOS=linux GOARCH=arm64 go build -o bumpx-linux-arm64 ./cmd/bumpx/
	GOOS=darwin GOARCH=amd64 go build -o bumpx-darwin-amd64 ./cmd/bumpx/
	GOOS=darwin GOARCH=arm64 go build -o bumpx-darwin-arm64 ./cmd/bumpx/
	GOOS=windows GOARCH=amd64 go build -o bumpx-windows-amd64.exe ./cmd/bumpx/
