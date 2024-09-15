.PHONY: all fmt vet lint test clean

# Default target
all: fmt vet lint test

# Format all Go files
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run golangci-lint
lint:
	golangci-lint run

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -rf ./bin

# Build the application
build:
	go build -o bin/$(BINARY_NAME) .

# Run the application
run: build
	./bin/$(BINARY_NAME)
