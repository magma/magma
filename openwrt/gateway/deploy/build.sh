#!/bin/bash

GO_VERSION_MAJOR=1
GO_VERSION_MINOR=13
LD_FLAGS="-s -w" # assume stripped non-debug binary
MY_FULL_PATH="$(cd "$(dirname "${0}")" && pwd)"
MAGMA_PATH=$(sed -E 's|/magma/.*|/magma|'  <<< "${MY_FULL_PATH}")
BIN_DIR="${HOME}/bin/magma/arm"
BIN_DEPLOY_DIR="/sbin"
IP="192.168.1.1"
USER="root"
PASSWORD="facebook"
BINARIES="magmad aaa_server"

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

check_compression_tool() {
  echo "Checking Ultimate Packer for eXecutables (upx)"
  if ! command -v upx > /dev/null; then
    echo "Need Ultimate Packer for eXecutables (upx)"
    echo "On mac, run \"brew install upx\""
    exit 1
  fi
}

set_compression_flags() {
  if [ "${compression_flags}" = "fast" ]; then
    compression_flags="-1"
  elif [ "${compression_flags}" = "slow" ]; then
    echo "Compression may take 12+ minutes..."
    compression_flags="--brute"
  else
    echo "Invalid compression option \"${compression_flags}\""
    usage
  fi
}

# Compress one binary
do_compress() {
  path="${bin_dir}/${1}"
  comp_path="${path}.compressed"
  echo "Compressing \"${path}\""
  rm -f "${comp_path}"
  upx "${compression_flags}" -o "${comp_path}" "${path}"
  mv "${comp_path}" "${path}"
}

usage() {
  echo "Usage: $0 [options]"
  echo "  --bin <path>       build path for the magma binaries"
  echo "                     default: ${BIN_DIR}"
  echo "  --build            build the magma binaries"
  echo "  --compress <speed> compress magma binaries"
  echo "                     speed: \"fast\" (seconds) or \"slow\" (minutes)"
  echo "                     requires --build"
  echo "  --debug            generate debug binaries"
  echo "                     requires --build, does not work with --compress"
  echo "  --deploy           deploy magma binaries to the gateway"
  echo "  --help             show this message"
  echo "  --ip               gateway ip address"
  echo "                     default: ${IP}"
  echo "  --passwd <pass>    gateway root password"
  echo "                     default: ${PASSWORD}"
  echo "  --show             show the magma binaries presently on the gateway"
  exit 1
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
    --compress|-c)
      compress=1
      shift
      compression_flags=$1
      ;;
    --debug)
      debug=1
      ;;
    --deploy|-d)
      deploy=1
      ;;
    --help|-h)
      usage
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
      ;;
  esac
  shift
done


# Sanity check the options
if [ -n "${compress}" ]; then
  if [ -z "${build}" ]; then
    echo "--compress requires --build"
    usage
  fi
  if [ -n "${debug}" ]; then
    echo "--debug does not work with --compress"
    usage
  fi
  set_compression_flags
fi

if [ -n "${debug}" ]; then
  if [ -z "${build}" ]; then
    echo "--debug requires --build"
    usage
  fi
  LD_FLAGS=""
fi

if [ -z "${build}" ] && [ -z "${deploy}" ] && [ -z "${show}" ]; then
  echo "Nothing to do."
  usage
fi

if [ -z "${ip}" ] || [ -z "${password}" ] || [ -z "${bin_dir}" ]; then
  echo "Invalid option value"
  usage
fi

# Build
# shellcheck disable=SC2164
if [ -n "${build}" ]; then
  check_golang_version

  echo "Creating binary dir: \"${bin_dir}\""
  mkdir -p "${bin_dir}"

  echo "Deleting existing binaries"
  for b in ${BINARIES}; do
    rm -f "${bin_dir}/${b}"
  done

  echo "Building magmad"
  cd "${magma_path}/orc8r/gateway/go" || exit 1
  GOOS=linux GOARCH=arm go build -o "${bin_dir}" -ldflags="${LD_FLAGS}" magma/gateway/services/magmad

  echo "Building AAA service"
  cd "${magma_path}/feg/gateway" || exit 1
  GOOS=linux GOARCH=arm go build -tags link_local_service,with_builtin_radius -o "${bin_dir}" -ldflags="${LD_FLAGS}" magma/feg/gateway/services/aaa/aaa_server
fi

# Compress binaries
if [ -n "${compress}" ]; then
  check_compression_tool
  for b in ${BINARIES}; do
    do_compress "${b}"
  done
fi

# Check deploy tools
if [ -n "${deploy}" ] || [ -n "${show}" ]; then
  check_deploy_tools
fi

# Deploy
if [ -n "${deploy}" ]; then
  for b in ${BINARIES}; do
    f="${bin_dir}/${b}"
    if [ ! -f "${f}" ]; then
      echo "Can't deploy. Missing binary: \"${f}\""
      return 1
    fi
    echo "stop ${b}"
    plink -4 -batch -pw "${password}" "${USER}@${ip}" "/etc/init.d/${b} stop 2>/dev/null >/dev/null;killall ${b} 2>/dev/null >/dev/null"
  done
  echo "Deploying magma binaries to gateway ${ip}"
  pscp -4 -batch -scp -pw "${password}" "${bin_dir}/magmad" "${bin_dir}/aaa_server" "${USER}@${ip}:${BIN_DEPLOY_DIR}/"
fi

# Show the deployed binaries on the gateway
# shellcheck disable=SC2005
if [ -n "${show}" ]; then
  echo "Showing magma binaries on gateway ${ip}"
  echo "$(plink -4 -batch -pw "${password}" "${USER}@${ip}" "ls -l ${BIN_DEPLOY_DIR}/aaa_server ${BIN_DEPLOY_DIR}/magmad")"
  echo "$(plink -4 -batch -pw "${password}" "${USER}@${ip}" "md5sum ${BIN_DEPLOY_DIR}/aaa_server ${BIN_DEPLOY_DIR}/magmad")"
fi

echo "Done"
