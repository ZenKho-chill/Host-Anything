#!/bin/bash
set -euo pipefail

EXIT_CODE=0

echo "Linting Go code..."
cd core
if ! golangci-lint run ./...; then
  EXIT_CODE=1
fi
cd ..

echo "Linting frontend..."
cd web
if ! npm run lint; then
  EXIT_CODE=1
fi
cd ..

exit $EXIT_CODE
