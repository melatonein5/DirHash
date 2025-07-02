#!/bin/bash

echo "🔍 Running Go quality checks..."

# Check formatting
echo "📝 Checking formatting..."
UNFORMATTED=$(gofmt -s -l . | grep -v installer/)
if [ -n "$UNFORMATTED" ]; then
    echo "❌ The following files need formatting:"
    echo "$UNFORMATTED"
    echo "Run: gofmt -s -w ."
    exit 1
else
    echo "✅ All files are properly formatted"
fi

# Run go vet
echo "🔍 Running go vet..."
if go vet ./src/...; then
    echo "✅ No vet issues found"
else
    echo "❌ Go vet found issues"
    exit 1
fi

# Run tests
echo "🧪 Running tests..."
if go test ./src/...; then
    echo "✅ All tests pass"
else
    echo "❌ Tests failed"
    exit 1
fi

# Generate coverage
echo "📊 Generating coverage report..."
go test -coverprofile=coverage.out ./src/...
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
echo "📈 Total coverage: $COVERAGE"

echo "🎉 All quality checks passed!"