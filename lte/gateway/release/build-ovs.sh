#!/bin/bash
#
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Generate the debian package from source for OVS
# Usage:
#   Run this script from a patched OVS source directory
#
# Example output:
#   libopenvswitch_2.8.0-1_amd64.deb
#   libopenvswitch-dev_2.8.0-1_amd64.deb
#   oai-gtp_4.9-6_amd64.deb
#   openvswitch-common_2.8.0-1_amd64.deb
#   openvswitch-datapath-dkms_2.8.0-1_all.deb
#   openvswitch-datapath-module-4.9.0-0.bpo.1-amd64_2.8.0-1_amd64.deb
#   openvswitch-datapath-source_2.8.0-1_all.deb
#   openvswitch-dbg_2.8.0-1_amd64.deb
#   openvswitch-pki_2.8.0-1_all.deb
#   openvswitch-switch_2.8.0-1_amd64.deb
#   openvswitch-test_2.8.0-1_all.deb
#   openvswitch-testcontroller_2.8.0-1_amd64.deb
#   openvswitch-vtep_2.8.0-1_amd64.deb
#   ovn-central_2.8.0-1_amd64.deb
#   ovn-common_2.8.0-1_amd64.deb
#   ovn-controller-vtep_2.8.0-1_amd64.deb
#   ovn-docker_2.8.0-1_amd64.deb
#   ovn-host_2.8.0-1_amd64.deb
#   python-openvswitch_2.8.0-1_all.deb

# /!\ Note this file is going to move elsewhere It's just temp.

set -e
WORK_DIR=/tmp/build-ovs
OVS_VERSION_SHORT="2.8.10"
OVS_VERSION="v${OVS_VERSION_SHORT}"
MAGMA_ROOT="$(realpath "$(dirname $0)"/../../../)"
GTP_PATCH_PATH="${MAGMA_ROOT}/third_party/gtp_ovs/kernel-4.9"
# Build time dependencies
BUILD_DEPS="graphviz debhelper dh-autoreconf python-all python-twisted-conch module-assistant git ruby-dev openssl pkg-config libssl-dev build-essential"
PATCHES="$(ls ${GTP_PATCH_PATH}/${OVS_VERSION_SHORT})"
FLOWBASED_PATH=$(readlink -f ${MAGMA_ROOT}/third_party/gtp_ovs/kernel-4.9/gtp-v4.9-backport/)
PATCH_ROOT=$(readlink -f "$GTP_PATCH_PATH/$OVS_VERSION_SHORT/")
VLAN_FIX="3cf2b424bb"
# be sure to increment this to enable upgrade from package repo when rebuilding identical upstream versions
LOCAL_REV=2

# The resulting package is placed in $OUTPUT_DIR
# or in the cwd.
if [ -z "$1" ]; then
  OUTPUT_DIR=$(pwd)
else
  OUTPUT_DIR=$1
  if [ ! -d "$OUTPUT_DIR" ]; then
    echo "error: $OUTPUT_DIR is not a valid directory. Exiting..."
    exit 1
  fi
fi

# build from source
if [ -d ${WORK_DIR} ]; then
  sudo rm -rf ${WORK_DIR}
fi
mkdir -p ${WORK_DIR}

# updating apt
sudo apt-get update
# install build time dependency
sudo apt-get install  ${BUILD_DEPS} -y
# make surew correct linux headers are installed
sudo apt-get -y install "linux-headers-$(uname -r)"
# Install fpm
sudo gem install fpm

# pull code and apply patches
cd ${WORK_DIR}
(wget "${OVS_ARCHIVE_URL}" && tar xzf "${OVS_ARCHIVE}") || git clone https://github.com/openvswitch/ovs.git
cd ovs
git checkout ${OVS_VERSION}
cp $PATCH_ROOT/*.patch $WORK_DIR/ovs
cp -r "${FLOWBASED_PATH}" "${WORK_DIR}/ovs/flow-based-gtp-linux-v4.9"
git apply ${PATCHES}
# vlan fix
git show $VLAN_FIX | git apply -3 -

# fakeroot finds version in debian/changelog
# -- build changelog entry with the correct version string and patch
cat <<EOF | git apply -
diff --git a/debian/changelog b/debian/changelog
index 824ed7d5f..b606d4331 100644
--- a/debian/changelog
+++ b/debian/changelog
@@ -1,3 +1,9 @@
+openvswitch (${OVS_VERSION_SHORT}-${LOCAL_REV}) unstable; urgency=low
+   [ local team ]
+   * New local version
+
+ -- local team <noreply@example.com>  $(date '+%a, %d %b %Y %H:%M:%S %z')
+
$(head -3 debian/changelog | sed 's/^/ /g')
EOF

./boot.sh
# Building OVS user packages
opts="--with-linux=/lib/modules/$(uname -r)/build"
deb_opts="nocheck parallel=$(nproc)"
fakeroot make -f debian/rules DATAPATH_CONFIGURE_OPTS="$opts" DEB_BUILD_OPTIONS="$deb_opts" binary

## Building OVS datapath kernel module package
cd ${WORK_DIR}/ovs
sudo mkdir -p /usr/src/linux
kvers=$(uname -r)
ksrc="/lib/modules/$kvers/build"
sudo make -f debian/rules.modules KSRC="$ksrc" KVERS="$kvers" binary-modules

cp /usr/src/*.deb ${WORK_DIR}/*.deb ${OUTPUT_DIR}
