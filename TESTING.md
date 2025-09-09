# Testing Guide for SJSU Badminton Discord Bot

This guide explains how to use test-driven development (TDD) to ensure code quality and reduce Railway deployment issues.

## 🧪 **Available Test Commands**

### Basic Testing
```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test -v ./internal/config
```

### Advanced Testing
```bash
# Run all tests with race detection
make test-race

# Generate test coverage report
make test-coverage

# Run all tests (basic + race + coverage)
make test-all

# Test build without running
make build-test
```

### Code Quality
```bash
# Run linter
make lint

# Fix formatting issues
make lint-fix

# Run pre-commit checks (format + lint + test + build)
make pre-commit
```

## 🔧 **Pre-Commit Workflow**

Before committing any changes, run the pre-commit script:

```bash
# Run comprehensive pre-commit checks
./scripts/pre-commit.sh

# Or use the Makefile target
make pre-commit
```

This will:
1. ✅ Format your code
2. ✅ Run the linter
3. ✅ Check code formatting
4. ✅ Run all tests
5. ✅ Test the build

## 📊 **Test Coverage**

Generate and view test coverage:

```bash
# Generate coverage report
make test-coverage

# View coverage in browser
open coverage.html
```

## 🚀 **Railway Deployment Workflow**

### 1. Local Development
```bash
# Make your changes
# ... edit code ...

# Run pre-commit checks
make pre-commit

# If all tests pass, commit
git add .
git commit -m "Your changes"
git push
```

### 2. Railway Auto-Deploy
- Railway automatically detects GitHub pushes
- Runs the same tests in CI/CD pipeline
- Only deploys if tests pass
- Reduces runtime errors

## 🐛 **Debugging Failed Tests**

### Test Failures
```bash
# Run specific test with verbose output
go test -v -run TestSpecificFunction ./internal/package

# Run tests with race detection
go test -race ./...

# Check for formatting issues
gofmt -s -l .
```

### Build Failures
```bash
# Test build
go build -v ./cmd/bot

# Check for unused imports
go vet ./...

# Check for formatting
gofmt -s -l .
```

## 📝 **Adding New Tests**

### 1. Create Test File
```bash
# For package internal/example
touch internal/example/example_test.go
```

### 2. Write Test Functions
```go
package example

import "testing"

func TestFunction(t *testing.T) {
    // Test implementation
    result := Function()
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

### 3. Run Tests
```bash
go test -v ./internal/example
```

## 🎯 **Test Categories**

### Unit Tests
- Test individual functions
- Mock external dependencies
- Fast execution
- High coverage

### Integration Tests
- Test component interactions
- Use real dependencies where possible
- Test data flow

### Configuration Tests
- Test environment variable handling
- Test default values
- Test validation logic

## 📈 **Benefits of TDD**

1. **Early Bug Detection** - Catch issues before deployment
2. **Faster Railway Deploys** - Fewer failed deployments
3. **Code Quality** - Enforced formatting and linting
4. **Confidence** - Know your code works before pushing
5. **Documentation** - Tests serve as usage examples

## 🔄 **Continuous Integration**

The GitHub Actions workflow automatically:
- Runs tests on every push
- Builds the Docker image
- Validates code quality
- Only allows merges if tests pass

## 📚 **Best Practices**

1. **Write Tests First** - TDD approach
2. **Test Edge Cases** - Error conditions, empty inputs
3. **Keep Tests Fast** - Use mocks for slow operations
4. **Test Configuration** - Environment variables, defaults
5. **Run Tests Locally** - Before every commit
6. **Use Descriptive Names** - Clear test function names

## 🚨 **Common Issues**

### Import Errors
```bash
go mod tidy
go mod download
```

### Test Failures
```bash
# Check for race conditions
go test -race ./...

# Check for timing issues
go test -timeout 30s ./...
```

### Build Failures
```bash
# Clean and rebuild
go clean -cache
go build ./cmd/bot
```

Remember: **Test locally, deploy confidently!** 🎾🤖
