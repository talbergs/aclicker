#!/bin/bash
echo "Running Go tests recursively for all packages..."
go test ./...
if [ $? -eq 0 ]; then
    echo "All tests passed!"
else
    echo "Some tests failed!"
    exit 1
fi
