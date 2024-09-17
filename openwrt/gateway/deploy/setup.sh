#!/bin/bash

MY_FULL_PATH="$(cd "$(dirname "${0}")" && pwd)"
MAGMA_PATH=$(sed -E 's|/magma/.*|/magma|'  <<< "${MY_FULL_PATH}")
SRC_ROOT="${MAGMA_PATH}/openwrt/gateway/configs"
TGT_CONFIG_DIR="/etc/magma/configs"
IP="192.168.1.1"
USER="root"
PASSWORD="facebook"
SERVICES="magmad aaa_server"

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
  echo "  --boot           address of the bootstrapper"
  echo "  --cloud          address of the cloud controller"
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
    --boot|-b)
      shift
      boot=$1
      ;;
    --cloud|-c)
      shift
      cloud=$1
      ;;
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
if [ -z "${ip}" ] || [ -z "${password}" ] || [ -z "${cloud}" ] || [ -z "${boot}" ]; then
  echo "Invalid option value"
  usage
fi

# Dry run
if [ -n "${dryrun}" ]; then
  echo "Would set cloud controller address: ${cloud}"
  echo "Would set bootstrapper address: ${boot}"
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

# Set the cloud controller address
echo "Setting cloud controller address"
plink -4 -batch -pw "${password}" "${USER}@${ip}" \
  "for f in \$(find ${TGT_CONFIG_DIR} -name \*.yml -type f); do sed -i s/{{\.ControllerAddr}}/${cloud}/ \${f}; done"

# Set the bootstrapper address
echo "Setting bootstrapper address"
plink -4 -batch -pw "${password}" "${USER}@${ip}" \
  "for f in \$(find ${TGT_CONFIG_DIR} -name \*.yml -type f); do sed -i s/{{\.BootstrapperAddr}}/${boot}/ \${f}; done"

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
