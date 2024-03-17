
# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

.PHONY: all build run test clean

.PHONY: modsync
modsync: 
	go mod tidy
	go mod vendor

.PHONY: docker
docker: modsync
	docker build -t gitfit-go .
	docker run --rm -it -p 8080:8080 gitfit-go:latest
		