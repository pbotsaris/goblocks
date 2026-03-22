.PHONY: test test-verbose test-cover lint build

# Run all tests
test:
	go test ./blocks/...

# Verbose test output
test-verbose:
	go test -v ./blocks/...

# Run tests with coverage
test-cover:
	go test -cover ./blocks/...

# Run tests with coverage report
test-cover-html:
	go test -coverprofile=coverage.out ./blocks/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter (install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
lint:
	golangci-lint run ./blocks/...

# Build check (ensures everything compiles)
build:
	go build ./blocks/...

# Run all checks (lint + test)
check: lint test
