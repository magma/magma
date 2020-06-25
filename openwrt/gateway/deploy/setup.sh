#!/bin/bash

MY_FULL_PATH="$(cd "$(dirname "${0}")" && pwd)"
MAGMA_PATH=$(sed -E 's|fbcode/magma/.*|fbcode/magma|'  <<< "${MY_FULL_PATH}")
SRC_ROOT="${MAGMA_PATH}/openwrt/gateway/configs"
IP="192.168.1.1"
USER="root"
PASSWORD="facebook"
SERVICES="magmad aaa_server radius"

check_deploy_tools() {
  echo "Checking PuTTY Secure Copy client tools"
  if ! command -v pscp > /dev/null || ! command -v plink > /dev/null; then
    echo "Need PuTTY Secure Copy client tools"
    echo "On mac, run \"brew install putty\""
    exit 1
  fi
}

usage() {
  echo "Usage: $0 [options]"
  echo " Initialize OpenWrt gateway with configuration and init scripts"
  echo "  --dryrun         show the gateway updates"
  echo "  --help           show this message"
  echo "  --ip             gateway ip address"
  echo "                   default: ${IP}"
  echo "  --noauto         disable magma services autostart"
  echo "  --passwd <pass>  gateway root password"
  echo "                   default: ${PASSWORD}"
  echo "  --start          start magma services as the last step"
  exit 1
}

ip=${IP}
password=${PASSWORD}
autostart=enable
while [ $# -gt 0 ]; do
  case $1 in
    --dryrun|-n)
      dryrun=1
      ;;
    --help|-h)
      usage
      ;;
    --ip)
      shift
      ip=$1
      ;;
    --noauto)
      autostart=disable
      ;;
    --passwd|-p)
      shift
      password=$1
      ;;
    --start|-s)
      start=1
      ;;
    *)
      echo "Unknown option \"$1\""
      usage
      ;;
  esac
  shift
done

# Sanity check the options
if [ -z "${ip}" ] || [ -z "${password}" ]; then
  echo "Invalid option value"
  usage
fi

# Dry run
if [ -n "${dryrun}" ]; then
  echo "Would update the following gateway files:"
  find "${SRC_ROOT}" -type f | sed "s|${SRC_ROOT}||"
  echo "autostart for magma services would be ${autostart}d"
  echo "magma services would $(if [ -z "${start}" ]; then echo "not"; fi) be started"
  exit 0
fi

# Check the deploy tools
echo "Checking deploy tools"
check_deploy_tools

# Stop magma services on the gateway
for s in ${SERVICES}; do
  echo "stop ${s}"
  plink -4 -batch -pw "${password}" "${USER}@${ip}" "/etc/init.d/${s} stop 2>/dev/null >/dev/null;killall ${s} 2>/dev/null >/dev/null"
done

# Update the gateway
echo "Copying files to gateway"
pscp -4 -batch -scp -p -r -pw "${password}" "${SRC_ROOT}/" "${USER}@${ip}:/"

# Enable or disable auto start
for s in ${SERVICES}; do
  echo "${autostart} ${s}"
  plink -4 -batch -pw "${password}" "${USER}@${ip}" "/etc/init.d/${s} ${autostart}"
done

# Start the magma services
if [ -n "${start}" ]; then
  for s in ${SERVICES}; do
    echo "start ${s}"
    plink -4 -batch -pw "${password}" "${USER}@${ip}" "/etc/init.d/${s} start"
  done
fi

echo "Done"
