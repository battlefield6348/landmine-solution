# Binary names
BINARY_NAME=solver
WASM_NAME=main.wasm

# Build paths
CLI_SRC=cmd/solver/main.go
WASM_SRC=cmd/wasm/main.go
DIST_DIR=docs

.PHONY: all build build-cli build-wasm clean help

all: build

## build: Build both CLI and WASM versions
build: build-cli build-wasm

## build-cli: Build the CLI solver
build-cli:
	@echo "Building CLI solver..."
	go build -o $(BINARY_NAME) $(CLI_SRC)

## build-wasm: Build the WebAssembly version
build-wasm:
	@echo "Building WebAssembly..."
	GOOS=js GOARCH=wasm go build -o $(DIST_DIR)/$(WASM_NAME) $(WASM_SRC)

## clean: Remove build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -f $(DIST_DIR)/$(WASM_NAME)

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
