# Makefile for post_hook Go application

GO ?= go

# Build directory
BUILD_DIR := build
BIN_NAME := claude_analysis
INSTALLER_NAME := installer
NODE_WIN_ZIP := node-v22.18.0-win-x64.zip
NODE_WIN_URL := https://nodejs.org/dist/v22.18.0/$(NODE_WIN_ZIP)
NODE_WIN_ARM64_ZIP := node-v22.18.0-win-arm64.zip
NODE_WIN_ARM64_URL := https://nodejs.org/dist/v22.18.0/$(NODE_WIN_ARM64_ZIP)
NODE_LINUX_AMD64_TXZ := node-v22.18.0-linux-x64.tar.xz
NODE_LINUX_AMD64_URL := https://nodejs.org/dist/v22.18.0/$(NODE_LINUX_AMD64_TXZ)
NODE_LINUX_ARM64_TXZ := node-v22.18.0-linux-arm64.tar.xz
NODE_LINUX_ARM64_URL := https://nodejs.org/dist/v22.18.0/$(NODE_LINUX_ARM64_TXZ)
NODE_DARWIN_AMD64_TGZ := node-v22.18.0-darwin-x64.tar.gz
NODE_DARWIN_AMD64_URL := https://nodejs.org/dist/v22.18.0/$(NODE_DARWIN_AMD64_TGZ)
NODE_DARWIN_ARM64_TGZ := node-v22.18.0-darwin-arm64.tar.gz
NODE_DARWIN_ARM64_URL := https://nodejs.org/dist/v22.18.0/$(NODE_DARWIN_ARM64_TGZ)

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
	CGO_ENABLED=0 $(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o $(BUILD_DIR)/$@ ./cmd/$@
	@echo "\033[32mSuccessfully built target: $@\033[0m"

# Build for multiple platforms
.PHONY: build-all build_linux_amd64 build_linux_arm64 build_windows_amd64 build_windows_arm64 build_darwin_amd64 build_darwin_arm64
build-all: build_linux_amd64 build_linux_arm64 build_windows_amd64 build_windows_arm64 build_darwin_amd64 build_darwin_arm64

build_linux_amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-linux-amd64 ./cmd/claude_analysis
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(INSTALLER_NAME)-linux-amd64 ./cmd/installer

build_linux_arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-linux-arm64 ./cmd/claude_analysis
	GOOS=linux GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(INSTALLER_NAME)-linux-arm64 ./cmd/installer

build_windows_amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-windows-amd64.exe ./cmd/claude_analysis
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(INSTALLER_NAME)-windows-amd64.exe ./cmd/installer

build_windows_arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-windows-arm64.exe ./cmd/claude_analysis
	GOOS=windows GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(INSTALLER_NAME)-windows-arm64.exe ./cmd/installer

build_darwin_amd64:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-darwin-amd64 ./cmd/claude_analysis
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(INSTALLER_NAME)-darwin-amd64 ./cmd/installer

build_darwin_arm64:
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(BIN_NAME)-darwin-arm64 ./cmd/claude_analysis
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_DIR)/$(INSTALLER_NAME)-darwin-arm64 ./cmd/installer

# Packaging to Claude-Code-Installer-{platform}.zip
.PHONY: package-all package_linux_amd64 package_linux_arm64 package_windows_amd64 package_windows_arm64 package_darwin_amd64 package_darwin_arm64
package-all: build-all package_linux_amd64 package_linux_arm64 package_windows_amd64 package_windows_arm64 package_darwin_amd64 package_darwin_arm64

package_linux_amd64: build_linux_amd64
	@cp $(BUILD_DIR)/$(BIN_NAME)-linux-amd64 $(BUILD_DIR)/claude_analysis
	@cp $(BUILD_DIR)/$(INSTALLER_NAME)-linux-amd64 $(BUILD_DIR)/installer
	@mkdir -p $(BUILD_DIR)
	@echo "Downloading $(NODE_LINUX_AMD64_TXZ) ..."
	@curl -fSL -o $(BUILD_DIR)/$(NODE_LINUX_AMD64_TXZ) $(NODE_LINUX_AMD64_URL) --silent
	@cd $(BUILD_DIR) && cp ../README*.md . && cp -r ../images . && \
	  zip -q -9 -r "Claude-Code-Installer-linux-amd64.zip" claude_analysis installer $(NODE_LINUX_AMD64_TXZ) README*.md images && rm -f claude_analysis installer README*.md $(NODE_LINUX_AMD64_TXZ) && rm -rf images
	@rm -f $(BUILD_DIR)/$(INSTALLER_NAME)-linux-amd64

package_linux_arm64: build_linux_arm64
	@cp $(BUILD_DIR)/$(BIN_NAME)-linux-arm64 $(BUILD_DIR)/claude_analysis
	@cp $(BUILD_DIR)/$(INSTALLER_NAME)-linux-arm64 $(BUILD_DIR)/installer
	@mkdir -p $(BUILD_DIR)
	@echo "Downloading $(NODE_LINUX_ARM64_TXZ) ..."
	@curl -fSL -o $(BUILD_DIR)/$(NODE_LINUX_ARM64_TXZ) $(NODE_LINUX_ARM64_URL) --silent
	@cd $(BUILD_DIR) && cp ../README*.md . && cp -r ../images . && \
	  zip -q -9 -r "Claude-Code-Installer-linux-arm64.zip" claude_analysis installer $(NODE_LINUX_ARM64_TXZ) README*.md images && rm -f claude_analysis installer README*.md $(NODE_LINUX_ARM64_TXZ) && rm -rf images
	@rm -f $(BUILD_DIR)/$(INSTALLER_NAME)-linux-arm64

