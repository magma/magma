#!/bin/bash

GO_VERSION_MAJOR=1
GO_VERSION_MINOR=13
LD_FLAGS="-s -w"
MY_FULL_PATH="$(cd "$(dirname "${0}")" && pwd)"
MAGMA_PATH=$(sed -E 's|fbcode/magma/.*|fbcode/magma|'  <<< "${MY_FULL_PATH}")
BIN_DIR="${HOME}/bin/magma/arm"
BIN_DEPLOY_DIR="/root"
IP="192.168.1.1"
USER="root"
PASSWORD="facebook"

# Note: "go version" output format is "xx xxxx go1.13.3 xxxx"
# shellcheck disable=SC2235
check_golang_version() {
  echo "Checking golang version"
  go_major="$(go version 2>/dev/null | cut -d' ' -f 3 | sed -E s/go// | cut -d'.' -f 1)"
  go_minor="$(go version 2>/dev/null | cut -d' ' -f 3 | sed -E s/go// | cut -d'.' -f 2)"
  if [ -z "${go_major}" ] || [ -z "${go_minor}" ] || \
    [ "${go_major}" -lt "${GO_VERSION_MAJOR}" ] || \
    ([ "${go_major}" -eq "${GO_VERSION_MAJOR}" ] && [ "${go_minor}" -lt "${GO_VERSION_MINOR}" ]); then
    echo "Need golang ${GO_VERSION_MAJOR}.${GO_VERSION_MINOR} or higher"
    echo "On mac, run \"brew install go\""
    exit 1
  fi
}

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
  echo "  --bin <path>     build path for the magma binaries"
  echo "                   default: ${BIN_DIR}"
  echo "  --build          build the magma binaries"
  echo "  --deploy         deploy magma binaries to the gateway"
  echo "                   default: build-only"
  echo "  --help           show this message"
  echo "  --ip             gateway ip address"
  echo "                   default: ${IP}"
  echo "  --passwd <pass>  gateway root password"
  echo "                   default: ${PASSWORD}"
  echo "  --show           show the magma binaries presently on the gateway"
}

bin_dir=${BIN_DIR}
ip=${IP}
magma_path=${MAGMA_PATH}
password=${PASSWORD}
while [ $# -gt 0 ]; do
  case $1 in
    --bin)
      shift
      bin_dir=$1
      ;;
    --build|-b)
      build=1
      ;;
    --deploy|-d)
      deploy=1
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    --ip)
      shift
      ip=$1
      ;;
    --passwd|-p)
      shift
      password=$1
      ;;
    --show|-s)
      show=1
      ;;
    *)
      echo "Unknown option \"$1\""
      usage
      exit 1
      ;;
  esac
  shift
done

# Sanity check options
if [ -z "${ip}" ] || [ -z "${password}" ] || [ -z "${bin_dir}" ]; then
  echo "Invalid option value"
  usage
  exit 1
fi

# Build
# shellcheck disable=SC2164
if [ -n "${build}" ]; then
  check_golang_version

  echo "Creating binary dir: \"${bin_dir}\""
  mkdir -p "${bin_dir}"

  echo "Building magmad"
  cd "${magma_path}/orc8r/gateway/go" || exit 1
  GOOS=linux GOARCH=arm go build -o "${bin_dir}" -ldflags="${LD_FLAGS}" magma/gateway/services/magmad

  echo "Building AAA service"
  cd "${magma_path}/feg/gateway" || exit 1
  GOOS=linux GOARCH=arm go build -tags link_local_service -o "${bin_dir}" -ldflags="${LD_FLAGS}" magma/feg/gateway/services/aaa/aaa_server

  echo "Building radius server"
  cd "${magma_path}/feg/radius/src" || exit 1
  GOOS=linux GOARCH=arm go build -o "${bin_dir}" -ldflags="${LD_FLAGS}" .
fi

# Check deploy tools
if [ -n "${deploy}" ] || [ -n "${show}" ]; then
  check_deploy_tools
fi

# Deploy
if [ -n "${deploy}" ]; then
  if [ ! -d "${bin_dir}" ]; then
    echo "Can't deploy from missing directory \"${bin_dir}\""
  else
    echo "Deploying magma binaries to gateway ${ip}"
    pscp -4 -batch -scp -pw "${password}" "${bin_dir}/magmad" "${bin_dir}/aaa_server" "${bin_dir}/radius" "${USER}@${ip}:${BIN_DEPLOY_DIR}/"
  fi
fi

# Show the deployed binaries on the gateway
# shellcheck disable=SC2005
if [ -n "${show}" ]; then
  echo "Showing magma binaries on gateway ${ip}"
  echo "$(plink -4 -batch -pw "${password}" "${USER}@${ip}" "ls -l ${BIN_DEPLOY_DIR}/aaa_server ${BIN_DEPLOY_DIR}/magmad ${BIN_DEPLOY_DIR}/radius")"
  echo "$(plink -4 -batch -pw "${password}" "${USER}@${ip}" "md5sum ${BIN_DEPLOY_DIR}/aaa_server ${BIN_DEPLOY_DIR}/magmad ${BIN_DEPLOY_DIR}/radius")"
fi

echo "Done"
