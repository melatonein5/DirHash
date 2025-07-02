#!/bin/bash

# Generate coverage report
go test -coverprofile=coverage.out ./src/...

# Extract coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

# Create coverage badge URL
COLOR="red"
if (( $(echo "$COVERAGE > 80" | bc -l) )); then COLOR="green"; fi
if (( $(echo "$COVERAGE > 60" | bc -l) )) && (( $(echo "$COVERAGE <= 80" | bc -l) )); then COLOR="yellow"; fi

echo "Coverage: $COVERAGE%"
echo "Badge URL: https://img.shields.io/badge/Coverage-$COVERAGE%25-$COLOR"

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
echo "HTML report generated: coverage.html"