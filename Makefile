.PHONY: help run test test-coverage build docker-build docker-run clean

# Default target
help:
    @echo "Available commands:"
	@echo "	make run		- Run the application locally"
	@echo "	make test		- Run unit tests"
	@echo "	make test-coverage	- Run tests with coverage report"
	@echo "	make build		- Build the binary"
	@echo "	make docker-build	- Build Docker image"
	@echo "	make docker-run		- Run Docker container"
	@echo "	make clean		- Clean build artifacts"

# Run the application
run:
	@go run cmd/server/main.go

# Run tests
test:
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@go test -v -cover -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build binary
build:
	@mkdir -p bin
	@go build -o bin/pack-calculator cmd/server/main.go

# Build Docker image
docker-build:
	@docker build -f deployments/Dockerfile -t pack-calculator:latest .

# Run Docker container
docker-run:
	@docker run -p 8080:8080 --rm pack-calculator:latest

# Clean build artifacts
clean:
	@rm -rf bin/ coverage.out coverage.html