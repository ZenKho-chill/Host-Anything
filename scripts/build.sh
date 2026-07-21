#!/bin/bash
set -euo pipefail

echo "Building Go binary..."
cd core
GOOS=linux GOARCH=amd64 go build -o ../bin/hostanything-linux-amd64 ./cmd/hostanything
GOOS=linux GOARCH=arm64 go build -o ../bin/hostanything-linux-arm64 ./cmd/hostanything
cd ..

echo "Building frontend..."
cd web
npm install
npm run build
cd ..

echo "Build successful!"
