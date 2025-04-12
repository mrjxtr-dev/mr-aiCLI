.PHONY: all build build-all clean

BINARY_NAME := mrai
OUTPUT_DIR := bin

# OS/Arch combinations
TARGETS := \
	linux/amd64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

# Default target
all: build-all

# Build for current system
build:
	@echo "Building for current system..."
	@mkdir -p $(OUTPUT_DIR)
	@go build -o $(OUTPUT_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete: $(OUTPUT_DIR)/$(BINARY_NAME)"

# Build for all targets
build-all:
	@echo "Cross-compiling..."
	@mkdir -p $(OUTPUT_DIR)
	@$(foreach target,$(TARGETS), \
		GOOS=$(word 1,$(subst /, ,$(target))) \
		GOARCH=$(word 2,$(subst /, ,$(target))) \
		go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-$(subst /,-,$(target))$(if $(findstring windows,$(target)),.exe,) main.go &&) true
	@echo "All builds complete."

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -rf $(OUTPUT_DIR)
	@echo "Clean complete."

