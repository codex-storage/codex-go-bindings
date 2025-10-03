# Makefile for Codex Go Bindings

NIM_CODEX_DIR := vendor/nim-codex

.PHONY: all clean update build-libcodex build

all: build

update:
	@echo "Fetching submodules..."
	@git submodule update --init --recursive
	@echo "Updating nim-codex..."
	@$(MAKE) -C $(NIM_CODEX_DIR) update

libcodex:
	@echo "Building libcodex..."
	@$(MAKE) -C $(NIM_CODEX_DIR) libcodex

build:
	@echo "Building Codex Go Bindings..."
	go build ./...

clean:
	@echo "Cleaning up..."
	@git submodule deinit -f $(NIM_CODEX_DIR)
	@rm -f codex-go-bindings