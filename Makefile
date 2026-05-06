.PHONY: build run dev clean install

# Build the application
build:
	@echo "Building WaWeb V2..."
	@go build -o waweb main.go

# Run the built binary
run: build
	@echo "Starting WaWeb V2..."
	@./waweb

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without hot reload..."; \
		go run main.go; \
	fi

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f waweb
	@rm -rf published_sites/*
	@echo "Done!"

# Run tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Check for errors
vet:
	@echo "Vetting code..."
	@go vet ./...
