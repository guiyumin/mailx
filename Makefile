.PHONY: build clean

APP_NAME := maily
BUILD_DIR := build
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -s -w \
	-X maily/internal/version.Version=$(VERSION) \
	-X maily/internal/version.Commit=$(COMMIT) \
	-X maily/internal/version.Date=$(DATE)

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) .


clean:
	rm -rf $(BUILD_DIR)


.PHONY: push
push:
	git push origin main --tags
