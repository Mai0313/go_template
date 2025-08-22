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

# Packaging to Go-Template-{platform}.zip
.PHONY: package-all package_linux_amd64 package_linux_arm64 package_windows_amd64 package_windows_arm64 package_darwin_amd64 package_darwin_arm64
package-all: build-all package_linux_amd64 package_linux_arm64 package_windows_amd64 package_windows_arm64 package_darwin_amd64 package_darwin_arm64

package_linux_amd64: build_linux_amd64
	@mkdir -p $(BUILD_DIR)
	@cd $(BUILD_DIR) && cp ../README*.md . && \
	  zip -q -9 -r "Go-Template-linux-amd64.zip" $(BIN_NAME)-linux-amd64 $(CLI_NAME)-linux-amd64 README*.md && rm -f README*.md

package_linux_arm64: build_linux_arm64
	@mkdir -p $(BUILD_DIR)
	@cd $(BUILD_DIR) && cp ../README*.md . && \
	  zip -q -9 -r "Go-Template-linux-arm64.zip" $(BIN_NAME)-linux-arm64 $(CLI_NAME)-linux-arm64 README*.md && rm -f README*.md

package_windows_amd64: build_windows_amd64
	@mkdir -p $(BUILD_DIR)
	@cd $(BUILD_DIR) && cp ../README*.md . && \
	  zip -q -9 -r "Go-Template-windows-amd64.zip" $(BIN_NAME)-windows-amd64.exe $(CLI_NAME)-windows-amd64.exe README*.md && rm -f README*.md

package_windows_arm64: build_windows_arm64
	@mkdir -p $(BUILD_DIR)
	@cd $(BUILD_DIR) && cp ../README*.md . && \
	  zip -q -9 -r "Go-Template-windows-arm64.zip" $(BIN_NAME)-windows-arm64.exe $(CLI_NAME)-windows-arm64.exe README*.md && rm -f README*.md

package_darwin_amd64: build_darwin_amd64
	@mkdir -p $(BUILD_DIR)
	@cd $(BUILD_DIR) && cp ../README*.md . && \
	  zip -q -9 -r "Go-Template-darwin-amd64.zip" $(BIN_NAME)-darwin-amd64 $(CLI_NAME)-darwin-amd64 README*.md && rm -f README*.md

package_darwin_arm64: build_darwin_arm64
	@mkdir -p $(BUILD_DIR)
	@cd $(BUILD_DIR) && cp ../README*.md . && \
	  zip -q -9 -r "Go-Template-darwin-arm64.zip" $(BIN_NAME)-darwin-arm64 $(CLI_NAME)-darwin-arm64 README*.md && rm -f README*.md

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

# Test (if you add tests later)
.PHONY: test
test: ## Run tests
	$(GO) test -cover -v ./...

.PHONY: test-verbose
test-verbose: ## Run verbose tests
	$(GO) test -cover -v ./tests -run TestParser_FromTestConversationJSONL_PrintsFullPayload -count=1
