# Source directory
SRC_DIR := src

# Build directory
BUILD_DIR := build

# Name of the binary output
BINARY_NAME := mediametaparser

# Operating systems and architectures
OS_ARCH := linux-amd64 linux-386 darwin-amd64 windows-amd64 windows-386

.PHONY: all $(OS_ARCH) clean

all: $(OS_ARCH)

# General rule for building each OS-ARCH combination
$(OS_ARCH):
	@if [ ! -d $(BUILD_DIR) ]; then \
		mkdir $(BUILD_DIR); \
	fi
	@(cd $(SRC_DIR); \
    GOOS=$(word 1,$(subst -, ,$@)) GOARCH=$(word 2,$(subst -, ,$@)) go build -o '../$(BUILD_DIR)/$(BINARY_NAME)-$@' main.go)

# Clean up
clean:
	@echo Cleaning up...
	@rm -rf $(BUILD_DIR)/
