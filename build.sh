#!/bin/bash
echo "Building Clicker2..."
go build -o clicker2 .
if [ $? -eq 0 ]; then
    echo "Build successful!"
else
    echo "Build failed!"
    exit 1
fi
