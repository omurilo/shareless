.PHONY: all build run docker-run docker-down clean watch docker-build

# Get the current git commit hash
COMMIT_HASH := $(shell git rev-parse --short HEAD)

# Build the Go application
build:
	@echo "Building the Go application..."
	@bash scripts/build.sh

# Run the application
run:
	@go run cmd/shareless/main.go

# Create DB container
docker-run:
	@if docker compose -f deploy/docker-compose.yaml up --build >/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose -f deploy/docker-compose.yaml down >/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf bin/shareless

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

# Build the Docker image with the commit hash as the tag
docker-build:
	@echo "Building the Docker image with tag: omurilo/shareless:$(COMMIT_HASH)..."
	@docker build -t omurilo/shareless:$(COMMIT_HASH) -f deploy/Dockerfile .
