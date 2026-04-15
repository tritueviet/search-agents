# Search Agents Makefile
.PHONY: help install build test clean run-tor stop-tor lint

# Variables
BINARY_NAME=sagent
API_BINARY_NAME=sagent-api
BUILD_DIR=build
GO=go
GOPROXY=https://proxy.golang.org,direct

# Build flags
LDFLAGS=-ldflags "-s -w"
GOFLAGS=-trimpath

help: ## Show this help message
	@echo "Search Agents - Makefile Commands"
	@echo "=================================="
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install Tor and required system tools
	@echo "🔧 Installing required tools..."
	@echo "Installing Tor proxy..."
	echo '1122' | sudo -S apt update -qq || true
	echo '1122' | sudo -S apt install -y tor jq
	@echo "Configuring Tor..."
	echo 'SocksPort 9050' | sudo tee -a /etc/tor/torrc > /dev/null 2>&1 || true
	@echo "Starting Tor service..."
	echo '1122' | sudo -S systemctl enable tor
	echo '1122' | sudo -S systemctl start tor
	@sleep 3
	@echo "Verifying Tor is running on port 9050..."
	@ss -tlnp | grep -q 9050 && echo "✅ Tor is running on port 9050" || echo "⚠️  Tor not started automatically, run 'make run-tor'"
	@echo ""
	@echo "✅ Installation complete!"
	@echo "Usage: ./sagent images \"query\" -m 5"

build: ## Build CLI and API binaries
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/sagent
	@echo "🔨 Building $(API_BINARY_NAME)..."
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(API_BINARY_NAME) ./cmd/server
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME), $(BUILD_DIR)/$(API_BINARY_NAME)"

build-cli: ## Build CLI only
	@echo "🔨 Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/sagent
	@echo "✅ CLI built: $(BINARY_NAME)"

build-api: ## Build API server only
	@echo "🔨 Building $(API_BINARY_NAME)..."
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(API_BINARY_NAME) ./cmd/server
	@echo "✅ API built: $(API_BINARY_NAME)"

test: ## Run all tests
	@echo "🧪 Running tests..."
	$(GO) test -v -race -coverprofile=coverage.out ./...
	@echo ""
	@echo "📊 Coverage report:"
	$(GO) tool cover -func=coverage.out | tail -1

test-integration: ## Run integration tests only
	@echo "🧪 Running integration tests..."
	$(GO) test -v -timeout 120s ./tests/...

test-unit: ## Run unit tests only (excluding integration)
	@echo "🧪 Running unit tests..."
	$(GO) test -v -race ./internal/... ./pkg/...

test-text: ## Test text search only
	@echo "🧪 Testing text search..."
	$(GO) test -v -run TestTextSearch ./tests/...

test-images: ## Test images search (requires Tor)
	@echo "🧪 Testing images search..."
	$(GO) test -v -run TestImagesSearch ./tests/...

test-videos: ## Test videos search (requires Tor)
	@echo "🧪 Testing videos search..."
	$(GO) test -v -run TestVideosSearch ./tests/...

test-news: ## Test news search (requires Tor)
	@echo "🧪 Testing news search..."
	$(GO) test -v -run TestNewsSearch ./tests/...

test-books: ## Test books search (requires Tor)
	@echo "🧪 Testing books search..."
	$(GO) test -v -run TestBooksSearch ./tests/...

lint: ## Run linter
	@echo "🔍 Running linter..."
	$(GO) vet ./...
	@echo "✅ Linting complete"

run-tor: ## Start Tor service
	@echo "🚀 Starting Tor service..."
	echo '1122' | sudo -S systemctl start tor
	@sleep 2
	@ss -tlnp | grep -q 9050 && echo "✅ Tor started on port 9050" || echo "❌ Failed to start Tor"

stop-tor: ## Stop Tor service
	@echo "🛑 Stopping Tor service..."
	echo '1122' | sudo -S systemctl stop tor
	@echo "✅ Tor stopped"

status-tor: ## Check Tor status
	@echo "📊 Tor Status:"
	@ss -tlnp | grep 9050 && echo "✅ Tor is running on port 9050" || echo "❌ Tor is not running"

run: build-cli ## Build and run CLI with example
	@echo "🚀 Running $(BINARY_NAME)..."
	./$(BINARY_NAME) text "golang tutorial" -m 3

run-api: build-api ## Build and run API server
	@echo "🚀 Starting API server on port 8000..."
	./$(API_BINARY_NAME) --host 0.0.0.0 --port 8000

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME) $(API_BINARY_NAME)
	rm -f coverage.out
	$(GO) clean -cache
	@echo "✅ Clean complete"

deps: ## Download and verify dependencies
	@echo "📦 Downloading dependencies..."
	$(GO) mod download
	$(GO) mod verify
	$(GO) mod tidy
	@echo "✅ Dependencies updated"
