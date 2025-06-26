# GitHub Profiler Makefile

BINARY_NAME=github-profiler
VERSION=1.0.0
BUILD_DIR=build
GO_VERSION=1.24.4

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: all build clean install test lint deps help

all: clean deps build

## Build the binary
build:
	@echo "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

## Build for multiple platforms
build-all:
	@echo "$(BLUE)Building for multiple platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@echo "$(GREEN)✓ Linux AMD64 build complete$(NC)"
	
	# Linux ARM64
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "$(GREEN)✓ Linux ARM64 build complete$(NC)"
	
	# macOS AMD64
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@echo "$(GREEN)✓ macOS AMD64 build complete$(NC)"
	
	# macOS ARM64 (Apple Silicon)
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "$(GREEN)✓ macOS ARM64 build complete$(NC)"
	
	# Windows AMD64
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "$(GREEN)✓ Windows AMD64 build complete$(NC)"
	
	@echo "$(GREEN)✓ All builds complete in $(BUILD_DIR)/$(NC)"

## Install dependencies
deps:
	@echo "$(BLUE)Installing dependencies...$(NC)"
	@go mod tidy
	@go mod download
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

## Update all dependencies to latest versions
deps-update:
	@echo "$(BLUE)Updating all dependencies to latest versions...$(NC)"
	@go get -u all
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

## Install the binary to system PATH
install: build
	@echo "$(BLUE)Installing $(BINARY_NAME) to /usr/local/bin...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ $(BINARY_NAME) installed successfully$(NC)"
	@echo "$(YELLOW)You can now use: $(BINARY_NAME) <username>$(NC)"

## Uninstall the binary from system PATH
uninstall:
	@echo "$(BLUE)Uninstalling $(BINARY_NAME)...$(NC)"
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ $(BINARY_NAME) uninstalled$(NC)"

## Run tests
test:
	@echo "$(BLUE)Running tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✓ Tests complete$(NC)"

## Run linter
lint:
	@echo "$(BLUE)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
		echo "$(GREEN)✓ Linting complete$(NC)"; \
	else \
		echo "$(YELLOW)golangci-lint not found, skipping...$(NC)"; \
		echo "$(YELLOW)Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

## Clean build artifacts
clean:
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f *.html
	@echo "$(GREEN)✓ Clean complete$(NC)"

## Run the application with example user
demo:
	@echo "$(BLUE)Running demo with user 'octocat'...$(NC)"
	@./$(BUILD_DIR)/$(BINARY_NAME) octocat

## Show version information
version:
	@echo "$(BLUE)GitHub Profiler v$(VERSION)$(NC)"
	@echo "$(BLUE)Go version: $(shell go version)$(NC)"

## Create release archives
release: build-all
	@echo "$(BLUE)Creating release archives...$(NC)"
	@mkdir -p $(BUILD_DIR)/releases
	
	@cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	@cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	@cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64
	@cd $(BUILD_DIR) && tar -czf releases/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@cd $(BUILD_DIR) && zip -q releases/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe
	
	@echo "$(GREEN)✓ Release archives created in $(BUILD_DIR)/releases/$(NC)"

## Show help
help:
	@echo "$(BLUE)GitHub Profiler - Beautiful CLI tool for GitHub user analysis$(NC)"
	@echo ""
	@echo "$(YELLOW)Available commands:$(NC)"
	@echo "  $(GREEN)build$(NC)      - Build the binary"
	@echo "  $(GREEN)build-all$(NC)  - Build for multiple platforms"
	@echo "  $(GREEN)install$(NC)    - Install binary to system PATH"
	@echo "  $(GREEN)uninstall$(NC)  - Remove binary from system PATH"
	@echo "  $(GREEN)deps$(NC)       - Install Go dependencies"
	@echo "  $(GREEN)deps-update$(NC) - Update all dependencies to latest versions"
	@echo "  $(GREEN)test$(NC)       - Run tests"
	@echo "  $(GREEN)lint$(NC)       - Run linter"
	@echo "  $(GREEN)clean$(NC)      - Clean build artifacts"
	@echo "  $(GREEN)demo$(NC)       - Run demo with 'octocat' user"
	@echo "  $(GREEN)version$(NC)    - Show version information"
	@echo "  $(GREEN)release$(NC)    - Create release archives"
	@echo "  $(GREEN)help$(NC)       - Show this help message"
	@echo ""
	@echo "$(YELLOW)Usage examples:$(NC)"
	@echo "  make build && ./build/$(BINARY_NAME) octocat"
	@echo "  make install && $(BINARY_NAME) --help"
	@echo "  $(BINARY_NAME) username --format json"
	@echo "  $(BINARY_NAME) username --format html"
	@echo "  GITHUB_TOKEN=your_token $(BINARY_NAME) username"
