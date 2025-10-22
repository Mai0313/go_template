# Makefile for post_hook Go application

GO ?= go
GOPATH_BIN := $(shell $(GO) env GOPATH)/bin

ifneq (,$(wildcard .env))
include .env
endif

# Version variables
VERSION_RAW := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
VERSION := $(shell echo $(VERSION_RAW) | sed 's/^v//')
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS := -X go-template/core/version.Version=$(VERSION) \
		   -X go-template/core/version.BuildTime=$(BUILD_TIME) \
		   -X go-template/core/version.GitCommit=$(GIT_COMMIT)

# Build directory
BUILD_DIR := build
BIN_NAME := go-template
PLATFORMS := linux/amd64 linux/arm64 windows/amd64 windows/arm64 darwin/amd64 darwin/arm64

# Automatically find all command directories
CMDS := $(notdir $(wildcard cmd/*))

.PHONY: all
all: $(CMDS) ## Build all commands.

.PHONY: build
build: $(CMDS) ## Build default commands

help: # Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

$(CMDS):
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 $(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o $(BUILD_DIR)/$@ ./cmd/$@
	@echo "\033[32mSuccessfully built target: $@ (version: $(VERSION))\033[0m"

# ---

.PHONY: build-all $(addprefix build_,$(subst /,_,$(PLATFORMS)))
build-all: $(addprefix build_,$(subst /,_,$(PLATFORMS)))

# Generic cross-platform build function
define build_platform
build_$(subst /,_,$(1)):
	@mkdir -p $(BUILD_DIR)
	$(eval GOOS_VAL := $(word 1,$(subst /, ,$(1))))
	$(eval GOARCH_VAL := $(word 2,$(subst /, ,$(1))))
	$(eval SUFFIX := $(if $(filter windows,$(GOOS_VAL)),.exe,))
	@$(foreach cmd,$(CMDS),GOOS=$(GOOS_VAL) GOARCH=$(GOARCH_VAL) $(GO) build -ldflags "$(LDFLAGS) -s -w" -o $(BUILD_DIR)/$(cmd)-$(GOOS_VAL)-$(GOARCH_VAL)$(SUFFIX) ./cmd/$(cmd);)
	@echo "\033[32mSuccessfully built for $(1)\033[0m"
endef

# Generate build rules for each platform
$(foreach platform,$(PLATFORMS),$(eval $(call build_platform,$(platform))))

# ---

.PHONY: package-all $(addprefix package_,$(subst /,_,$(PLATFORMS)))
package-all: build-all $(addprefix package_,$(subst /,_,$(PLATFORMS))) ## Build and package all platforms into zip files

# Generic cross-platform package function
define package_platform
package_$(subst /,_,$(1)):
	$(eval GOOS_VAL := $(word 1,$(subst /, ,$(1))))
	$(eval GOARCH_VAL := $(word 2,$(subst /, ,$(1))))
	$(eval SUFFIX := $(if $(filter windows,$(GOOS_VAL)),.exe,))
	@$(foreach cmd,$(CMDS),cd $(BUILD_DIR) && zip -qm $(cmd)-$(GOOS_VAL)-$(GOARCH_VAL).zip $(cmd)-$(GOOS_VAL)-$(GOARCH_VAL)$(SUFFIX) && cd ..;)
	@echo "\033[32mSuccessfully packaged for $(1)\033[0m"
endef

# Generate package rules for each platform
$(foreach platform,$(PLATFORMS),$(eval $(call package_platform,$(platform))))

# ---

# Clean build artifacts
.PHONY: clean
clean: ## Remove build artifacts
	@rm -rf $(BUILD_DIR) coverage.out
	@find . -type f -name "*.DS_Store" -ls -delete
	@find . -type f -name "*.zip" -ls -delete
	@$(GO) clean -cache
	@$(GO) clean -testcache
	@$(GO) clean -fuzzcache
	@git fetch --prune
	@git gc --prune=now --aggressive

# ---

# Run the application (for testing)
.PHONY: run
run: build ## Build and run the application
	./$(BUILD_DIR)/$(BIN_NAME)

# ---

# Install to system (optional)
.PHONY: install
install: build ## Install binary to /usr/local/bin
	sudo cp $(BUILD_DIR)/$(BIN_NAME) /usr/local/bin/$(BIN_NAME)

# ---

# Format code
.PHONY: fmt
fmt: ## Format Go code
	$(GO) fmt ./...

# ---

# Test (if you add tests later)
.PHONY: test
test: ## Run tests
	$(GO) test -cover -coverprofile=coverage.out ./...

# ---

.PHONY: test-verbose
test-verbose: ## Run verbose tests
	$(GO) test -cover -coverprofile=coverage.out -v ./...

# ---

# Dead code detection using staticcheck (U1000) and golang.org/x/tools deadcode
.PHONY: lint-deadcode
lint-deadcode: ## Detect dead/unused code (staticcheck U1000 + deadcode)
	@{ test -x "$(GOPATH_BIN)/staticcheck" || { echo "Installing staticcheck..."; $(GO) install honnef.co/go/tools/cmd/staticcheck@latest; }; }
	@{ test -x "$(GOPATH_BIN)/deadcode" || { echo "Installing deadcode..."; $(GO) install golang.org/x/tools/cmd/deadcode@latest; }; }
	@echo "Running staticcheck (U1000) ..."
	@"$(GOPATH_BIN)/staticcheck" -checks U1000 ./... || true
	@echo "\nRunning deadcode ..."
	@"$(GOPATH_BIN)/deadcode" ./... || true
