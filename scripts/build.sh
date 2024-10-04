#!/bin/bash

# Exit on error
set -e

# Build the Go application
echo "Building the Go application..."
mkdir -p bin/shareless
go build -o bin/shareless ./...

echo "Build completed successfully!"
