#!/usr/bin/env bash
# ==============================================================================
# KenPanel v2.3 — Clean Uninstaller & Environment Purge Script
# Specification Section 27.3: Clean uninstallation with confirmation and backup
# ==============================================================================
set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}================================================================${NC}"
echo -e "${CYAN}               KenPanel Clean Uninstaller v2.3                  ${NC}"
echo -e "${CYAN}================================================================${NC}"

if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}[ERROR] This uninstaller must be run as root.${NC}" 1>&2
   exit 1
fi

PURGE_DATA=false
FORCE=false

while [[ $# -gt 0 ]]; do
  case $1 in
    --purge-data)
      PURGE_DATA=true
      shift
      ;;
    -f|--force)
      FORCE=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

if [[ "$FORCE" != true ]]; then
  echo -e "${YELLOW}[WARNING] This will stop and remove all KenPanel control plane and node agent services.${NC}"
  if [[ "$PURGE_DATA" == true ]]; then
    echo -e "${RED}[DANGER] --purge-data specified! All databases, site definitions and configuration history will be permanently erased.${NC}"
  else
    echo -e "${GREEN}[INFO] User data in /var/lib/kenpanel will be PRESERVED.${NC}"
  fi
  read -rp "Are you sure you wish to continue? (y/N): " confirm
  if [[ "$confirm" != [yY] && "$confirm" != [yY][eE][sS] ]]; then
    echo "Aborted by user."
    exit 0
  fi
fi

echo -e "${CYAN}[1/4] Stopping and disabling systemd services...${NC}"
systemctl stop kenpanel.service 2>/dev/null || true
systemctl stop kenpanel-agent.service 2>/dev/null || true
systemctl disable kenpanel.service 2>/dev/null || true
systemctl disable kenpanel-agent.service 2>/dev/null || true

echo -e "${CYAN}[2/4] Removing systemd unit files...${NC}"
rm -f /etc/systemd/system/kenpanel.service
rm -f /etc/systemd/system/kenpanel-agent.service
systemctl daemon-reload

echo -e "${CYAN}[3/4] Removing binaries and CLI symlinks...${NC}"
rm -f /usr/local/bin/kenpanel
rm -f /usr/local/bin/kenpanel-agent
rm -f /usr/local/bin/kenpanel-cli
rm -f /usr/local/bin/kenpanel-transplant
rm -f /usr/local/bin/kenpanel-capsule
rm -f /usr/local/bin/kenpanel-mailbridge
rm -f /usr/local/bin/kenpanel-recovery

echo -e "${CYAN}[4/4] Processing configuration and data directories...${NC}"
rm -rf /etc/kenpanel

if [[ "$PURGE_DATA" == true ]]; then
  echo -e "${RED}Purging /var/lib/kenpanel and log directories...${NC}"
  rm -rf /var/lib/kenpanel
  rm -rf /var/log/kenpanel
else
  echo -e "${GREEN}Preserved /var/lib/kenpanel for future reinstallation or disaster extraction.${NC}"
fi

echo -e "${GREEN}================================================================${NC}"
echo -e "${GREEN}       KenPanel has been successfully uninstalled.              ${NC}"
echo -e "${GREEN}================================================================${NC}"
