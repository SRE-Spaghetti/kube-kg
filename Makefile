# Makefile for the kube-kg project

.PHONY: all build run test clean lint security

# Default target: builds the application
all: build

# Build the Go application binary
# Output will be an executable named 'kube-kg' in the project root
build:
	@echo "Building kube-kg..."
	go build -o kube-kg ./cmd/kube-kg

# Run the application
# This depends on the 'build' target to ensure the binary is up-to-date
run: build
	@echo "Running kube-kg..."
	./kube-kg

# Run all tests in the project
test:
	@echo "Running tests..."
	go test ./...

# Lint the code using golangci-lint
# To install golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
lint:
	@echo "Linting code..."
	golangci-lint run ./...

# Run a Trivy filesystem scan for security vulnerabilities
# Requires Trivy to be installed. See https://aquasecurity.github.io/trivy/v0.51/getting-started/installation/
security:
	@echo "Running Trivy filesystem scan..."
	trivy fs .

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f kube-kg
