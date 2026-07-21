#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

LOG_FILE="/tmp/ha_install.log"
> "$LOG_FILE"

echo -e "${BLUE}==> Host Anything Enterprise Setup${NC}"

# 1. Check for Sudo/Root privileges
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Error: Please run this script with sudo or as root.${NC}"
  echo -e "Example: sudo bash install.sh"
  exit 1
fi

spinner() {
    local pid=$1
    local msg=$2
    local spin='-\|/'
    local i=0
    while kill -0 $pid 2>/dev/null; do
        i=$(( (i + 1) % 4 ))
        printf "\r\033[K[${YELLOW}${spin:$i:1}${NC}] $msg"
        sleep 0.1
    done
    wait $pid
    local status=$?
    if [ $status -eq 0 ]; then
        printf "\r\033[K[${GREEN}✔${NC}] $msg\n"
    else
        printf "\r\033[K[${RED}✘${NC}] $msg (failed, check $LOG_FILE)\n"
        exit 1
    fi
}

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
        # === PHASE 1: DOWNLOAD (PARALLEL) ===
        echo -e "${BLUE}--> Phase 1: Downloading dependencies (Parallel)...${NC}"
        
        (
            # Update apt first
            apt-get update -y >> "$LOG_FILE" 2>&1
            
            # Start downloads in background
            pids=()
            for dep in "${MISSING_DEPS[@]}"; do
                if [ "$dep" == "nodejs" ]; then
                    curl -fsSL https://nodejs.org/dist/v20.15.1/node-v20.15.1-linux-x64.tar.xz -o /tmp/node.tar.xz >> "$LOG_FILE" 2>&1 &
                    pids+=($!)
                elif [ "$dep" == "docker" ]; then
                    apt-get install -y --download-only docker.io >> "$LOG_FILE" 2>&1 &
                    pids+=($!)
                elif [ "$dep" == "golang" ]; then
                    curl -fsSL https://go.dev/dl/go1.22.5.linux-amd64.tar.gz -o /tmp/go.tar.gz >> "$LOG_FILE" 2>&1 &
                    pids+=($!)
                fi
            done
            
            # Wait for all downloads to finish
            for pid in "${pids[@]}"; do
                wait $pid
            done
        ) &
        spinner $! "Downloading packages"

        # === PHASE 2: INSTALL (PARALLEL) ===
        echo -e "${BLUE}--> Phase 2: Installing dependencies (Parallel)...${NC}"
        
        (
            pids=()
            for dep in "${MISSING_DEPS[@]}"; do
                if [ "$dep" == "nodejs" ]; then
                    (
                        tar -xf /tmp/node.tar.xz -C /usr/local --strip-components=1
                        rm -f /tmp/node.tar.xz
                    ) >> "$LOG_FILE" 2>&1 &
                    pids+=($!)
                elif [ "$dep" == "golang" ]; then
                    (
                        rm -rf /usr/local/go
                        tar -xzf /tmp/go.tar.gz -C /usr/local
                        ln -sf /usr/local/go/bin/go /usr/local/bin/go
                        ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt
                        rm -f /tmp/go.tar.gz
                    ) >> "$LOG_FILE" 2>&1 &
                    pids+=($!)
                elif [ "$dep" == "docker" ]; then
                    (
                        DEBIAN_FRONTEND=noninteractive apt-get install -y docker.io
                        systemctl enable --now docker
                    ) >> "$LOG_FILE" 2>&1 &
                    pids+=($!)
                fi
            done
            
            # Wait for all installations to finish
            for pid in "${pids[@]}"; do
                wait $pid
            done
        ) &
        spinner $! "Installing packages"
        
        echo -e "${GREEN}--> Dependencies successfully installed!${NC}"
    else
        echo -e "${RED}Setup aborted. Please install missing dependencies manually.${NC}"
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

# Ensure data directory exists (fix "File exists" error)
if [ -f "/var/lib/hostanything/data" ]; then
    rm -f "/var/lib/hostanything/data"
fi
if [ -f "/var/lib/hostanything" ]; then
    rm -f "/var/lib/hostanything"
fi
mkdir -p /var/lib/hostanything/data
if [ -n "$SUDO_USER" ]; then
    chown -R $SUDO_USER:$SUDO_USER /var/lib/hostanything/data || true
fi

# 4. Build Phase (Parallel)
echo -e "\n${BLUE}--> Phase 3: Building Host Anything (Parallel)...${NC}"

(
    # Build Backend
    (
        if [ -n "$SUDO_USER" ]; then
            sudo -u $SUDO_USER -E make build
        else
            make build
        fi
    ) >> "$LOG_FILE" 2>&1 &
    PID_BACKEND=$!

    # Build Frontend
    (
        if [ -n "$SUDO_USER" ]; then
            sudo -u $SUDO_USER -E make build-web
        else
            make build-web
        fi
    ) >> "$LOG_FILE" 2>&1 &
    PID_FRONTEND=$!

    wait $PID_BACKEND
    wait $PID_FRONTEND
) &
spinner $! "Compiling Source Code"


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
