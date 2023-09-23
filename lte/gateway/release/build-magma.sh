#!/bin/bash
#
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script builds Magma based on the current state of your repo. It needs to
# be run inside the VM.

set -ex
shopt -s extglob
SCRIPT_DIR="$(dirname "$(realpath "$0")")"

# Please update the version number accordingly for beta/stable builds
# Test builds are versioned automatically by fabfile.py
VERSION=1.8.0 # magma version number
SCTPD_MIN_VERSION=1.8.0 # earliest version of sctpd with which this version is compatible
DHCP_CLI_MIN_VERSION=1.8.0 # earliest version of dhcp_cli with which this version is compatible

# RelWithDebInfo or Debug
BUILD_TYPE=RelWithDebInfo

# Cmdline options that overwrite the version configs above
COMMIT_HASH=""  # hash of top magma commit
COMMIT_COUNT="" # count of commits on git main
CERT_FILE="$MAGMA_ROOT/.cache/test_certs/rootCA.pem"
CONTROL_PROXY_FILE="$MAGMA_ROOT/lte/gateway/configs/control_proxy.yml"
OS="ubuntu"

while [[ $# -gt 0 ]]
do
key="$1"
case $key in
    -v|--version)
    VERSION="$2"
    shift  # pass argument or value
    ;;
    -h|--hash)
    COMMIT_HASH="$2"
    shift
    ;;
    --commit-count)
    COMMIT_COUNT="$2"
    shift
    ;;
    -t|--type)
    BUILD_TYPE="$2"
    shift  # pass argument or value
    ;;
    -c|--cert)
    CERT_FILE="$2"
    shift
    ;;
    -p|--proxy)
    CONTROL_PROXY_FILE="$2"
    shift
    ;;
    --os)
    OS="$2"
    shift
    ;;
    *)
    echo "Error: unknown cmdline option:" $key
    echo "Usage: $0 [-v|--version V] [-i|--iteration I] [-h|--hash HASH]
    [-t|--type Debug|RelWithDebInfo] [-c|--cert <path to cert .pem file>]
    [-p|--proxy <path to control_proxy config .yml file]>
    [-u|--build-buntu build packages for ubuntu>"
    exit 1
    ;;
esac
shift  # past argument or value
done

case $BUILD_TYPE in
    Debug)
    ;;
    RelWithDebInfo)
    ;;
    *)
    echo "Error: unknown type option:" $BUILD_TYPE
    echo "Usage: [-t|--type Debug|RelWithDebInfo]"
    exit 1
    ;;
esac

case $OS in
    ubuntu)
    echo "Ubuntu package build"
    ;;
    *)
    echo "Error: unknown OS option:" $OS
    echo "Usage: [--os ubuntu]"
    exit 1
    ;;
esac


# Default options
BUILD_DATE=`date -u +"%Y%m%d%H%M%S"`
ARCH=amd64
PKGFMT=deb
MAGMA_PKGNAME=magma
SCTPD_PKGNAME=magma-sctpd
DHCP_CLI_PKGNAME=magma-dhcp-cli

# Magma system dependencies: anything that we depend on at the top level, add
# here.
MAGMA_DEPS=(
    "grpc-dev >= 1.15.0"
    "lighttpd >= 1.4.45"
    "libxslt1.1"
    "nghttp2-proxy >= 1.18.1"
    "python3-protobuf >= 3.20.3"
    "redis-server >= 3.2.0"
    "sudo"
    "dnsmasq >= 2.7"
    "net-tools" # for ifconfig
    "python3-pip"
    "python3-apt" # The version in pypi is abandoned and broken on stretch
    "python3-aioeventlet" # The version in pypi got deleted
    "libsystemd-dev"
    "libyaml-cpp-dev" # install yaml parser
    "libgoogle-glog-dev"
    "python-redis"
    "magma-cpp-redis"
    "libfolly-dev" # required for C++ services
    "libdouble-conversion-dev" # required for folly
    "libboost-chrono-dev" # required for folly
    "ntpdate" # required for eventd time synchronization
    "tshark" # required for call tracing
    "libtins-dev" # required for Connection tracker
    "libmnl-dev" # required for Connection tracker
    "getenvoy-envoy" # for envoy dep
    "uuid-dev" # for liagentd
    "libprotobuf17 >= 3.0.0"
    "nlohmann-json3-dev"
    "sentry-native"   # sessiond
    "td-agent-bit >= 1.7.8"
    # eBPF compile and load tools for kernsnoopd and AGW datapath
    # Ubuntu bcc lib (bpfcc-tools) is pretty old, use magma repo package
    "bcc-tools"
    "wireguard"
    "${DHCP_CLI_PKGNAME} >= ${DHCP_CLI_MIN_VERSION}"
    )

