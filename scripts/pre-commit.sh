#!/bin/bash

# Pre-commit hook for SJSU Badminton Bot
# This script runs before each commit to ensure code quality

set -e

echo "🔍 Running pre-commit checks..."

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: go.mod not found. Please run this script from the project root."
    exit 1
fi

# Format code
echo "📝 Formatting code..."
go fmt ./...

# Run linter
echo "�� Running linter..."
go vet ./...

# Check for formatting issues
echo "🎨 Checking code formatting..."
if [ "$(gofmt -s -l . | wc -l)" -gt 0 ]; then
    echo "❌ Code formatting issues found:"
    gofmt -s -l .
    echo "Run 'make lint-fix' to fix formatting issues."
    exit 1
fi

# Run tests
echo "🧪 Running tests..."
go test -v ./...

# Build test
echo "🔨 Testing build..."
go build -v ./cmd/bot

echo "✅ All pre-commit checks passed!"
echo "🚀 Ready to commit!"