package_windows_amd64: build_windows_amd64
	@cp $(BUILD_DIR)/$(BIN_NAME)-windows-amd64.exe $(BUILD_DIR)/claude_analysis.exe
	@cp $(BUILD_DIR)/$(INSTALLER_NAME)-windows-amd64.exe $(BUILD_DIR)/installer.exe
	@mkdir -p $(BUILD_DIR)
	@echo "Downloading $(NODE_WIN_ZIP) ..."
	@curl -fSL -o $(BUILD_DIR)/$(NODE_WIN_ZIP) $(NODE_WIN_URL) --silent
	@cd $(BUILD_DIR) && cp ../README*.md . && cp -r ../images . && \
	  zip -q -9 -r "Claude-Code-Installer-windows-amd64.zip" claude_analysis.exe installer.exe $(NODE_WIN_ZIP) README*.md images && rm -f claude_analysis.exe installer.exe README*.md $(NODE_WIN_ZIP) && rm -rf images
	@rm -f $(BUILD_DIR)/$(INSTALLER_NAME)-windows-amd64.exe

package_windows_arm64: build_windows_arm64
	@cp $(BUILD_DIR)/$(BIN_NAME)-windows-arm64.exe $(BUILD_DIR)/claude_analysis.exe
	@cp $(BUILD_DIR)/$(INSTALLER_NAME)-windows-arm64.exe $(BUILD_DIR)/installer.exe
	@mkdir -p $(BUILD_DIR)
	@echo "Downloading $(NODE_WIN_ARM64_ZIP) ..."
	@curl -fSL -o $(BUILD_DIR)/$(NODE_WIN_ARM64_ZIP) $(NODE_WIN_ARM64_URL) --silent
	@cd $(BUILD_DIR) && cp ../README*.md . && cp -r ../images . && \
	  zip -q -9 -r "Claude-Code-Installer-windows-arm64.zip" claude_analysis.exe installer.exe $(NODE_WIN_ARM64_ZIP) README*.md images && rm -f claude_analysis.exe installer.exe README*.md $(NODE_WIN_ARM64_ZIP) && rm -rf images
	@rm -f $(BUILD_DIR)/$(INSTALLER_NAME)-windows-arm64.exe

package_darwin_amd64: build_darwin_amd64
	@cp $(BUILD_DIR)/$(BIN_NAME)-darwin-amd64 $(BUILD_DIR)/claude_analysis
	@cp $(BUILD_DIR)/$(INSTALLER_NAME)-darwin-amd64 $(BUILD_DIR)/installer
	@mkdir -p $(BUILD_DIR)
	@echo "Downloading $(NODE_DARWIN_AMD64_TGZ) ..."
	@curl -fSL -o $(BUILD_DIR)/$(NODE_DARWIN_AMD64_TGZ) $(NODE_DARWIN_AMD64_URL) --silent
	@cd $(BUILD_DIR) && cp ../README*.md . && cp -r ../images . && \
	  zip -q -9 -r "Claude-Code-Installer-darwin-amd64.zip" claude_analysis installer $(NODE_DARWIN_AMD64_TGZ) README*.md images && rm -f claude_analysis installer README*.md $(NODE_DARWIN_AMD64_TGZ) && rm -rf images
	@rm -f $(BUILD_DIR)/$(INSTALLER_NAME)-darwin-amd64

package_darwin_arm64: build_darwin_arm64
	@cp $(BUILD_DIR)/$(BIN_NAME)-darwin-arm64 $(BUILD_DIR)/claude_analysis
	@cp $(BUILD_DIR)/$(INSTALLER_NAME)-darwin-arm64 $(BUILD_DIR)/installer
	@mkdir -p $(BUILD_DIR)
	@echo "Downloading $(NODE_DARWIN_ARM64_TGZ) ..."
	@curl -fSL -o $(BUILD_DIR)/$(NODE_DARWIN_ARM64_TGZ) $(NODE_DARWIN_ARM64_URL) --silent
	@cd $(BUILD_DIR) && cp ../README*.md . && cp -r ../images . && \
	  zip -q -9 -r "Claude-Code-Installer-darwin-arm64.zip" claude_analysis installer $(NODE_DARWIN_ARM64_TGZ) README*.md images && rm -f claude_analysis installer README*.md $(NODE_DARWIN_ARM64_TGZ) && rm -rf images
	@rm -f $(BUILD_DIR)/$(INSTALLER_NAME)-darwin-arm64

# Clean build artifacts
.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

# Run the application (for testing)
.PHONY: run
run: ## Build and run the application
	run: build
	./$(BUILD_DIR)/$(BIN_NAME)

# Install to system (optional)
.PHONY: install
install: ## Install binary to /usr/local/bin
	install: build
	sudo cp $(BUILD_DIR)/$(BIN_NAME) /usr/local/bin/$(BIN_NAME)

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
