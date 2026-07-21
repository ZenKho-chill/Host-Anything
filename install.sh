#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}==> Setting up Host Anything...${NC}"

# 1. Build the Go Backend
echo -e "\n${BLUE}--> Building core backend...${NC}"
make build

# 2. Build the Web UI
echo -e "\n${BLUE}--> Building web UI...${NC}"
make build-web

echo -e "\n${GREEN}==> Setup complete!${NC}"
echo -e "${GREEN}==> Starting Host Anything on port 8080...${NC}"
echo "Press Ctrl+C to stop."
echo ""

# Start the binary
./bin/hostanything
