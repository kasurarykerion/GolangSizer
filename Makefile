# Open source image resizer coded by kasuraSH
# Makefile for building and testing the image resizer

# Binary name
BINARY_NAME=golangresizer
BINARY_PATH=bin/$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags following Power of 10 Rule 10: Compile with all warnings
BUILDFLAGS=-v -ldflags="-s -w"
TESTFLAGS=-v -race -coverprofile=coverage.out

.PHONY: all build clean test fmt vet deps help

all: deps fmt vet build

## build: Build the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	$(GOBUILD) $(BUILDFLAGS) -o $(BINARY_PATH) ./cmd/golangresizer

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/
	@rm -f coverage.out

## test: Run tests with race detection
test:
	@echo "Running tests..."
	$(GOTEST) $(TESTFLAGS) ./...

## fmt: Format Go source code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## vet: Run go vet for static analysis
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## help: Display this help message
help:
	@echo "Available targets:"
	@grep -E '^##' Makefile | sed 's/##//'
