# Makefile for Codex Go Bindings

# Directories
THIRD_PARTY_DIR := ../third_party
NIM_CODEX_REPO := https://github.com/codex-storage/nim-codex
NIM_CODEX_DIR := $(THIRD_PARTY_DIR)/nim-codex

.PHONY: all clean prepare build-libcodex build

# Default target
all: build

# Prepare third_party directory and clone nim-codex
prepare:
	@echo "Creating third_party directory..."
	@mkdir -p $(THIRD_PARTY_DIR)

	@echo "Cloning nim-codex repository..."
	@if [ ! -d "$(NIM_CODEX_DIR)" ]; then \
		cd $(THIRD_PARTY_DIR) && \
		git clone $(NIM_CODEX_REPO) && \
		cd $(NIM_CODEX_DIR) && \
		make update; \
	else \
		echo "nim-codex repository already exists."; \
	fi

# Build libcodex
build-libcodex: prepare
	@echo "Building libcodex..."
	@cd $(NIM_CODEX_DIR) && make libcodex

# Build Codex Go Bindings
build: build-libcodex
	@echo "Building Codex Go Bindings..."
	go build ./...

# Clean up generated files
clean:
	@echo "Cleaning up..."
	@rm -rf $(THIRD_PARTY_DIR)
	@rm -f codex-go-bindings