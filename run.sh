#!/bin/bash

# Start the WaWeb V2 server
echo "Starting WaWeb V2 Server..."

# Check if .env exists, if not copy from example
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
fi

# Build and run
go run main.go
