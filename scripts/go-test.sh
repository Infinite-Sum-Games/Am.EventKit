#!/bin/sh

if find . -name "*_test.go" | grep . >/dev/null; then
  echo "Running tests..."
  go test ./...
else
  echo "No test files found. Skipping tests."
fi
