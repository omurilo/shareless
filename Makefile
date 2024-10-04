.PHONY: build clean docker

# Get the current git commit hash
COMMIT_HASH := $(shell git rev-parse --short HEAD)

# Build the Go application
build:
	@echo "Building the Go application..."
	@bash scripts/build.sh

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f bin/shareless

# Build the Docker image with the commit hash as the tag
docker:
	@echo "Building the Docker image with tag: @omurilo/shareless:$(COMMIT_HASH)..."
	@docker build -t @omurilo/shareless:$(COMMIT_HASH) .
