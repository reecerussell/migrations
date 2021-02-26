#!/bin/bash

echo "Waiting for mssql to be ready..."

sleep 45s

echo "Running tests..."
go test ./... -timeout 30s -race -coverprofile=coverage.out -covermode=atomic
test_exit_code=$?

mkdir -p /tests
mv cp.out /tests/cp.out

if [[ $test_exit_code -ne 0 ]]; then
    echo "Tests failed!"
    exit 1
fi

echo "Tests passed!"