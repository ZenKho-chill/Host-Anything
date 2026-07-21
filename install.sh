#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}==> Host Anything Enterprise Setup${NC}"

# 1. Check for Sudo/Root privileges
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Error: Please run this script with sudo or as root.${NC}"
  echo -e "Example: sudo bash install.sh"
  exit 1
fi

# 2. Dependency Scanner
MISSING_DEPS=()
command -v go >/dev/null 2>&1 || MISSING_DEPS+=("golang")
command -v node >/dev/null 2>&1 || MISSING_DEPS+=("nodejs")
command -v docker >/dev/null 2>&1 || MISSING_DEPS+=("docker")

if [ ${#MISSING_DEPS[@]} -ne 0 ]; then
    echo -e "${YELLOW}Missing dependencies detected: ${MISSING_DEPS[*]}${NC}"
    read -p "Do you want to attempt automatic installation? (y/n) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${BLUE}--> Updating apt packages...${NC}"
        apt-get update -y
        
        for dep in "${MISSING_DEPS[@]}"; do
            echo -e "${BLUE}--> Installing $dep...${NC}"
            if [ "$dep" == "nodejs" ]; then
                # Install Node.js 20.x via NodeSource
                curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
                apt-get install -y nodejs
            elif [ "$dep" == "docker" ]; then
                apt-get install -y docker.io
                systemctl enable --now docker
            elif [ "$dep" == "golang" ]; then
                apt-get install -y golang
            else
                apt-get install -y "$dep"
            fi
        done
        echo -e "${GREEN}--> Dependencies installed!${NC}"
    else
        echo -e "${RED}Setup aborted. Please install the missing dependencies manually and run this script again.${NC}"
        exit 1
    fi
fi

# 3. Super Admin Credentials Setup
echo -e "\n${BLUE}==> Super Admin Account Setup${NC}"
echo -e "These credentials will be used to initialize the master database."

read -p "Admin Username (default: admin): " ADMIN_USER
ADMIN_USER=${ADMIN_USER:-admin}

# Read password silently
while true; do
    read -s -p "Admin Password: " ADMIN_PASS
    echo ""
    if [ -z "$ADMIN_PASS" ]; then
        echo -e "${RED}Password cannot be empty. Try again.${NC}"
    else
        break
    fi
done

export HA_ADMIN_USERNAME="$ADMIN_USER"
export HA_ADMIN_PASSWORD="$ADMIN_PASS"
export HA_DB_PATH="/var/lib/hostanything/data/hostanything.db"

# Ensure data directory exists
mkdir -p /var/lib/hostanything/data
chown -R $SUDO_USER:$SUDO_USER /var/lib/hostanything/data || true

# 4. Build the Go Backend
echo -e "\n${BLUE}--> Building core backend...${NC}"
# Drop privileges to build if running via sudo
if [ -n "$SUDO_USER" ]; then
    sudo -u $SUDO_USER -E make build
else
    make build
fi

# 5. Build the Web UI
echo -e "\n${BLUE}--> Building web UI...${NC}"
if [ -n "$SUDO_USER" ]; then
    sudo -u $SUDO_USER -E make build-web
else
    make build-web
fi

echo -e "\n${GREEN}==> Setup complete!${NC}"
echo -e "${GREEN}==> Starting Host Anything on port 8080...${NC}"
echo "Press Ctrl+C to stop."
echo ""

# Start the binary
if [ -n "$SUDO_USER" ]; then
    sudo -u $SUDO_USER -E ./bin/hostanything
else
    ./bin/hostanything
fi
