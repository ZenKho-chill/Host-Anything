#!/bin/bash
set -euo pipefail

echo "Starting backend in dev mode..."
cd core
go run -race ./cmd/hostanything &
BACKEND_PID=$!
cd ..

echo "Starting frontend dev server..."
cd web
npm run dev &
FRONTEND_PID=$!
cd ..

trap 'echo "Shutting down..."; kill $BACKEND_PID $FRONTEND_PID; exit' SIGINT SIGTERM

wait
