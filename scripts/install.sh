#!/bin/sh
# KenPanel Unified Control Plane - Official Bootstrap Installer
# Website: https://kenpanel.kentralo.com
# Usage: curl -fsSL https://kenpanel.kentralo.com/install.sh | sudo bash [OPTIONS]

set -e

KENPANEL_VERSION="2.3.0"
KENPANEL_RELEASE_URL="https://kenpanel.kentralo.com/releases"
INSTALL_DIR="/opt/kenpanel"
CONFIG_DIR="/etc/kenpanel"
DATA_DIR="/var/lib/kenpanel"
LOG_DIR="/var/log/kenpanel"
RUN_DIR="/run/kenpanel"

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}[SUCCESS]${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1"
}

error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
    exit 1
}

# 1. Root Check
if [ "$(id -u)" -ne 0 ]; then
    error "KenPanel installer must be executed as root (use 'sudo bash')."
fi

echo "============================================================"
echo "          KenPanel Unified Control Plane Installer          "
echo "                 Version: ${KENPANEL_VERSION}               "
echo "        One panel. Your servers. Your rules.                "
echo "============================================================"
echo ""

# 2. Architecture & OS Detection
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)
        KP_ARCH="x86_64"
        ;;
    aarch64|arm64)
        KP_ARCH="arm64"
        ;;
    *)
        error "Unsupported architecture: $ARCH. KenPanel requires x86_64 or arm64."
        ;;
esac

OS=""
OS_VER=""
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS="$ID"
    OS_VER="$VERSION_ID"
else
    error "Could not detect operating system via /etc/os-release."
fi

info "Detected System: ${OS} ${OS_VER} (${KP_ARCH})"

# Compatibility validation
SUPPORTED=false
case "$OS" in
    ubuntu)
        if [ "$OS_VER" = "22.04" ] || [ "$OS_VER" = "24.04" ]; then SUPPORTED=true; fi
        ;;
    debian)
        if [ "$OS_VER" = "11" ] || [ "$OS_VER" = "12" ]; then SUPPORTED=true; fi
        ;;
    almalinux|rocky|rhel)
        case "$OS_VER" in
            9*|8*) SUPPORTED=true ;;
        esac
        ;;
    alpine)
        SUPPORTED=true
        ;;
esac

if [ "$SUPPORTED" = false ]; then
    warn "Distribution ${OS} ${OS_VER} is not in the certified Tier-1 matrix. Proceeding in generic adapter mode."
fi

# 3. Preflight Resource Verification
info "Running preflight resource checks..."
MEM_TOTAL_KB=$(grep MemTotal /proc/meminfo | awk '{print $2}')
MEM_TOTAL_MB=$((MEM_TOTAL_KB / 1024))

if [ "$MEM_TOTAL_MB" -lt 1024 ]; then
    warn "System has only ${MEM_TOTAL_MB}MB RAM. KenPanel recommends at least 2048MB for production environments."
else
    success "Memory check passed: ${MEM_TOTAL_MB}MB available."
fi

DISK_AVAIL_KB=$(df -k / | awk 'NR==2 {print $4}')
DISK_AVAIL_MB=$((DISK_AVAIL_KB / 1024))
if [ "$DISK_AVAIL_MB" -lt 4096 ]; then
    error "Root partition has only ${DISK_AVAIL_MB}MB free space. KenPanel requires at least 4096MB."
else
    success "Disk space check passed: ${DISK_AVAIL_MB}MB free."
fi

# 4. Create Dedicated Users and Directories
info "Configuring system user and directory structure..."
if ! id -u kenpanel >/dev/null 2>&1; then
    useradd -r -s /bin/false -d "$DATA_DIR" kenpanel
fi

mkdir -p "$INSTALL_DIR" "$CONFIG_DIR" "$DATA_DIR" "$LOG_DIR" "$RUN_DIR"
chown -R kenpanel:kenpanel "$DATA_DIR" "$LOG_DIR" "$RUN_DIR"
chmod 750 "$DATA_DIR" "$LOG_DIR"
chmod 700 "$CONFIG_DIR"

# 5. Initialize Configuration if Not Present
if [ ! -f "$CONFIG_DIR/kenpanel.conf" ]; then
    info "Generating initial security configuration..."
    JWT_SECRET=$(head -c 32 /dev/urandom | base64)
    cat <<EOF > "$CONFIG_DIR/kenpanel.conf"
# KenPanel Core Configuration - Generated at $(date -u +"%Y-%m-%dT%H:%M:%SZ")
listen_addr = "127.0.0.1:8443"
data_dir = "$DATA_DIR"
log_dir = "$LOG_DIR"
jwt_secret = "$JWT_SECRET"
allow_registration = false
telemetry_enabled = false
EOF
    chmod 600 "$CONFIG_DIR/kenpanel.conf"
    success "Configuration generated at $CONFIG_DIR/kenpanel.conf"
fi

# 6. Service Management (Systemd / OpenRC)
if [ -d /run/systemd/system ]; then
    info "Installing systemd unit files..."
    cat <<EOF > /etc/systemd/system/kenpanel.service
[Unit]
Description=KenPanel Unified Control Plane
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
User=kenpanel
Group=kenpanel
WorkingDirectory=$DATA_DIR
ExecStart=$INSTALL_DIR/bin/kenpanel server --config $CONFIG_DIR/kenpanel.conf
Restart=on-failure
RestartSec=5s
LimitNOFILE=65536
CapabilityBoundingSet=CAP_NET_BIND_SERVICE
AmbientCapabilities=CAP_NET_BIND_SERVICE

[Install]
WantedBy=multi-user.target
EOF

    cat <<EOF > /etc/systemd/system/kenpanel-agent.service
[Unit]
Description=KenPanel Privileged Node Daemon
After=network.target

[Service]
Type=simple
User=root
ExecStart=$INSTALL_DIR/bin/kenpanel-agent daemon --socket $RUN_DIR/agent.sock
Restart=on-failure
RestartSec=5s
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload || true
    info "Systemd units registered: kenpanel.service, kenpanel-agent.service."
fi

echo ""
echo "============================================================"
success "KenPanel preparation complete!"
echo "To initialize the first administrator and start services, run:"
echo "  kenpanel-cli admin init"
echo "  systemctl start kenpanel kenpanel-agent"
echo "Access the control panel at: https://$(hostname -I 2>/dev/null | awk '{print $1}' || echo '127.0.0.1'):8443"
echo "============================================================"
