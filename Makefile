.PHONY: build build-web dev dev-web test test-coverage lint lint-web fmt clean package-deb security-scan docs help

## Defaults
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0
BIN_NAME = hostanything
BIN_DIR = bin

## Help target
help: ## Show this help
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

## Build Targets
build: ## Build the Go backend binary
	@echo "Building Go backend..."
	cd core && CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ../$(BIN_DIR)/$(BIN_NAME) ./cmd/$(BIN_NAME)

build-web: ## Build the frontend
	@echo "Building frontend..."
	cd web && npm install && npm run build

## Development Targets
dev: ## Run the Go backend dev server
	@echo "Starting Go dev server..."
	cd core && go run ./cmd/$(BIN_NAME)

dev-web: ## Run the frontend dev server
	@echo "Starting frontend dev server..."
	cd web && npm run dev

## Testing Targets
test: ## Run Go tests
	@echo "Running tests..."
	cd core && go test -v ./...

test-coverage: ## Run Go tests with coverage report
	@echo "Running tests with coverage..."
	cd core && go test -coverprofile=../coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated at coverage.html"

## Linting & Formatting
lint: ## Run golangci-lint
	@echo "Linting Go code..."
	cd core && golangci-lint run ./...

lint-web: ## Run frontend linters
	@echo "Linting frontend code..."
	cd web && npm run lint

fmt: ## Format Go code
	@echo "Formatting Go code..."
	cd core && go fmt ./...
	@if command -v goimports >/dev/null; then \
		cd core && goimports -w .; \
	else \
		echo "goimports not installed, skipping..."; \
	fi

## Clean Targets
clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BIN_DIR)/
	rm -f coverage.txt coverage.html
	rm -rf web/dist/ web/build/

## Packaging Targets
package-deb: build ## Build .deb package
	@echo "Building Debian package..."
	bash scripts/package-deb.sh

## Security Targets
security-scan: ## Run security scans
	@echo "Running govulncheck..."
	cd core && govulncheck ./...
	@echo "Running trivy scan..."
	trivy fs .

## Docs Targets
docs: ## Serve documentation locally
	@echo "Serving docs..."
	# Placeholder for docs server (e.g., mkdocs serve)
	@echo "Docs target is a placeholder."