# OAI runtime dependencies
OAI_DEPS=(
    "libconfig9"
    "oai-asn1c"
    "oai-gnutls >= 3.1.23"
    "oai-nettle >= 1.0.1"
    "prometheus-cpp-dev >= 1.0.2"
    "liblfds710"
    "libsctp-dev"
    "magma-sctpd >= ${SCTPD_MIN_VERSION}"
    "libczmq-dev >= 4.0.2-7"
    "libasan5"
    "oai-freediameter >= 0.0.2"
    )

# OVS runtime dependencies
OVS_DEPS=(
      "magma-libfluid >= 0.1.0.7"
      "libopenvswitch >= 2.15.4-9-magma"
      "openvswitch-switch >= 2.15.4-9-magma"
      "openvswitch-common >= 2.15.4-9-magma"
      "openvswitch-datapath-dkms >= 2.15.4-9-magma"
      )

# generate string for FPM
SYSTEM_DEPS=""
for dep in "${MAGMA_DEPS[@]}"
do
    SYSTEM_DEPS=${SYSTEM_DEPS}" -d '"${dep}"'"
done
for dep in "${OAI_DEPS[@]}"
do
    SYSTEM_DEPS=${SYSTEM_DEPS}" -d '"${dep}"'"
done
for dep in "${OVS_DEPS[@]}"
do
    SYSTEM_DEPS=${SYSTEM_DEPS}" -d '"${dep}"'"
done

RELEASE_DIR=${MAGMA_ROOT}/lte/gateway/release
POSTINST=${RELEASE_DIR}/magma-postinst

# python environment
PY_VERSION=python3.8
PY_PKG_LOC=dist-packages
PY_DEST=/usr/local/lib/${PY_VERSION}/${PY_PKG_LOC}
PY_PROTOS=${PYTHON_BUILD}/gen/
PY_LTE=${MAGMA_ROOT}/lte/gateway/python
PY_ORC8R=${MAGMA_ROOT}/orc8r/gateway/python
PY_TMP_BUILD=/tmp/build-${MAGMA_PKGNAME}
PY_TMP_BUILD_SUFFIX=/usr/lib/python3/${PY_PKG_LOC}

PWD=`pwd`

glob_files () {
    # Given a list of files represented by the pattern in $1, and a package
    # output location in $2, generate a string of file locations that can be
    # passed to FPM. If $1 is a glob, you MUST surround it with quotes!
    #
    # For example, if you wanted to have all the files maching
    # foo/bar/*.yml end up in /etc/magma/, you would call:
    #
    # glob_files "foo/bar/*.yml" /etc/magma
    #
    # which would return:
    #
    # foo/bar/baz.yml=/etc/magma foo/bar/qux.yml=/etc/magma
    #
    # This is useful because fpm only accepts individual files or entire
    # directories for the dir package source type.
    RES=""
    for f in $1
    do
        RES="$RES $f=$2"
    done

    echo $RES
}

# The resulting package is placed in $OUTPUT_DIR
# or in the cwd.
if [ -z "$1" ]; then
  OUTPUT_DIR=${PWD}
else
  OUTPUT_DIR=$1
  if [ ! -d "$OUTPUT_DIR" ]; then
    echo "error: $OUTPUT_DIR is not a valid directory. Exiting..."
    exit 1
  fi
fi

# Build OAI and sessiond C/C++ services
cd "${MAGMA_ROOT}/lte/gateway"
OAI_BUILD="${C_BUILD}/core/oai"
SESSIOND_BUILD="${C_BUILD}/session_manager"
CONNECTIOND_BUILD="${C_BUILD}/connection_tracker"
LI_AGENT_BUILD="${C_BUILD}/li_agent"
SCTPD_BUILD="${C_BUILD}/sctpd/src"

make build_oai BUILD_TYPE="${BUILD_TYPE}"
make build_session_manager BUILD_TYPE="${BUILD_TYPE}"
make build_sctpd BUILD_TYPE="${BUILD_TYPE}"
make build_connection_tracker BUILD_TYPE="${BUILD_TYPE}"
make build_li_agent BUILD_TYPE="${BUILD_TYPE}"

# Build Magma Envoy Controller service
cd "${MAGMA_ROOT}/feg/gateway"
make install_envoy_controller

# Next, gather up the python files and put them into a build path.
#
# Note: Debian-based distributions install packages by default into a
# dist-packages directory, which is different than other distros, which drop
# packages into the site-packages directory.

