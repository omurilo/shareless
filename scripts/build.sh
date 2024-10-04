#!/bin/bash

# Exit on error
set -e

# Build the Go application
echo "Building the Go application..."
go build -o bin/shareless cmd/shareless/main.go

echo "Build completed successfully!"
