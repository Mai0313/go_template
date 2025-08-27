# Makefile for Go Template Project

GO ?= go

# Build directory
BUILD_DIR := build
BIN_NAME := app
CLI_NAME := cli

# Automatically find all command directories
CMDS := $(notdir $(wildcard cmd/*))

.PHONY: all
all: $(CMDS) ## Build all commands.

help: # Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

$(CMDS):
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GO) build -v -ldflags '-s -w' -o $(BUILD_DIR)/$@ ./cmd/$@
	@echo "\033[32mSuccessfully built target: $@\033[0m"

# Build for multiple platforms
.PHONY: build-all build_linux_amd64 build_linux_arm64 build_windows_amd64 build_windows_arm64 build_darwin_amd64 build_darwin_arm64
build-all: build_linux_amd64 build_linux_arm64 build_windows_amd64 build_windows_arm64 build_darwin_amd64 build_darwin_arm64

build_linux_amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-linux-amd64 ./cmd/app
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(CLI_NAME)-linux-amd64 ./cmd/cli

build_linux_arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-linux-arm64 ./cmd/app
	GOOS=linux GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(CLI_NAME)-linux-arm64 ./cmd/cli

build_windows_amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-windows-amd64.exe ./cmd/app
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(CLI_NAME)-windows-amd64.exe ./cmd/cli

build_windows_arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-windows-arm64.exe ./cmd/app
	GOOS=windows GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(CLI_NAME)-windows-arm64.exe ./cmd/cli

build_darwin_amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-darwin-amd64 ./cmd/app
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(CLI_NAME)-darwin-amd64 ./cmd/cli

build_darwin_arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-darwin-arm64 ./cmd/app
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(CLI_NAME)-darwin-arm64 ./cmd/cli

# Clean build artifacts
.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

# Run the application (for testing)
.PHONY: run
run: app ## Build and run the application
	./$(BUILD_DIR)/app

# Install to system (optional)
.PHONY: install
install: app ## Install binary to /usr/local/bin
	sudo cp $(BUILD_DIR)/app /usr/local/bin/app

# Format code
.PHONY: fmt
fmt: ## Format Go code
	$(GO) fmt ./...

# Test (runs all test types with coverage and reports)
.PHONY: test
test: ## Run all tests with coverage and generate reports
	@echo "🚀 Starting comprehensive test suite..."
	@echo "Installing go-junit-report if not present..."
	@$(GO) install github.com/jstemmer/go-junit-report/v2@latest
	@echo "Creating reports directory..."
	@mkdir -p .github/reports

	@echo "📊 Running unit tests with coverage..."
	$(GO) test -coverprofile=.github/reports/coverage.out -covermode=atomic -v ./... 2>&1 | $(shell go env GOPATH)/bin/go-junit-report -out .github/reports/test-results.xml
	$(GO) tool cover -html=.github/reports/coverage.out -o .github/reports/coverage.html

	@echo "🔗 Running integration tests..."
	$(GO) test -v ./tests/integration/

	@echo "⚡ Running benchmark tests..."
	$(GO) test -bench=. -benchmem ./tests/benchmark/ -run=^$$

	@echo "🎯 Running fuzz tests..."
	@echo "  - Testing FuzzConfigParsing..."
	@$(GO) test -fuzz=FuzzConfigParsing -fuzztime=3s ./tests/fuzz/ -run=^$$ || true
	@echo "  - Testing FuzzConfigLoad..."
	@$(GO) test -fuzz=FuzzConfigLoad -fuzztime=3s ./tests/fuzz/ -run=^$$ || true
	@echo "  - Testing FuzzConfigPath..."
	@$(GO) test -fuzz=FuzzConfigPath -fuzztime=2s ./tests/fuzz/ -run=^$$ || true
	@echo "  - Testing FuzzJSONFieldValues..."
	@$(GO) test -fuzz=FuzzJSONFieldValues -fuzztime=2s ./tests/fuzz/ -run=^$$ || true

	@echo "🔍 Running code linting..."
	$(GO) vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping additional linting"; \
	fi

	@echo "✅ All tests completed!"
	@echo "📄 XML test report: .github/reports/test-results.xml"
	@echo "📈 Coverage report: .github/reports/coverage.html"
