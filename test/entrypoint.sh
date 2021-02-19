#!/bin/bash

echo "Waiting for mssql to be ready..."

sleep 45s

echo "Running tests..."
go test ./... -race -coverprofile=cp.out
test_exit_code=$?

mv cp.out /tests/cp.out

if [[ $test_exit_code -ne 0 ]]; then
    echo "Tests failed!"
    exit 1
fi