# clean python build dir
if [ -d ${PY_TMP_BUILD} ]; then
    rm -r ${PY_TMP_BUILD}
fi

FULL_VERSION=${VERSION}-$(date +%s)-${COMMIT_HASH}
COMMIT_HASH_WITH_VERSION="magma@${VERSION}.${COMMIT_COUNT}-${COMMIT_HASH}"

# first do python protos and then build the python packages.
# library will be dropped in $PY_TMP_BUILD/usr/lib/python3/dist-packages
# scripts will be dropped in $PY_TMP_BUILD/usr/bin.
# Use pydep to generate the lockfile and python deps
# update magma.lockfile if needed (see Makefile)
# adjust mtime of a setup.py to force update
# (e.g. `touch ${PY_LTE}/setup.py`)
pushd "${RELEASE_DIR}" || exit 1
make os_release=$OS -e magma.lockfile
popd

cd ${PY_ORC8R}
make protos
PKG_VERSION=${FULL_VERSION} ${PY_VERSION} setup.py install --root ${PY_TMP_BUILD} --install-layout deb \
    --no-compile --single-version-externally-managed

ORC8R_PY_DEPS=`${RELEASE_DIR}/pydep lockfile ${RELEASE_DIR}/magma.lockfile.$OS`

cd ${PY_LTE}
make protos
make swagger
PKG_VERSION=${FULL_VERSION} ${PY_VERSION} setup.py install --root ${PY_TMP_BUILD} --install-layout deb \
    --no-compile --single-version-externally-managed
${RELEASE_DIR}/pydep finddep -l ${RELEASE_DIR}/magma.lockfile.$OS setup.py
LTE_PY_DEPS=`${RELEASE_DIR}/pydep lockfile ${RELEASE_DIR}/magma.lockfile.$OS`

MAGMA_BUILD_PATH=${OUTPUT_DIR}/${MAGMA_PKGNAME}_${FULL_VERSION}_${ARCH}.${PKGFMT}
SCTPD_BUILD_PATH=${OUTPUT_DIR}/${SCTPD_PKGNAME}_${FULL_VERSION}_${ARCH}.${PKGFMT}
DHCP_CLI_BUILD_PATH=${OUTPUT_DIR}/${DHCP_CLI_PKGNAME}_${FULL_VERSION}_${ARCH}.${PKGFMT}

cd $PWD
# remove old packages
if [ -f "${MAGMA_BUILD_PATH}" ]; then
  rm "${MAGMA_BUILD_PATH}"
fi
if [ -f "${SCTPD_BUILD_PATH}" ]; then
  rm "${SCTPD_BUILD_PATH}"
fi
if [ -f "${DHCP_CLI_BUILD_PATH}" ]; then
  rm "${DHCP_CLI_BUILD_PATH}"
fi

SERVICE_DIR="/etc/systemd/system/"
ANSIBLE_FILES="${MAGMA_ROOT}/lte/gateway/deploy/roles/magma/files"

SCTPD_VERSION_FILE=$(mktemp)
SCTPD_MIN_VERSION_FILE=$(mktemp)
COMMIT_HASH_FILE=$(mktemp)

# files to be removed should be safely named (no special chars from mktemp)
# use current value (see https://github.com/koalaman/shellcheck/wiki/SC2064)
# shellcheck disable=SC2064
trap "rm -f '${SCTPD_VERSION_FILE}' '${SCTPD_MIN_VERSION_FILE}' '${COMMIT_HASH_FILE}'" EXIT

echo "${FULL_VERSION}" > "${SCTPD_VERSION_FILE}"
echo "${SCTPD_MIN_VERSION}" > "${SCTPD_MIN_VERSION_FILE}"
echo "COMMIT_HASH=\"${COMMIT_HASH_WITH_VERSION}\"" > "${COMMIT_HASH_FILE}"

# meta info
DESCRIPTION_SCTPD="Magma SCTPD"
DESCRIPTION_AGW="Magma Access Gateway"
if [ "${BUILD_TYPE}" != "RelWithDebInfo" ]; then
    DESCRIPTION_SCTPD="${DESCRIPTION_SCTPD} - dev build"
    DESCRIPTION_AGW="${DESCRIPTION_AGW} - dev build"
fi
URL="https://github.com/magma/magma/"
VENDOR="magma"
LICENSE="BSD-3-Clause"
MAINTAINER="The Magma Authors <main@lists.magmacore.org>"

BUILDCMD="fpm \
-s dir \
-t ${PKGFMT} \
-a ${ARCH} \
-n ${SCTPD_PKGNAME} \
-v ${FULL_VERSION} \
--provides ${SCTPD_PKGNAME} \
--replaces ${SCTPD_PKGNAME} \
--package ${SCTPD_BUILD_PATH} \
--description '${DESCRIPTION_SCTPD}' \
--url '${URL}' \
--vendor '${VENDOR}' \
--license '${LICENSE}' \
--maintainer '${MAINTAINER}' \
--exclude '*/.ignoreme' \
${SCTPD_BUILD}/sctpd=/usr/local/sbin/ \
${SCTPD_VERSION_FILE}=/usr/local/share/sctpd/version \
$(glob_files "${SERVICE_DIR}/sctpd.service" /etc/systemd/system/sctpd.service) \
${MAGMA_ROOT}/LICENSE=/usr/share/doc/${SCTPD_PKGNAME}/"

eval "$BUILDCMD"

BUILDCMD="fpm \
-s dir \
-t ${PKGFMT} \
-a ${ARCH} \
-n ${MAGMA_PKGNAME} \
-v ${FULL_VERSION} \
--provides ${MAGMA_PKGNAME} \
--replaces ${MAGMA_PKGNAME} \
--package ${MAGMA_BUILD_PATH} \
--description '${DESCRIPTION_AGW}' \
--url '${URL}' \
--vendor '${VENDOR}' \
--license '${LICENSE}' \
--maintainer '${MAINTAINER}' \
--after-install ${POSTINST} \
--exclude '*/.ignoreme' \
--config-files /etc/sysctl.d/99-magma.conf \
${ORC8R_PY_DEPS} \
${LTE_PY_DEPS} \
${SYSTEM_DEPS} \
${OAI_BUILD}/oai_mme/mme=/usr/local/bin/ \
${SESSIOND_BUILD}/sessiond=/usr/local/bin/ \
${CONNECTIOND_BUILD}/src/connectiond=/usr/local/bin/ \
${LI_AGENT_BUILD}/src/liagentd=/usr/local/bin/ \
${GO_BUILD}/envoy_controller=/usr/local/bin/ \
${SCTPD_MIN_VERSION_FILE}=/usr/local/share/magma/sctpd_min_version \
${COMMIT_HASH_FILE}=/usr/local/share/magma/commit_hash \
$(glob_files "${SERVICE_DIR}/magma@.service" /etc/systemd/system/magma@.service) \
$(glob_files "${SERVICE_DIR}/magma@control_proxy.service" /etc/systemd/system/magma@control_proxy.service) \
$(glob_files "${SERVICE_DIR}/magma@magmad.service" /etc/systemd/system/magma@magmad.service) \
$(glob_files "${SERVICE_DIR}/magma@mme.service" /etc/systemd/system/magma@mme.service) \
$(glob_files "${SERVICE_DIR}/magma@sessiond.service" /etc/systemd/system/magma@sessiond.service) \
$(glob_files "${SERVICE_DIR}/magma@connectiond.service" /etc/systemd/system/magma@connectiond.service) \
$(glob_files "${SERVICE_DIR}/magma@mobilityd.service" /etc/systemd/system/magma@mobilityd.service) \
$(glob_files "${SERVICE_DIR}/magma@pipelined.service" /etc/systemd/system/magma@pipelined.service) \
$(glob_files "${SERVICE_DIR}/magma_dp@envoy.service" /etc/systemd/system/magma_dp@envoy.service) \
$(glob_files "${SERVICE_DIR}/magma@envoy_controller.service" /etc/systemd/system/magma@envoy_controller.service) \
$(glob_files "${SERVICE_DIR}/magma@redirectd.service" /etc/systemd/system/magma@redirectd.service) \
$(glob_files "${SERVICE_DIR}/magma@dnsd.service" /etc/systemd/system/magma@dnsd.service) \
$(glob_files "${SERVICE_DIR}/magma@lighttpd.service" /etc/systemd/system/magma@lighttpd.service) \
$(glob_files "${SERVICE_DIR}/magma@redis.service" /etc/systemd/system/magma@redis.service) \
$(glob_files "${SERVICE_DIR}/magma@td-agent-bit.service" /etc/systemd/system/magma@td-agent-bit.service) \
$(glob_files "${MAGMA_ROOT}/lte/gateway/configs/!(control_proxy.yml|pipelined.yml|sessiond.yml|connectiond.yml)" /etc/magma/) \
$(glob_files "${MAGMA_ROOT}/lte/gateway/configs/pipelined.yml_prod" /etc/magma/pipelined.yml) \
$(glob_files "${MAGMA_ROOT}/lte/gateway/configs/sessiond.yml_prod" /etc/magma/sessiond.yml) \
$(glob_files "${MAGMA_ROOT}/lte/gateway/configs/templates/*" /etc/magma/templates/) \
$(glob_files "${MAGMA_ROOT}/orc8r/gateway/configs/templates/*" /etc/magma/templates/) \
$(glob_files "${MAGMA_ROOT}/lte/gateway/python/magma/kernsnoopd/ebpf/*" /var/opt/magma/ebpf/kernsnoopd/) \
$(glob_files "${MAGMA_ROOT}/lte/gateway/python/magma/pipelined/ebpf/*" /var/opt/magma/ebpf/) \
${CONTROL_PROXY_FILE}=/etc/magma/ \
$(glob_files "${ANSIBLE_FILES}/magma_modules_load" /etc/modules-load.d/magma.conf) \
$(glob_files "${ANSIBLE_FILES}/configure_envoy_namespace.sh" /usr/local/bin/ ) \
$(glob_files "${ANSIBLE_FILES}/envoy.yaml" /var/opt/magma/ ) \
$(glob_files "${ANSIBLE_FILES}/logrotate_oai.conf" /etc/logrotate.d/oai) \
$(glob_files "${ANSIBLE_FILES}/logrotate_rsyslog.conf" /etc/logrotate.d/rsyslog.magma) \
$(glob_files "${ANSIBLE_FILES}/local-cdn/*" /var/www/local-cdn/) \
${ANSIBLE_FILES}/99-magma.conf=/etc/sysctl.d/ \
${ANSIBLE_FILES}/magma_ifaces_gtp=/etc/network/interfaces.d/gtp \
${ANSIBLE_FILES}/20auto-upgrades=/etc/apt/apt.conf.d/20auto-upgrades \
${ANSIBLE_FILES}/coredump=/usr/local/bin/ \
${MAGMA_ROOT}/orc8r/tools/ansible/roles/fluent_bit/files/60-fluent-bit.conf=/etc/rsyslog.d/60-fluent-bit.conf \
${ANSIBLE_FILES}/set_irq_affinity=/usr/local/bin/ \
${ANSIBLE_FILES}/ovs-kmod-upgrade.sh=/usr/local/bin/ \
${ANSIBLE_FILES}/magma-bridge-reset.sh=/usr/local/bin/ \
${ANSIBLE_FILES}/magma-setup-wg.sh=/usr/local/bin/ \
${ANSIBLE_FILES}/magma-create-gtp-port.sh=/usr/local/bin/ \
${PY_PROTOS}=${PY_DEST} \
$(glob_files "${PY_TMP_BUILD}/${PY_TMP_BUILD_SUFFIX}/${MAGMA_PKGNAME}*" ${PY_DEST}) \
$(glob_files "${PY_TMP_BUILD}/${PY_TMP_BUILD_SUFFIX}/*.egg-info" ${PY_DEST}) \
$(glob_files "${PY_TMP_BUILD}/usr/bin/*" /usr/local/bin/) \
${MAGMA_ROOT}/LICENSE=/usr/share/doc/${MAGMA_PKGNAME}/ \
" # Leave this quote on a new line to mark end of BUILDCMD

eval "$BUILDCMD"

DESCRIPTION_DHCP="Magma DHCP helper CLI"
LICENSE_DHCP="GPL-2.0"
SCAPY_PACKAGE="python3-scapy"
SCAPY_VERSION="2.4.3"

BUILDCMD="fpm \
-s dir \
-t ${PKGFMT} \
-a ${ARCH} \
-n ${DHCP_CLI_PKGNAME} \
-v ${FULL_VERSION} \
--provides ${DHCP_CLI_PKGNAME} \
--replaces ${DHCP_CLI_PKGNAME} \
--package ${DHCP_CLI_BUILD_PATH} \
--description '${DESCRIPTION_DHCP}' \
--url '${URL}' \
--vendor '${VENDOR}' \
--license '${LICENSE_DHCP}' \
--maintainer '${MAINTAINER}' \
--depends '${SCAPY_PACKAGE} >= ${SCAPY_VERSION}' \
${MAGMA_ROOT}/lte/gateway/python/dhcp_helper_cli/dhcp_helper_cli.py=/usr/local/bin/ \
${MAGMA_ROOT}/lte/gateway/python/dhcp_helper_cli/LICENSE=/usr/share/doc/${DHCP_CLI_PKGNAME}/ \
"

eval "$BUILDCMD"
